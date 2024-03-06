package ante

import (
	"context"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.AnteDecorator      = MsgVersioningGateKeeper{}
	_ baseapp.CircuitBreaker = MsgVersioningGateKeeper{}
)

// MsgVersioningGateKeeper dictates which transactions are accepted for an app version
type MsgVersioningGateKeeper struct {
	acceptedMsgs map[uint64]map[string]struct{}
}

func NewMsgVersioningGateKeeper(acceptedList map[uint64]map[string]struct{}) *MsgVersioningGateKeeper {
	return &MsgVersioningGateKeeper{
		acceptedMsgs: acceptedList,
	}
}

// AnteHandle implements the ante.Decorator interface
func (mgk MsgVersioningGateKeeper) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	acceptedMsgs, exists := mgk.acceptedMsgs[ctx.BlockHeader().Version.App]
	if !exists {
		return ctx, sdkerrors.ErrNotSupported.Wrapf("app version %d is not supported", ctx.BlockHeader().Version.App)
	}
	for _, msg := range tx.GetMsgs() {
		msgTypeURL := sdk.MsgTypeURL(msg)
		_, exists := acceptedMsgs[msgTypeURL]
		if !exists {
			return ctx, sdkerrors.ErrNotSupported.Wrapf("transaction type %s is not supported in version %d", msgTypeURL, ctx.BlockHeader().Version.App)
		}
	}

	return next(ctx, tx, simulate)
}

func (mgk MsgVersioningGateKeeper) IsAllowed(ctx context.Context, msgName string) (bool, error) {
	appVersion := sdk.UnwrapSDKContext(ctx).BlockHeader().Version.App
	acceptedMsgs, exists := mgk.acceptedMsgs[appVersion]
	if !exists {
		return false, sdkerrors.ErrNotSupported.Wrapf("app version %d is not supported", appVersion)
	}
	_, exists = acceptedMsgs[msgName]
	if !exists {
		return false, nil
	}
	return true, nil
}
