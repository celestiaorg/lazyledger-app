package shares

import (
	"bytes"
	"testing"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	testutilns "github.com/celestiaorg/celestia-app/testutil/namespace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	coretypes "github.com/tendermint/tendermint/types"
)

// TestPadFirstIndexedBlob ensures that we are adding padding to the first share
// instead of calculating the value.
func TestPadFirstIndexedBlob(t *testing.T) {
	tx := tmrand.Bytes(300)
	blob := tmrand.Bytes(300)
	index := 100
	indexedTx, err := coretypes.MarshalIndexWrapper(tx, 100)
	require.NoError(t, err)

	bd := coretypes.Data{
		Txs: []coretypes.Tx{indexedTx},
		Blobs: []coretypes.Blob{
			{
				NamespaceID:  testutilns.RandomBlobNamespace(),
				Data:         blob,
				ShareVersion: appconsts.ShareVersionZero,
			},
		},
		SquareSize: 64,
	}

	shares, err := Split(bd, true)
	require.NoError(t, err)

	resShare, err := shares[index].RawData()
	require.NoError(t, err)

	require.True(t, bytes.Contains(resShare, blob))
}

func TestSequenceLen(t *testing.T) {
	type testCase struct {
		name    string
		share   Share
		wantLen uint32
		wantErr bool
	}
	sparseNamespaceID := bytes.Repeat([]byte{1}, appconsts.NamespaceSize)
	firstShare := append(sparseNamespaceID,
		[]byte{
			1,           // info byte
			0, 0, 0, 10, // sequence len
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, // data
		}...)
	firstShareWithLongSequence := append(sparseNamespaceID,
		[]byte{
			1,           // info byte
			0, 0, 1, 67, // sequence len
		}...)
	continuationShare := append(sparseNamespaceID,
		[]byte{
			0, // info byte
		}...)
	compactShare := append([]byte(appconsts.TxNamespaceID),
		[]byte{
			1,           // info byte
			0, 0, 0, 10, // sequence len
		}...)
	noInfoByte := []byte(appconsts.TxNamespaceID)
	noSequenceLen := append([]byte(appconsts.TxNamespaceID),
		[]byte{
			1, // info byte
		}...)
	testCases := []testCase{
		{
			name:    "first share",
			share:   firstShare,
			wantLen: 10,
			wantErr: false,
		},
		{
			name:    "first share with long sequence",
			share:   firstShareWithLongSequence,
			wantLen: 323,
			wantErr: false,
		},
		{
			name:    "continuation share",
			share:   continuationShare,
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "compact share",
			share:   compactShare,
			wantLen: 10,
			wantErr: false,
		},
		{
			name:    "no info byte returns error",
			share:   noInfoByte,
			wantLen: 0,
			wantErr: true,
		},
		{
			name:    "no sequence len returns error",
			share:   noSequenceLen,
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			len, err := tc.share.SequenceLen()

			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			if tc.wantLen != len {
				t.Errorf("want %d, got %d", tc.wantLen, len)
			}
		})
	}
}

func TestRawData(t *testing.T) {
	type testCase struct {
		name    string
		share   Share
		want    []byte
		wantErr bool
	}
	sparseNamespaceID := bytes.Repeat([]byte{1}, appconsts.NamespaceSize)
	firstSparseShare := append(
		sparseNamespaceID,
		[]byte{
			1,           // info byte
			0, 0, 0, 10, // sequence len
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, // data
		}...)
	continuationSparseShare := append(
		sparseNamespaceID,
		[]byte{
			0,                             // info byte
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, // data
		}...)
	firstCompactShare := append([]byte(appconsts.TxNamespaceID),
		[]byte{
			1,           // info byte
			0, 0, 0, 10, // sequence len
			0, 0, 0, 15, // reserved bytes
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, // data
		}...)
	continuationCompactShare := append([]byte(appconsts.TxNamespaceID),
		[]byte{
			0,          // info byte
			0, 0, 0, 0, // reserved bytes
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, // data
		}...)
	noSequenceLen := append([]byte(appconsts.TxNamespaceID),
		[]byte{
			1, // info byte
		}...)
	notEnoughSequenceLenBytes := append([]byte(appconsts.TxNamespaceID),
		[]byte{
			1,        // info byte
			0, 0, 10, // sequence len
		}...)
	testCases := []testCase{
		{
			name:  "first sparse share",
			share: firstSparseShare,
			want:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:  "continuation sparse share",
			share: continuationSparseShare,
			want:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:  "first compact share",
			share: firstCompactShare,
			want:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:  "continuation compact share",
			share: continuationCompactShare,
			want:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:    "no sequence len returns error",
			share:   noSequenceLen,
			wantErr: true,
		},
		{
			name:    "not enough sequence len bytes returns error",
			share:   notEnoughSequenceLenBytes,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rawData, err := tc.share.RawData()
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tc.want, rawData)
		})
	}
}

func TestIsCompactShare(t *testing.T) {
	type testCase struct {
		name  string
		share Share
		want  bool
	}

	txShare, _ := zeroPadIfNecessary(appconsts.TxNamespaceID, appconsts.ShareSize)
	pfbTxShare, _ := zeroPadIfNecessary(appconsts.PayForBlobNamespaceID, appconsts.ShareSize)
	blobShare, _ := zeroPadIfNecessary(namespaceOne, appconsts.ShareSize)

	testCases := []testCase{
		{
			name:  "tx share",
			share: txShare,
			want:  true,
		},
		{
			name:  "pfb tx share",
			share: pfbTxShare,
			want:  true,
		},
		{
			name:  "blob share",
			share: blobShare,
			want:  false,
		},
	}

	for _, tc := range testCases {
		got := tc.share.IsCompactShare()
		assert.Equal(t, tc.want, got)
	}
}

func TestIsPadding(t *testing.T) {
	type testCase struct {
		name    string
		share   Share
		want    bool
		wantErr bool
	}
	emptyShare := Share{}
	blobShare, _ := zeroPadIfNecessary(
		append(
			namespaceOne,
			[]byte{
				1,          // info byte
				0, 0, 0, 1, // sequence len
				0xff, // data
			}...,
		),
		appconsts.ShareSize)

	testCases := []testCase{
		{
			name:    "empty share",
			share:   emptyShare,
			wantErr: true,
		},
		{
			name:  "blob share",
			share: blobShare,
			want:  false,
		},
		{
			name:  "namespace padding",
			share: NamespacePaddingShare(namespaceOne),
			want:  true,
		},
		{
			name:  "tail padding",
			share: TailPaddingShare(),
			want:  true,
		},
		{
			name:  "reserved padding",
			share: ReservedPaddingShare(),
			want:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.share.IsPadding()
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
