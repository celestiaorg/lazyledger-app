package upgrade

import (
	"context"
	"encoding/binary"

	sdkmath "cosmossdk.io/math"
	"github.com/celestiaorg/celestia-app/x/upgrade/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Keeper implements the MsgServer and QueryServer interfaces
var (
	_ types.MsgServer   = &Keeper{}
	_ types.QueryServer = Keeper{}

	defaultSignalTheshold = Fraction{Numerator: 5, Denominator: 6}
)

type Fraction struct {
	Numerator   int64
	Denominator int64
}

// SignalThreshold is the fraction of voting power that is required
// to signal for a version change. It is set to 5/6 as the middle point
// between 2/3 and 3/3 providing 1/6 fault tolerance to halting the
// network during an upgrade period. It can be modified through a
// hard fork change that modified the app version
func SignalThreshold(version uint64) Fraction {
	switch version {
	default:
		return defaultSignalTheshold
	}
}

type Keeper struct {
	// we use the same upgrade store key so existing IBC client state can
	// safely be ported over without any migration
	storeKey storetypes.StoreKey

	// in memory copy of the upgrade height if any. This is local per node
	// and configured from the config. Used just for V2
	upgradeHeight int64

	// quorumVersion is the version that has received a quorum of validators
	// to signal for it. This variable is relevant just for the scope of the
	// lifetime of the block
	quorumVersion uint64

	// staking keeper is used to fetch validators to calculate the total power
	// signalled to a version
	stakingKeeper StakingKeeper
}

// NewKeeper constructs an upgrade keeper
func NewKeeper(
	storeKey storetypes.StoreKey,
	upgradeHeight int64,
	stakingKeeper StakingKeeper,
) Keeper {
	return Keeper{
		storeKey:      storeKey,
		upgradeHeight: upgradeHeight,
		stakingKeeper: stakingKeeper,
	}
}

// SignalVersion is a method required by the MsgServer interface
func (k Keeper) SignalVersion(ctx context.Context, req *types.MsgSignalVersion) (*types.MsgSignalVersionResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	// the signalled version must be either the current version (for cancelling an upgrade)
	// or the very next version (for accepting an upgrade)
	currentVersion := sdkCtx.BlockHeader().Version.App
	if req.Version != currentVersion && req.Version != currentVersion+1 {
		return nil, types.ErrInvalidVersion
	}

	// check that the validator exists
	_, found := k.stakingKeeper.GetValidator(sdkCtx, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	k.SetValidatorVersion(sdkCtx, valAddr, req.Version)
	return &types.MsgSignalVersionResponse{}, nil
}

// TryUpgrade is a method required by the MsgServer interface
// It tallies the voting power that has voted on each version.
// If one version has quorum, it is set as the quorum version
// which the application can use as signal to upgrade to that version.
func (k *Keeper) TryUpgrade(ctx context.Context, _ *types.MsgTryUpgrade) (*types.MsgTryUpgradeResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	threshold := k.GetVotingPowerThreshold(sdkCtx)
	hasQuorum, version := k.TallyVotingPower(sdkCtx, threshold.Int64())
	if hasQuorum {
		k.quorumVersion = version
	}
	return &types.MsgTryUpgradeResponse{}, nil
}

// VersionTally is a method required by the QueryServer interface
func (k Keeper) VersionTally(ctx context.Context, req *types.QueryVersionTallyRequest) (*types.QueryVersionTallyResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	totalVotingPower := k.stakingKeeper.GetLastTotalPower(sdkCtx)
	currentVotingPower := sdk.NewInt(0)
	store := sdkCtx.KVStore(k.storeKey)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		valAddress := sdk.ValAddress(iterator.Key())
		power := k.stakingKeeper.GetLastValidatorPower(sdkCtx, valAddress)
		version := VersionFromBytes(iterator.Value())
		if version == req.Version {
			currentVotingPower = currentVotingPower.AddRaw(power)
		}
	}
	threshold := k.GetVotingPowerThreshold(sdkCtx)
	return &types.QueryVersionTallyResponse{
		VotingPower:      currentVotingPower.Uint64(),
		ThresholdPower:   threshold.Uint64(),
		TotalVotingPower: totalVotingPower.Uint64(),
	}, nil
}

// SetValidatorVersion saves a signalled version for a validator using the keeper's store key
func (k Keeper) SetValidatorVersion(ctx sdk.Context, valAddress sdk.ValAddress, version uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(valAddress, VersionToBytes(version))
}

// DeleteValidatorVersion deletes a signalled version for a validator using the keeper's store key
func (k Keeper) DeleteValidatorVersion(ctx sdk.Context, valAddress sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(valAddress)
}

// TallyVotingPower tallies the voting power for each version and returns true if
// any version has reached the quorum in voting power
func (k Keeper) TallyVotingPower(ctx sdk.Context, threshold int64) (bool, uint64) {
	output := make(map[uint64]int64)
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		valAddress := sdk.ValAddress(iterator.Key())
		// check that the validator is still part of the bonded set
		val, found := k.stakingKeeper.GetValidator(ctx, valAddress)
		if !found {
			// if it no longer exists, delete the version
			k.DeleteValidatorVersion(ctx, valAddress)
		}
		// if the validator is not bonded, skip it's voting power
		if !found || !val.IsBonded() {
			continue
		}
		power := k.stakingKeeper.GetLastValidatorPower(ctx, valAddress)
		version := VersionFromBytes(iterator.Value())
		if _, ok := output[version]; !ok {
			output[version] = power
		} else {
			output[version] += power
		}
		if output[version] >= threshold {
			return true, version
		}
	}
	return false, 0
}

// GetVotingPowerThreshold returns the voting power threshold required to
// upgrade to a new version.
func (k Keeper) GetVotingPowerThreshold(ctx sdk.Context) sdkmath.Int {
	// contract: totalVotingPower should not exceed MaxUint64
	totalVotingPower := k.stakingKeeper.GetLastTotalPower(ctx)
	thresholdFraction := SignalThreshold(ctx.BlockHeader().Version.App)
	product := totalVotingPower.MulRaw(thresholdFraction.Numerator)
	if product.ModRaw(thresholdFraction.Denominator).IsZero() {
		return product.QuoRaw(thresholdFraction.Denominator)
	}
	return product.QuoRaw(thresholdFraction.Denominator).AddRaw(1)
}

// ShouldUpgradeToV2 returns true if the current height is one before
// the locally provided upgrade height that is passed as a flag
// NOTE: This is only used to upgrade to v2 and should be deprecated
// in v3
func (k Keeper) ShouldUpgradeToV2(height int64) bool {
	return k.upgradeHeight == height+1
}

// ShouldUpgrade returns true if the signalling mechanism has concluded
// that the network is ready to upgrade. It also returns the version
// that the node should upgrade to.
func (k *Keeper) ShouldUpgrade() (bool, uint64) {
	return k.quorumVersion != 0, k.quorumVersion
}

// ResetTally resets the tally after a version change. It iterates over the store,
// and deletes any versions that are less than the provided version. It also
// resets the quorumVersion to 0.
func (k *Keeper) ResetTally(ctx sdk.Context, version uint64) {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		v := VersionFromBytes(iterator.Value())
		if v <= version {
			store.Delete(iterator.Key())
		}
	}
	k.quorumVersion = 0
}

func VersionToBytes(version uint64) []byte {
	return binary.BigEndian.AppendUint64(nil, version)
}

func VersionFromBytes(version []byte) uint64 {
	return binary.BigEndian.Uint64(version)
}
