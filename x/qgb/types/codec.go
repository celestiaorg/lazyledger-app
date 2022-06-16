package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgDataCommitmentConfirm{}, "qgb/DataCommitmentConfirm", nil)
	cdc.RegisterConcrete(&MsgValsetConfirm{}, "qgb/MsgValSetConfirm", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgDataCommitmentConfirm{},
		&MsgValsetConfirm{},
	)

	// TODO check if it works without
	registry.RegisterInterface("AttestationRequestI", (*AttestationRequestI)(nil))

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
