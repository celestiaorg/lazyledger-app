package testnode

import (
	"encoding/json"
	"time"

	"github.com/celestiaorg/celestia-app/app"
	blobtypes "github.com/celestiaorg/celestia-app/x/blob/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// GenesisOption allows for arbitrary changes to be made on the genesis state
// after initial accounts have been added. It accepts the genesis state as input
// and is expected to return the modifed genesis as output.
type GenesisOption func(state map[string]json.RawMessage) map[string]json.RawMessage

// SetBlobParams will set the provided blob params as genesis state.
func SetBlobParams(codec codec.Codec, params blobtypes.Params) GenesisOption {
	return func(state map[string]json.RawMessage) map[string]json.RawMessage {
		blobGenState := blobtypes.DefaultGenesis()
		blobGenState.Params = params
		state[blobtypes.ModuleName] = codec.MustMarshalJSON(blobGenState)
		return state
	}
}

// ImmediateProposals sets the thresholds for getting a gov proposal to very low
// levels.
func ImmediateProposals(codec codec.Codec) GenesisOption {
	return func(state map[string]json.RawMessage) map[string]json.RawMessage {
		gs := v1.DefaultGenesisState()
		gs.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(app.BondDenom, sdk.NewInt(10)))
		gs.TallyParams.Quorum = "0.000001"
		gs.TallyParams.Threshold = "0.000001"
		vp := time.Second * 5
		gs.VotingParams.VotingPeriod = &vp
		state[govtypes.ModuleName] = codec.MustMarshalJSON(gs)
		return state
	}
}

// FundAccounts adds a set of accounts to the genesis and then sets their balance as provided.
// This is good in the case where you have a separate keyring you want to test against and not
// use the one generated by the testnet infra.
func FundAccounts(codec codec.Codec, addresses []sdk.AccAddress, balance sdk.Coin) GenesisOption {
	return func(state map[string]json.RawMessage) map[string]json.RawMessage {
		// set the accounts in the genesis state
		var authGenState authtypes.GenesisState
		codec.MustUnmarshalJSON(state[authtypes.ModuleName], &authGenState)

		genAccounts := make([]authtypes.GenesisAccount, len(addresses))
		genBalances := make([]banktypes.Balance, len(addresses))
		for idx, addr := range addresses {
			genAccounts[idx] = authtypes.NewBaseAccount(addr, nil, uint64(idx+len(authGenState.Accounts)), 0)
			genBalances[idx] = banktypes.Balance{Address: addr.String(), Coins: sdk.NewCoins(balance)}
		}

		accounts, err := authtypes.PackAccounts(genAccounts)
		if err != nil {
			panic(err)
		}

		authGenState.Accounts = append(authGenState.Accounts, accounts...)
		state[authtypes.ModuleName] = codec.MustMarshalJSON(&authGenState)

		// set the balances in the genesis state
		var bankGenState banktypes.GenesisState
		codec.MustUnmarshalJSON(state[banktypes.ModuleName], &bankGenState)

		bankGenState.Balances = append(bankGenState.Balances, genBalances...)
		state[banktypes.ModuleName] = codec.MustMarshalJSON(&bankGenState)
		return state
	}
}
