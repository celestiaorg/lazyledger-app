package types

import (
	"bytes"
	"testing"

	sdkerrors "cosmossdk.io/errors"
	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWirePayForData_ValidateBasic(t *testing.T) {
	type test struct {
		name    string
		msg     *MsgWirePayForData
		wantErr *sdkerrors.Error
	}

	// valid wpfd
	validMsg := validWirePayForData(t)

	// wpfd with bad namespace id
	badIDMsg := validWirePayForData(t)
	badIDMsg.MessageNamespaceId = []byte{1, 2, 3, 4, 5, 6, 7}

	// wpfd that uses reserved namespace id
	reservedMsg := validWirePayForData(t)
	reservedMsg.MessageNamespaceId = []byte{0, 0, 0, 0, 0, 0, 0, 100}

	// wpfd that uses parity shares namespace id
	paritySharesMsg := validWirePayForData(t)
	paritySharesMsg.MessageNamespaceId = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	// wpfd that uses parity shares namespace id
	tailPaddingMsg := validWirePayForData(t)
	tailPaddingMsg.MessageNamespaceId = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE}

	// wpfd that has a wrong msg size
	invalidDeclaredMsgSizeMsg := validWirePayForData(t)
	invalidDeclaredMsgSizeMsg.MessageSize = 999

	// wpfd with bad message share commitment
	badCommitMsg := validWirePayForData(t)
	badCommitMsg.MessageShareCommitment = &ShareCommitAndSignature{
		ShareCommitment: []byte{1, 2, 3, 4},
	}

	tests := []test{
		{
			name:    "valid msg",
			msg:     validMsg,
			wantErr: nil,
		},
		{
			name:    "bad ns ID",
			msg:     badIDMsg,
			wantErr: ErrInvalidNamespaceLen,
		},
		{
			name:    "reserved ns id",
			msg:     reservedMsg,
			wantErr: ErrReservedNamespace,
		},
		{
			name:    "bad declared message size",
			msg:     invalidDeclaredMsgSizeMsg,
			wantErr: ErrDeclaredActualDataSizeMismatch,
		},
		{
			name:    "bad commitment",
			msg:     badCommitMsg,
			wantErr: ErrInvalidShareCommit,
		},
		{
			name:    "parity shares namespace id",
			msg:     paritySharesMsg,
			wantErr: ErrParitySharesNamespace,
		},
		{
			name:    "tail padding namespace id",
			msg:     tailPaddingMsg,
			wantErr: ErrTailPaddingNamespace,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
				space, code, log := sdkerrors.ABCIInfo(err, false)
				assert.Equal(t, tt.wantErr.Codespace(), space)
				assert.Equal(t, tt.wantErr.ABCICode(), code)
				t.Log(log)
			}
		})
	}
}

func TestProcessWirePayForData(t *testing.T) {
	type test struct {
		name          string
		namespace     []byte
		msg           []byte
		minSquareSize uint64
		expectErr     bool
		modify        func(*MsgWirePayForData) *MsgWirePayForData
	}

	dontModify := func(in *MsgWirePayForData) *MsgWirePayForData {
		return in
	}

	kb := generateKeyring(t, "test")

	signer := NewKeyringSigner(kb, "test", "chain-id")

	tests := []test{
		{
			name:          "single share square size 8",
			namespace:     []byte{1, 1, 1, 1, 1, 1, 1, 1},
			msg:           bytes.Repeat([]byte{1}, totalMsgSize(appconsts.SparseShareContentSize)),
			minSquareSize: 8,
			modify:        dontModify,
		},
		{
			name:          "12 shares square size 4",
			namespace:     []byte{1, 1, 1, 1, 1, 1, 1, 2},
			msg:           bytes.Repeat([]byte{2}, totalMsgSize(appconsts.SparseShareContentSize*12)),
			minSquareSize: 4,
			modify:        dontModify,
		},
		{
			name:          "nil signature",
			namespace:     []byte{1, 1, 1, 1, 1, 1, 1, 2},
			msg:           bytes.Repeat([]byte{2}, totalMsgSize(appconsts.SparseShareContentSize*12)),
			minSquareSize: 4,
			modify: func(wpfd *MsgWirePayForData) *MsgWirePayForData {
				wpfd.MessageShareCommitment.Signature = nil
				return wpfd
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		wpfd, err := NewWirePayForData(tt.namespace, tt.msg, int(tt.minSquareSize))
		require.NoError(t, err, tt.name)
		err = wpfd.SignMessageShareCommitment(signer)
		assert.NoError(t, err)

		wpfd = tt.modify(wpfd)

		message, spfd, sig, err := ProcessWirePayForData(wpfd, tt.minSquareSize)
		if tt.expectErr {
			assert.Error(t, err, tt.name)
			continue
		}

		// ensure that the shared fields are identical
		assert.Equal(t, tt.msg, message.Data, tt.name)
		assert.Equal(t, tt.namespace, message.NamespaceId, tt.name)
		assert.Equal(t, wpfd.Signer, spfd.Signer, tt.name)
		assert.Equal(t, wpfd.MessageNamespaceId, spfd.MessageNamespaceId, tt.name)
		assert.Equal(t, wpfd.MessageShareCommitment.ShareCommitment, spfd.MessageShareCommitment, tt.name)
		assert.Equal(t, wpfd.MessageShareCommitment.Signature, sig, tt.name)
	}
}
