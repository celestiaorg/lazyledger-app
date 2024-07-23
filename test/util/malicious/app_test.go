package malicious

import (
	"testing"

	"github.com/celestiaorg/celestia-app/v3/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/v3/pkg/da"
	"github.com/celestiaorg/celestia-app/v3/pkg/wrapper"
	"github.com/celestiaorg/celestia-app/v3/test/util/blobfactory"
	"github.com/celestiaorg/celestia-app/v3/test/util/testfactory"
	"github.com/celestiaorg/celestia-app/v3/test/util/testnode"
	"github.com/celestiaorg/go-square/shares"
	"github.com/celestiaorg/go-square/square"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"
)

// TestOutOfOrderNMT tests that the malicious NMT implementation is able to
// generate the same root as the ordered NMT implementation when the leaves are
// added in the same order and can generate roots when leaves are out of
// order.
func TestOutOfOrderNMT(t *testing.T) {
	squareSize := uint64(64)
	c := NewConstructor(squareSize)
	goodConstructor := wrapper.NewConstructor(squareSize)

	orderedTree := goodConstructor(0, 0)
	maliciousOrderedTree := c(0, 0)
	maliciousUnorderedTree := c(0, 0)
	data := testfactory.GenerateRandNamespacedRawData(64)

	// compare the roots generated by pushing ordered data
	for _, d := range data {
		err := orderedTree.Push(d)
		require.NoError(t, err)
		err = maliciousOrderedTree.Push(d)
		require.NoError(t, err)
	}

	goodOrderedRoot, err := orderedTree.Root()
	require.NoError(t, err)
	malOrderedRoot, err := maliciousOrderedTree.Root()
	require.NoError(t, err)
	require.Equal(t, goodOrderedRoot, malOrderedRoot)

	// test the new tree with unordered data
	for i := range data {
		j := tmrand.Intn(len(data))
		data[i], data[j] = data[j], data[i]
	}

	for _, d := range data {
		err := maliciousUnorderedTree.Push(d)
		require.NoError(t, err)
	}

	root, err := maliciousUnorderedTree.Root()
	require.NoError(t, err)
	require.Len(t, root, 90)                   // two namespaces + 32 bytes of hash = 90 bytes
	require.NotEqual(t, goodOrderedRoot, root) // quick sanity check to ensure the roots are different
}

// TestMaliciousTestNode runs a single validator network using the malicious
// node. This will begin to produce out of order blocks after block height of 5.
func TestMaliciousTestNode(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping MaliciousTestNode in short mode.")
	}
	accounts := testfactory.RandomAccountNames(5)
	cfg := OutOfOrderNamespaceConfig(5).
		WithFundedAccounts(accounts...)

	cctx, _, _ := testnode.NewNetwork(t, cfg)
	_, err := cctx.WaitForHeight(6)
	require.NoError(t, err)

	// submit a multiblob tx where each blob is using a random namespace. This
	// will result in the first two blobs being swapped in the square as per the
	// malicious square builder.
	client, err := testnode.NewTxClientFromContext(cctx)
	require.NoError(t, err)
	blobs := blobfactory.ManyRandBlobs(tmrand.NewRand(), 10_000, 10_000, 10_000, 10_000, 10_000, 10_000, 10_000)
	txres, err := client.SubmitPayForBlob(cctx.GoContext(), blobs, blobfactory.DefaultTxOpts()...)
	require.NoError(t, err)
	require.Equal(t, abci.CodeTypeOK, txres.Code)

	// fetch the block that included in the tx
	inclusionHeight := txres.Height
	block, err := cctx.Client.Block(cctx.GoContext(), &inclusionHeight)
	require.NoError(t, err)

	// check that we can recalculate the data root using the malicious code but
	// not the correct code
	s, err := Construct(block.Block.Txs.ToSliceOfBytes(), appconsts.LatestVersion, appconsts.DefaultSquareSizeUpperBound, OutOfOrderExport)
	require.NoError(t, err)

	rawSquare := shares.ToBytes(s)
	eds, err := ExtendShares(rawSquare)
	require.NoError(t, err)

	dah, err := da.NewDataAvailabilityHeader(eds)
	require.NoError(t, err)
	require.Equal(t, block.Block.DataHash.Bytes(), dah.Hash())

	correctSquare, err := square.Construct(block.Block.Txs.ToSliceOfBytes(),
		appconsts.DefaultSquareSizeUpperBound,
		appconsts.DefaultSubtreeRootThreshold,
	)
	require.NoError(t, err)

	goodEds, err := da.ExtendShares(shares.ToBytes(correctSquare))
	require.NoError(t, err)

	goodDah, err := da.NewDataAvailabilityHeader(goodEds)
	require.NoError(t, err)
	require.NotEqual(t, block.Block.DataHash.Bytes(), goodDah.Hash())
}
