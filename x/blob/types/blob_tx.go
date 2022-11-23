package types

import (
	"bytes"

	"github.com/celestiaorg/celestia-app/app/encoding"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// ProcessedBlobTx caches the unmarshalled result of the BlobTx
type ProcessedBlobTx struct {
	Tx    sdk.Tx
	Blobs []tmproto.Blob
	PFBs  []*MsgPayForBlob
}

// ProcessBlobTx validates the contents of the BlobTx and performs the
// malleation process by separating the blobs from the MsgPayForBlob.
func ProcessBlobTx(encfg encoding.Config, bTx *tmproto.BlobTx) (ProcessedBlobTx, error) {
	sdkTx, err := encfg.TxConfig.TxDecoder()(bTx.Tx)
	if err != nil {
		return ProcessedBlobTx{}, err
	}

	msgs := sdkTx.GetMsgs()
	pfbs := make([]*MsgPayForBlob, 0)
	for _, msg := range msgs {
		if sdk.MsgTypeURL(msg) == URLMsgPayForBlob {
			pfb, ok := msg.(*MsgPayForBlob)
			if !ok {
				continue
			}
			pfbs = append(pfbs, pfb)
		}
	}

	if len(pfbs) == 0 {
		return ProcessedBlobTx{}, ErrNoPFBInBlobTx
	}
	if len(pfbs) != len(bTx.Blobs) {
		return ProcessedBlobTx{}, ErrMismatchedNumberOfPFBorBlob
	}

	protoBlobs := make([]tmproto.Blob, len(pfbs))
	for i, pfb := range pfbs {
		err = pfb.ValidateBasic()
		if err != nil {
			return ProcessedBlobTx{}, err
		}

		blob := bTx.Blobs[i].Data

		// make sure that the blob size matches the actual size of the blob
		if pfb.BlobSize != uint64(len(blob)) {
			return ProcessedBlobTx{}, ErrDeclaredActualDataSizeMismatch.Wrapf(
				"declared: %d vs actual: %d",
				pfb.BlobSize,
				len(blob),
			)
		}

		// verify that the commitment of the blob matches that of the PFB
		calculatedCommit, err := CreateCommitment(pfb.NamespaceId, blob)
		if err != nil {
			return ProcessedBlobTx{}, err // todo: wrap this error with an sdkerror error
		}
		if !bytes.Equal(calculatedCommit, pfb.ShareCommitment) {
			return ProcessedBlobTx{}, ErrInvalidShareCommit
		}

		protoBlobs[i] = tmproto.Blob{NamespaceId: pfb.NamespaceId, Data: blob}
	}

	return ProcessedBlobTx{
		Tx:    sdkTx,
		Blobs: protoBlobs,
		PFBs:  pfbs,
	}, nil
}

func (btx ProcessedBlobTx) ValidateBasic() error {
	btx.Tx.ValidateBasic()

	return nil
}

func (btx ProcessedBlobTx) GetMsgs() []sdk.Msg {
	return btx.Tx.GetMsgs()
}
