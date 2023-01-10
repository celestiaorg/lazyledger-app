package app

import (
	"math"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/pkg/shares"
	coretypes "github.com/tendermint/tendermint/types"
)

// estimateSquareSize uses the provided block data to over estimate the square
// size and the starting share index of non-reserved namespaces. The estimates
// returned are liberal in the sense that we assume close to worst case and
// round up.
//
// NOTE: The estimation process does not have to be perfect. We can overestimate
// because the cost of padding is limited.
func estimateSquareSize(txs []parsedTx) (squareSize uint64, nonreserveStart int) {
	txSharesUsed := estimateTxSharesUsed(txs)
	pfbTxSharesUsed := estimatePFBTxSharesUsed(appconsts.DefaultMaxSquareSize, txs)
	blobSharesUsed := 0

	for _, ptx := range txs {
		if len(ptx.normalTx) != 0 {
			continue
		}
		blobSharesUsed += shares.SparseSharesNeeded(uint32(ptx.blobTx.DataUsed()))
	}

	// assume that we have to add a lot of padding by simply doubling the number
	// of shares used
	//
	// TODO: use a more precise estimation that doesn't over
	// estimate as much
	totalSharesUsed := uint64(txSharesUsed + pfbTxSharesUsed + blobSharesUsed)
	totalSharesUsed *= 2
	minSize := uint64(math.Sqrt(float64(totalSharesUsed)))
	squareSize = shares.RoundUpPowerOfTwo(minSize)
	if squareSize >= appconsts.DefaultMaxSquareSize {
		squareSize = appconsts.DefaultMaxSquareSize
	}
	if squareSize <= appconsts.DefaultMinSquareSize {
		squareSize = appconsts.DefaultMinSquareSize
	}

	return squareSize, txSharesUsed + pfbTxSharesUsed
}

// estimateTxSharesUsed estimates the number of shares used by ordinary
// transactions (i.e. all transactions that aren't PFBs).
func estimateTxSharesUsed(ptxs []parsedTx) int {
	txBytes := 0
	for _, pTx := range ptxs {
		if pTx.isNormalTx() {
			txLen := len(pTx.normalTx)
			txLen += shares.DelimLen(uint64(txLen))
			txBytes += txLen
		}
	}
	return estimateCompactShares(txBytes)
}

// estimatePFBTxSharesUsed estimates the number of shares used by PFB
// transactions.
func estimatePFBTxSharesUsed(squareSize uint64, ptxs []parsedTx) int {
	maxWTxOverhead := maxWrappedTxOverhead(squareSize)
	txBytes := 0
	for _, pTx := range ptxs {
		if pTx.isBlobTx() {
			txLen := len(pTx.blobTx.Tx) + maxWTxOverhead
			txLen += shares.DelimLen(uint64(txLen))
			txBytes += txLen
		}
	}
	return estimateCompactShares(txBytes)
}

// estimateCompactShares estimates the number of shares used by compact shares
func estimateCompactShares(totalBytes int) int {
	if totalBytes == 0 {
		return 0
	}
	if totalBytes <= appconsts.FirstCompactShareContentSize {
		return 1
	}
	// account for the first share
	sharesUsed := 1
	totalBytes -= appconsts.FirstCompactShareContentSize

	// account for continuation shares
	sharesUsed += (totalBytes / appconsts.ContinuationCompactShareContentSize)
	if totalBytes%appconsts.ContinuationCompactShareContentSize != 0 {
		sharesUsed++
	}

	return sharesUsed
}

// maxWrappedTxOverhead calculates the maximum amount of overhead introduced by
// wrapping a transaction with a shares index
//
// TODO: make more efficient by only generating these numbers once or something
// similar. This function alone can take up to 5ms.
func maxWrappedTxOverhead(squareSize uint64) int {
	maxTxLen := squareSize * squareSize * appconsts.ContinuationCompactShareContentSize
	wtx, err := coretypes.MarshalIndexWrapper(
		uint32(squareSize*squareSize),
		make([]byte, maxTxLen))
	if err != nil {
		panic(err)
	}
	return len(wtx) - int(maxTxLen)
}
