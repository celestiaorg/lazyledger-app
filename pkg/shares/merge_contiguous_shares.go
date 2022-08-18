package shares

import (
	"encoding/binary"
	"errors"

	"github.com/tendermint/tendermint/pkg/consts"
)

// processContiguousShares takes raw shares and extracts out transactions,
// intermediate state roots, or evidence. The returned [][]byte do have
// namespaces or length delimiters and are ready to be unmarshalled
func processContiguousShares(shares [][]byte) (txs [][]byte, err error) {
	if len(shares) == 0 {
		return nil, nil
	}

	ss := newShareStack(shares)
	return ss.resolve()
}

// shareStack hold variables for peel
type shareStack struct {
	shares [][]byte
	txLen  uint64
	txs    [][]byte
	cursor int
}

func newShareStack(shares [][]byte) *shareStack {
	return &shareStack{shares: shares}
}

func (ss *shareStack) resolve() ([][]byte, error) {
	if len(ss.shares) == 0 {
		return nil, nil
	}
	err := ss.peel(ss.shares[0][consts.NamespaceSize+consts.ShareReservedBytes:], true)
	return ss.txs, err
}

// peel recursively parses each chunk of data (either a transaction,
// intermediate state root, or evidence) and adds it to the underlying slice of data.
func (ss *shareStack) peel(share []byte, delimited bool) (err error) {
	if delimited {
		var txLen uint64
		share, txLen, err = ParseDelimiter(share)
		if err != nil {
			return err
		}
		if txLen == 0 {
			return nil
		}
		ss.txLen = txLen
	}
	// safeLen describes the point in the share where it can be safely split. If
	// split beyond this point, it is possible to break apart a length
	// delimiter, which will result in incorrect share merging
	safeLen := len(share) - binary.MaxVarintLen64
	if safeLen < 0 {
		safeLen = 0
	}
	if ss.txLen <= uint64(safeLen) {
		ss.txs = append(ss.txs, share[:ss.txLen])
		share = share[ss.txLen:]
		return ss.peel(share, true)
	}
	// add the next share to the current share to continue merging if possible
	if len(ss.shares) > ss.cursor+1 {
		ss.cursor++
		share := append(share, ss.shares[ss.cursor][consts.NamespaceSize+consts.ShareReservedBytes:]...)
		return ss.peel(share, false)
	}
	// collect any remaining data
	if ss.txLen <= uint64(len(share)) {
		ss.txs = append(ss.txs, share[:ss.txLen])
		share = share[ss.txLen:]
		return ss.peel(share, true)
	}
	return errors.New("failure to parse block data: transaction length exceeded data length")
}
