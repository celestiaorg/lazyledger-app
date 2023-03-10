package shares

import (
	"encoding/binary"
	"fmt"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/celestiaorg/nmt/namespace"
	coretypes "github.com/tendermint/tendermint/types"
)

type ShareRange struct {
	// Start is the index of the first share occupied by this range.
	Start int
	// End is the index of the last share occupied by this range.
	End int
}

// CompactShareSplitter will write raw data compactly across a progressively
// increasing set of shares. It is used to lazily split block data such as
// transactions or intermediate state roots into shares.
type CompactShareSplitter struct {
	shares []Share
	// pendingShare Share
	shareBuilder *Builder
	namespace    namespace.ID
	done         bool
	shareVersion uint8
	// shareRanges is a map from a transaction key to the range of shares it
	// occupies. The range assumes this compact share splitter is the only
	// thing in the data square (e.g. the range for the first tx starts at index
	// 0).
	shareRanges map[coretypes.TxKey]ShareRange
}

// NewCompactShareSplitter returns a CompactShareSplitter using the provided
// namespace and shareVersion.
func NewCompactShareSplitter(ns namespace.ID, shareVersion uint8) *CompactShareSplitter {
	return &CompactShareSplitter{
		shares:       []Share{},
		namespace:    ns,
		shareVersion: shareVersion,
		shareRanges:  map[coretypes.TxKey]ShareRange{},
		shareBuilder: NewBuilder(ns, shareVersion, true, true),
	}
}

// WriteTx adds the delimited data for the provided tx to the underlying compact
// share splitter.
func (css *CompactShareSplitter) WriteTx(tx coretypes.Tx) {
	rawData, err := MarshalDelimitedTx(tx)
	if err != nil {
		panic(fmt.Sprintf("included Tx in mem-pool that can not be encoded %v", tx))
	}

	startShare := len(css.shares)
	css.write(rawData)
	endShare := css.Count() - 1

	css.shareRanges[tx.Key()] = ShareRange{
		Start: startShare,
		End:   endShare,
	}
}

// write adds the delimited data to the underlying compact shares.
func (css *CompactShareSplitter) write(rawData []byte) {
	if css.done {
		// remove the last element
		if !css.shareBuilder.IsEmptyShare() {
			css.shares = css.shares[:len(css.shares)-1]
		}
		css.done = false
	}

	if err := css.shareBuilder.MaybeWriteReservedBytes(); err != nil {
		panic(err)
	}

	for {

		rawDataLeftOver := css.shareBuilder.AddData(rawData)
		if rawDataLeftOver == nil {
			break
		}
		css.stackPending()

		rawData = rawDataLeftOver
	}

	if css.shareBuilder.AvailableBytes() == 0 {
		css.stackPending()
	}
}

// stackPending will build & add the pending share to accumulated shares
func (css *CompactShareSplitter) stackPending() {
	pendingShare, err := css.shareBuilder.Build()
	if err != nil {
		panic(err)
	}
	css.shares = append(css.shares, *pendingShare)

	// Now we need to create a new builder
	css.shareBuilder = NewBuilder(css.namespace, css.shareVersion, true, false)
}

// Export finalizes and returns the underlying compact shares and a map of
// shareRanges. All share ranges in the map of shareRanges will be offset (i.e.
// incremented) by the shareRangeOffset provided. shareRangeOffset should be 0
// for the first compact share sequence in the data square (transactions) but
// should be some non-zero number for subsequent compact share sequences (e.g.
// pfb txs).
func (css *CompactShareSplitter) Export(shareRangeOffset int) ([]Share, map[coretypes.TxKey]ShareRange) {
	// apply the shareRangeOffset to all share ranges
	shareRanges := make(map[coretypes.TxKey]ShareRange, len(css.shareRanges))

	if css.isEmpty() {
		return []Share{}, shareRanges
	}

	for k, v := range css.shareRanges {
		shareRanges[k] = ShareRange{
			Start: v.Start + shareRangeOffset,
			End:   v.End + shareRangeOffset,
		}
	}

	// in case Export is called multiple times
	if css.done {
		return css.shares, shareRanges
	}

	var bytesOfPadding int
	// add the pending share to the current shares before returning
	if !css.shareBuilder.IsEmptyShare() {
		bytesOfPadding = css.shareBuilder.ZeroPadIfNecessary()
		css.stackPending()
	}

	sequenceLen := css.sequenceLen(bytesOfPadding)
	css.writeSequenceLen(sequenceLen)
	css.done = true
	return css.shares, shareRanges
}

// writeSequenceLen writes the sequence length to the first share.
func (css *CompactShareSplitter) writeSequenceLen(sequenceLen uint32) {
	if css.isEmpty() {
		return
	}

	// We may find a more efficient way to write seqLen
	b := NewBuilder(css.namespace, css.shareVersion, true, true)
	b.ImportRawShare(css.shares[0])

	if err := b.WriteSequenceLen(sequenceLen); err != nil {
		panic(err)
	}

	firstShare, err := b.Build()
	if err != nil {
		panic(err)
	}

	// replace existing first share with new first share
	css.shares[0] = *firstShare
}

// sequenceLen returns the total length in bytes of all units (transactions or
// intermediate state roots) written to this splitter. sequenceLen does not
// include the number of bytes occupied by the namespace ID, the share info
// byte, or the reserved bytes. sequenceLen does include the unit length
// delimiter prefixed to each unit.
func (css *CompactShareSplitter) sequenceLen(bytesOfPadding int) uint32 {
	if len(css.shares) == 0 {
		return 0
	}
	if len(css.shares) == 1 {
		return uint32(appconsts.FirstCompactShareContentSize) - uint32(bytesOfPadding)
	}

	continuationSharesCount := len(css.shares) - 1
	continuationSharesSequenceLen := continuationSharesCount * appconsts.ContinuationCompactShareContentSize
	return uint32(appconsts.FirstCompactShareContentSize + continuationSharesSequenceLen - bytesOfPadding)
}

// isEmpty returns whether this compact share splitter is empty.
func (css *CompactShareSplitter) isEmpty() bool {
	return len(css.shares) == 0 && css.shareBuilder.IsEmptyShare()
}

// Count returns the number of shares that would be made if `Export` was invoked
// on this compact share splitter.
func (css *CompactShareSplitter) Count() (shareCount int) {
	if !css.shareBuilder.IsEmptyShare() && !css.done {
		// pending share is non-empty, so it will be zero padded and added to shares during export
		return len(css.shares) + 1
	}
	return len(css.shares)
}

// MarshalDelimitedTx prefixes a transaction with the length of the transaction
// encoded as a varint.
func MarshalDelimitedTx(tx coretypes.Tx) ([]byte, error) {
	lenBuf := make([]byte, binary.MaxVarintLen64)
	length := uint64(len(tx))
	n := binary.PutUvarint(lenBuf, length)
	return append(lenBuf[:n], tx...), nil
}
