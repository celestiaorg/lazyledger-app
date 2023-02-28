package wrapper

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/celestiaorg/nmt"
	"github.com/celestiaorg/rsmt2d"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestPushErasuredNamespacedMerkleTree(t *testing.T) {
	testCases := []struct {
		name       string
		squareSize int
	}{
		{"extendedSquareSize = 16", 8},
		{"extendedSquareSize = 256", 128},
	}
	for _, tc := range testCases {
		tc := tc
		tree := NewErasuredNamespacedMerkleTree(uint64(tc.squareSize), 0)

		// push test data to the tree
		for _, d := range generateErasuredData(t, tc.squareSize, appconsts.DefaultCodec()) {
			// push will panic if there's an error
			tree.Push(d)
		}
	}
}

func TestRootErasuredNamespacedMerkleTree(t *testing.T) {
	// check that the root is different from a standard nmt tree this should be
	// the case, because the ErasuredNamespacedMerkleTree should add namespaces
	// to the second half of the tree
	size := 8
	data := generateRandNamespacedRawData(size, appconsts.NamespaceSize, appconsts.ContinuationSparseShareContentSize)
	tree := NewErasuredNamespacedMerkleTree(uint64(size), 0)
	nmtTree := nmt.New(sha256.New())

	for _, d := range data {
		tree.Push(d)
		err := nmtTree.Push(d)
		if err != nil {
			t.Error(err)
		}
	}

	assert.NotEqual(t, nmtTree.Root(), tree.Root())
}

func TestErasureNamespacedMerkleTreePanics(t *testing.T) {
	testCases := []struct {
		name  string
		pFunc assert.PanicTestFunc
	}{
		{
			"push over square size",
			assert.PanicTestFunc(
				func() {
					data := generateErasuredData(t, 16, appconsts.DefaultCodec())
					tree := NewErasuredNamespacedMerkleTree(uint64(15), 0)
					for _, d := range data {
						tree.Push(d)
					}
				}),
		},
		{
			"push in incorrect lexigraphic order",
			assert.PanicTestFunc(
				func() {
					data := generateErasuredData(t, 16, appconsts.DefaultCodec())
					tree := NewErasuredNamespacedMerkleTree(uint64(16), 0)
					for i := len(data) - 1; i > 0; i-- {
						tree.Push(data[i])
					}
				},
			),
		},
	}
	for _, tc := range testCases {
		tc := tc
		assert.Panics(t, tc.pFunc, tc.name)

	}
}

func TestExtendedDataSquare(t *testing.T) {
	squareSize := 4
	// data for a 4X4 square
	raw := generateRandNamespacedRawData(
		squareSize*squareSize,
		appconsts.NamespaceSize,
		appconsts.ContinuationSparseShareContentSize+1, // we +1 here to keep the generated data to be 512 bytes in len
	)

	_, err := rsmt2d.ComputeExtendedDataSquare(raw, appconsts.DefaultCodec(), NewConstructor(uint64(squareSize)))
	assert.NoError(t, err)
}

func TestErasuredNamespacedMerkleTree(t *testing.T) {
	// check that the Tree() returns exact underlying nmt tree
	squareSize := 8
	data := generateRandNamespacedRawData(squareSize, appconsts.NamespaceSize, appconsts.ContinuationSparseShareContentSize)
	tree := NewErasuredNamespacedMerkleTree(uint64(squareSize), 0)

	for _, d := range data {
		tree.Push(d)
	}

	assert.Equal(t, tree.Tree(), tree.tree)
	assert.Equal(t, tree.Tree().Root(), tree.tree.Root())
	assert.Equal(t, appconsts.NamespaceSize, int(tree.Tree().NamespaceSize()))
}

// generateErasuredData produces a slice that is twice as long as it erasures
// the data
func generateErasuredData(t *testing.T, numLeaves int, codec rsmt2d.Codec) [][]byte {
	raw := generateRandNamespacedRawData(
		numLeaves,
		appconsts.NamespaceSize,
		appconsts.ContinuationSparseShareContentSize+1, // we +1 here to keep the generated data to be 512 bytes in len
	)
	erasuredData, err := codec.Encode(raw)
	if err != nil {
		t.Error(err)
	}
	return append(raw, erasuredData...)
}

// this code is copy pasted from the plugin, and should likely be exported in the plugin instead
func generateRandNamespacedRawData(total int, nidSize int, leafSize int) [][]byte {
	data := make([][]byte, total)
	for i := 0; i < total; i++ {
		nid := make([]byte, nidSize)
		_, err := rand.Read(nid)
		if err != nil {
			panic(err)
		}
		data[i] = nid
	}

	sortByteArrays(data)
	for i := 0; i < total; i++ {
		d := make([]byte, leafSize)
		_, err := rand.Read(d)
		if err != nil {
			panic(err)
		}
		data[i] = append(data[i], d...)
	}

	return data
}

func sortByteArrays(src [][]byte) {
	sort.Slice(src, func(i, j int) bool { return bytes.Compare(src[i], src[j]) < 0 })
}

func TestErasuredNamespacedMerkleTree_Prove(t *testing.T) {
	nidSizes := []int{8, 16, 20, 32}
	for _, nidSize := range nidSizes {
		data := generateRandNamespacedRawData(appconsts.DefaultMaxSquareSize, nidSize, appconsts.ShareSize-nidSize)
		tree := NewErasuredNamespacedMerkleTree(appconsts.DefaultMaxSquareSize, 0, nmt.NamespaceIDSize(nidSize))
		for _, d := range data {
			tree.Push(d)
		}
		proof, err := tree.Prove(0)
		assert.NoError(t, err)
		fmt.Printf("nidSize=%v\n", nidSize)
		fmt.Printf("unencoded proof is %v bytes\n", proofSize(proof))

		publicProof := PublicProof{
			Start:                    int64(proof.Start()),
			End:                      int64(proof.End()),
			Nodes:                    proof.Nodes(),
			LeafHash:                 proof.LeafHash(),
			IsMaxNamespace_IDIgnored: proof.IsMaxNamespaceIDIgnored(),
		}
		b, err := proto.Marshal(&publicProof)
		assert.NoError(t, err)

		filename := fmt.Sprintf("/tmp/proof_%v.bin", nidSize)
		err = os.WriteFile(filename, b, 0644)

		f, err := os.Open(filename)
		assert.NoError(t, err)
		read := bufio.NewReader(f)
		d, err := ioutil.ReadAll(read)
		assert.NoError(t, err)
		compressedFilename := strings.Replace(filename, ".bin", ".gz", -1)
		f, err = os.Create(compressedFilename)
		assert.NoError(t, err)
		w := gzip.NewWriter(f)
		w.Write(d)
		w.Close()

		uncompressedFile, err := os.Stat(filename)
		assert.NoError(t, err)
		compressedFile, err := os.Stat(compressedFilename)
		assert.NoError(t, err)
		fmt.Printf("%v is %v bytes\n", filename, uncompressedFile.Size())
		fmt.Printf("%v is %v bytes\n", compressedFilename, compressedFile.Size())

		fmt.Println()
	}
}

func proofSize(proof nmt.Proof) int {
	size := 0
	for _, node := range proof.Nodes() {
		size += len(node)
	}
	return size
}

// type PublicProof struct {
// 	Start                   int
// 	End                     int
// 	Nodes                   [][]byte
// 	LeafHash                []byte
// 	IsMaxNamespaceIDIgnored bool
// }
