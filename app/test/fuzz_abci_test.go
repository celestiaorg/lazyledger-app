package app_test

import (
	"testing"

	"github.com/celestiaorg/celestia-app/app"
	"github.com/celestiaorg/celestia-app/app/encoding"
	"github.com/celestiaorg/celestia-app/testutil"
	"github.com/celestiaorg/celestia-app/testutil/blobfactory"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	core "github.com/tendermint/tendermint/proto/tendermint/types"
	coretypes "github.com/tendermint/tendermint/types"
)

// TestPrepareProposalConsistency produces blocks with random data using
// PrepareProposal and then tests those blocks by calling ProcessProposal. All
// blocks produced by PrepareProposal should be accepted by ProcessProposal. It
// doesn't use the standard go tools for fuzzing as those tools only support
// fuzzing limited types, instead we create blocks our selves using random
// transction.
func TestPrepareProposalConsistency(t *testing.T) {
	encConf := encoding.MakeConfig(app.ModuleEncodingRegisters...)
	testApp, _ := testutil.SetupTestAppWithGenesisValSet()
	for i := 0; i < 200; i++ {
		ProcessRandomProposal(t, 100, 50000, encConf, testApp)
	}
}

func ProcessRandomProposal(
	t *testing.T,
	count,
	maxSize int,
	cfg encoding.Config,
	capp *app.App,
) {
	txs := blobfactory.RandBlobTxsRandomlySized(cfg.TxConfig.TxEncoder(), count, maxSize)
	sendTxs := blobfactory.GenerateManyRawSendTxs(cfg.TxConfig, count)
	txs = append(txs, sendTxs...)
	resp := capp.PrepareProposal(abci.RequestPrepareProposal{
		BlockData: &core.Data{
			Txs: coretypes.Txs(txs).ToSliceOfBytes(),
		},
	})
	res := capp.ProcessProposal(abci.RequestProcessProposal{
		BlockData: resp.BlockData,
		Header: core.Header{
			DataHash: resp.BlockData.Hash,
		},
	})
	require.Equal(t, abci.ResponseProcessProposal_ACCEPT, res.Result)
}
