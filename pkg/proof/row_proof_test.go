package proof

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRowProofValidate(t *testing.T) {
	type testCase struct {
		name    string
		rp      RowProof
		root    []byte
		wantErr bool
	}
	testCases := []testCase{
		{
			name:    "empty row proof returns error",
			rp:      RowProof{},
			root:    root,
			wantErr: true,
		},
		{
			name:    "row proof with mismatched number of rows and row roots returns error",
			rp:      mismatchedRowRoots(),
			root:    root,
			wantErr: true,
		},
		{
			name:    "row proof with mismatched number of proofs returns error",
			rp:      mismatchedProofs(),
			root:    root,
			wantErr: true,
		},
		{
			name:    "row proof with mismatched number of rows returns error",
			rp:      mismatchedRows(),
			root:    root,
			wantErr: true,
		},
		{
			name:    "valid row proof returns no error",
			rp:      validRowProof(),
			root:    root,
			wantErr: false,
		},
		{
			name:    "valid row proof with incorrect root returns error",
			rp:      validRowProof(),
			root:    incorrectRoot,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.rp.Validate(tc.root)
			if tc.wantErr {
				assert.Error(t, got)
				return
			}
			assert.NoError(t, got)
		})
	}
}

// root is the root hash of the Merkle tree used in validRowProof
var root = []byte{0x82, 0x37, 0x91, 0xd2, 0x5d, 0x77, 0x7, 0x67, 0x35, 0x3, 0x90, 0x12, 0x10, 0xc4, 0x43, 0x8a, 0x8b, 0x78, 0x4b, 0xbf, 0x5b, 0x8f, 0xa6, 0x40, 0xa9, 0x51, 0xa7, 0xa9, 0xbd, 0x52, 0xd5, 0xf6}

var incorrectRoot = bytes.Repeat([]byte{0}, 32)

// validRowProof returns a row proof for one row. This test data copied from
// ceelestia-app's pkg/proof/proof_test.go TestNewShareInclusionProof: "1
// transaction share"
func validRowProof() RowProof {
	return RowProof{
		RowRoots: [][]byte{{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x9d, 0xe6, 0x38, 0x91, 0xc1, 0x6, 0xaf, 0x81, 0x75, 0x5a, 0x36, 0xf5, 0xb2, 0x62, 0x1e, 0xfa, 0xb9, 0xb8, 0x73, 0x87, 0xef, 0xe3, 0x6b, 0x33, 0xd8, 0xbf, 0xc9, 0x87, 0x1b, 0x8d, 0xfa, 0x8a}},
		Proofs: []*Proof{
			{
				Total:    128,
				Index:    0,
				LeafHash: []uint8{0x0, 0xcc, 0xfb, 0xff, 0x62, 0x10, 0x71, 0x61, 0x2f, 0xb9, 0x5a, 0xb1, 0xc3, 0x83, 0xff, 0x1d, 0x30, 0x31, 0x86, 0x42, 0xe4, 0x8e, 0x59, 0xe8, 0x8b, 0x92, 0x83, 0x11, 0x67, 0xb, 0xfc, 0x9a},
				Aunts:    [][]uint8{{0x5c, 0xc6, 0x3b, 0x1e, 0x91, 0xa4, 0xbf, 0x6a, 0xa7, 0xd2, 0x68, 0x1c, 0x44, 0xc1, 0xda, 0xa2, 0x22, 0xed, 0x33, 0xb8, 0xd0, 0x29, 0x48, 0xfc, 0xab, 0x8f, 0x71, 0x50, 0x9c, 0xbb, 0x15, 0xab}, {0xc6, 0x14, 0x2b, 0x33, 0x5d, 0xaa, 0xfa, 0x20, 0xdf, 0x8a, 0x9b, 0xe9, 0x29, 0x9b, 0x34, 0xcd, 0xeb, 0xe7, 0x35, 0x39, 0x5c, 0x58, 0xb1, 0x13, 0x1f, 0x4, 0xeb, 0xdc, 0x33, 0x99, 0xdf, 0x98}, {0xdb, 0x99, 0xe2, 0xdf, 0x86, 0x84, 0x24, 0x90, 0x44, 0x8e, 0x29, 0x26, 0xe1, 0xb2, 0xb0, 0x52, 0x42, 0xf9, 0x73, 0x7, 0x7f, 0xab, 0x1d, 0xa9, 0xad, 0x56, 0x10, 0xf0, 0x58, 0xdf, 0x8, 0xd7}, {0x48, 0xfd, 0xfc, 0x3b, 0x96, 0xa5, 0x19, 0xf5, 0x14, 0xf, 0x37, 0xfd, 0x95, 0xb3, 0x76, 0xfb, 0x7e, 0x5, 0x5b, 0x4d, 0x8b, 0x68, 0x16, 0x81, 0x51, 0x92, 0x44, 0x0, 0xe5, 0xf6, 0x49, 0x16}, {0xfb, 0x45, 0xdc, 0x2, 0x8b, 0xa9, 0x45, 0xfe, 0xa0, 0x7b, 0xeb, 0x62, 0x81, 0x84, 0x95, 0x19, 0x29, 0xf5, 0x78, 0x16, 0x15, 0xb8, 0xf2, 0xa3, 0x94, 0x96, 0xb1, 0x4c, 0x4c, 0xef, 0xf4, 0xd3}, {0x2c, 0x26, 0x82, 0xb1, 0x8c, 0x9f, 0xff, 0x50, 0xde, 0x67, 0x4e, 0x82, 0x3, 0x3, 0xd6, 0xdc, 0x7c, 0x7a, 0xea, 0x1a, 0xe3, 0x9, 0xf0, 0x1a, 0xc6, 0xcd, 0x19, 0x34, 0xc7, 0x54, 0x6, 0x14}, {0xe9, 0x41, 0x8b, 0x1, 0x9a, 0xd6, 0xd3, 0x13, 0x21, 0x14, 0x89, 0x98, 0xbb, 0x81, 0xda, 0xf7, 0xa, 0x36, 0x14, 0xcf, 0xc5, 0xac, 0xbf, 0xc3, 0x48, 0xb0, 0x88, 0x90, 0x45, 0x29, 0x80, 0x23}},
			},
		},
		StartRow: 0,
		EndRow:   0,
	}
}

func mismatchedRowRoots() RowProof {
	rp := validRowProof()
	rp.RowRoots = [][]byte{}
	return rp
}

func mismatchedProofs() RowProof {
	rp := validRowProof()
	rp.Proofs = []*Proof{}
	return rp
}

func mismatchedRows() RowProof {
	rp := validRowProof()
	rp.EndRow = 10
	return rp
}
