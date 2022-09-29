package shares

import (
	"bytes"
	"reflect"
	"testing"

	coretypes "github.com/tendermint/tendermint/types"
)

func TestSplitTxs(t *testing.T) {
	type testCase struct {
		name string
		txs  coretypes.Txs
		want [][]byte
	}
	testCases := []testCase{
		{
			name: "empty txs",
			txs:  coretypes.Txs{},
			want: [][]byte{},
		},
		{
			name: "one small tx",
			txs:  coretypes.Txs{coretypes.Tx{0xa}},
			want: [][]uint8{
				append([]uint8{
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, // namespace id
<<<<<<< HEAD
					0x1,                // info byte
					0x2, 0x0, 0x0, 0x0, // 1 byte (unit)  + 1 byte (unit length) = 2 bytes message length
					0x0, // BUG: reserved byte should be non-zero
					0x1, // unit length of first transaction
					0xa, // data of first transaction
				}, bytes.Repeat([]byte{0}, 240)...), // padding
=======
					0x1, // info byte
					0x0, // reserved byte
					0x1, // unit length of first transaction
					0xa, // data of first transaction
				}, bytes.Repeat([]byte{0}, 244)...), // padding
>>>>>>> 5dc1fb8 (test: transaction splitting (#813))
			},
		},
		{
			name: "two small txs",
			txs:  coretypes.Txs{coretypes.Tx{0xa}, coretypes.Tx{0xb}},
			want: [][]uint8{
				append([]uint8{
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, // namespace id
<<<<<<< HEAD
					0x1,                // info byte
					0x4, 0x0, 0x0, 0x0, // 2 bytes (first transaction) + 2 bytes (second transaction) = 4 bytes message length
					0x0, // BUG: reserved byte should be non-zero
=======
					0x1, // info byte
					0x0, // reserved byte
>>>>>>> 5dc1fb8 (test: transaction splitting (#813))
					0x1, // unit length of first transaction
					0xa, // data of first transaction
					0x1, // unit length of second transaction
					0xb, // data of second transaction
<<<<<<< HEAD
				}, bytes.Repeat([]byte{0}, 238)...), // padding
=======
				}, bytes.Repeat([]byte{0}, 242)...), // padding
>>>>>>> 5dc1fb8 (test: transaction splitting (#813))
			},
		},
		{
			name: "one large tx that spans two shares",
<<<<<<< HEAD
			txs:  coretypes.Txs{bytes.Repeat([]byte{0xC}, 241)},
			want: [][]uint8{
				append([]uint8{
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, // namespace id
					0x1,          // info byte
					243, 1, 0, 0, // 241 (unit) + 2 (unit length) = 243 message length
					0x0,    // BUG: reserved byte should be non-zero
					241, 1, // unit length of first transaction is 241
				}, bytes.Repeat([]byte{0xc}, 240)...), // data of first transaction
=======
			txs:  coretypes.Txs{bytes.Repeat([]byte{0xC}, 245)},
			want: [][]uint8{
				append([]uint8{
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, // namespace id
					0x1,    // info byte
					0x0,    // BUG reserved byte should be non-zero see https://github.com/celestiaorg/celestia-app/issues/802
					245, 1, // unit length of first transaction is 245
				}, bytes.Repeat([]byte{0xc}, 244)...), // data of first transaction
>>>>>>> 5dc1fb8 (test: transaction splitting (#813))
				append([]uint8{
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, // namespace id
					0x0, // info byte
					0x0, // reserved byte
					0xc, // continuation data of first transaction
				}, bytes.Repeat([]byte{0x0}, 245)...), // padding
			},
		},
		{
			name: "one small tx then one large tx that spans two shares",
<<<<<<< HEAD
			txs:  coretypes.Txs{coretypes.Tx{0xd}, bytes.Repeat([]byte{0xe}, 241)},
			want: [][]uint8{
				append([]uint8{
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, // namespace id
					0x1,          // info byte
					245, 1, 0, 0, // 2 bytes (first transaction) + 243 bytes (second transaction) = 245 bytes message length
					0x0,    // BUG: reserved byte should be non-zero
					1,      // unit length of first transaction
					0xd,    // data of first transaction
					241, 1, // unit length of second transaction is 241
				}, bytes.Repeat([]byte{0xe}, 238)...), // data of first transaction
				append([]uint8{
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, // namespace id
					0x0,           // info byte
					0x0,           // reserved byte
					0xe, 0xe, 0xe, // continuation data of second transaction
				}, bytes.Repeat([]byte{0x0}, 243)...), // padding
=======
			txs:  coretypes.Txs{coretypes.Tx{0xd}, bytes.Repeat([]byte{0xe}, 243)},
			want: [][]uint8{
				append([]uint8{
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, // namespace id
					0x1,    // info byte
					0x0,    // BUG reserved byte should be non-zero see https://github.com/celestiaorg/celestia-app/issues/802
					1,      // unit length of first transaction
					0xd,    // data of first transaction
					243, 1, // unit length of second transaction is 243
				}, bytes.Repeat([]byte{0xe}, 242)...), // data of first transaction
				append([]uint8{
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, // namespace id
					0x0, // info byte
					0x0, // reserved byte
					0xe, // continuation data of second transaction
				}, bytes.Repeat([]byte{0x0}, 245)...), // padding
>>>>>>> 5dc1fb8 (test: transaction splitting (#813))
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitTxs(tt.txs)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitTxs(%#v) got %#v, want %#v", tt.txs, got, tt.want)
			}
		})
	}
}
