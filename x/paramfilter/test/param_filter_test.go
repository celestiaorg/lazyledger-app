package test

import (
	"testing"

	"github.com/celestiaorg/celestia-app/app"
	testutil "github.com/celestiaorg/celestia-app/test/util"
	"github.com/celestiaorg/celestia-app/x/paramfilter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/stretchr/testify/require"
	tmlog "github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestParamFilter(t *testing.T) {
	app, _ := testutil.SetupTestAppWithGenesisValSet(app.DefaultConsensusParams())

	require.Greater(t, len(app.BlockedParams()), 0)

	// check that all blocked parameters are in the filter keeper
	pph := paramfilter.NewParamBlockList(app.BlockedParams()...)
	for _, p := range app.BlockedParams() {
		require.True(t, pph.IsBlocked(p[0], p[1]))
	}

	handler := pph.GovHandler(app.ParamsKeeper)
	ctx := sdk.NewContext(app.CommitMultiStore(), types.Header{}, false, tmlog.NewNopLogger())

	for _, p := range app.BlockedParams() {
		p := testProposal(proposal.NewParamChange(p[0], p[1], "value"))
		err := handler(ctx, p)
		require.Error(t, err)
		require.Contains(t, err.Error(), "parameter can not be modified")
	}
}

func testProposal(changes ...proposal.ParamChange) *proposal.ParameterChangeProposal {
	return proposal.NewParameterChangeProposal("title", "description", changes)
}
