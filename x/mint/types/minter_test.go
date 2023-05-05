package types

import (
	"math/rand"
	"testing"
	time "time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestCalculateInflationRate(t *testing.T) {
	minter := DefaultMinter()
	genesisTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	minter.GenesisTime = &genesisTime

	type testCase struct {
		year int
		want float64
	}

	testCases := []testCase{
		{0, 0.08},
		{1, 0.072},
		{2, 0.0648},
		{3, 0.05832},
		{4, 0.052488},
		{5, 0.0472392},
		{6, 0.04251528},
		{7, 0.038263752},
		{8, 0.0344373768},
		{9, 0.03099363912},
		{10, 0.027894275208},
		{11, 0.0251048476872},
		{12, 0.02259436291848},
		{13, 0.020334926626632},
		{14, 0.0183014339639688},
		{15, 0.01647129056757192},
		{16, 0.0150},
		{17, 0.0150},
		{18, 0.0150},
		{19, 0.0150},
		{20, 0.0150},
		{21, 0.0150},
		{22, 0.0150},
		{23, 0.0150},
		{24, 0.0150},
		{25, 0.0150},
		{26, 0.0150},
		{27, 0.0150},
		{28, 0.0150},
		{29, 0.0150},
		{30, 0.0150},
		{31, 0.0150},
		{32, 0.0150},
		{33, 0.0150},
		{34, 0.0150},
		{35, 0.0150},
		{36, 0.0150},
		{37, 0.0150},
		{38, 0.0150},
		{39, 0.0150},
		{40, 0.0150},
	}

	for _, tc := range testCases {
		blockTime := genesisTime.AddDate(tc.year, 0, 0)
		ctx := sdk.NewContext(nil, tmproto.Header{}, false, nil).WithBlockTime(blockTime)
		inflationRate := minter.CalculateInflationRate(ctx)
		got, err := inflationRate.Float64()
		assert.NoError(t, err)
		assert.Equal(t, tc.want, got, "want %v got %v year %v blockTime %v", tc.want, got, tc.year, blockTime)
	}
}

func TestCalculateBlockProvision(t *testing.T) {
	minter := DefaultMinter()

	type testCase struct {
		annualProvisions int64
		want             sdk.Coin
	}
	testCases := []testCase{
		{
			annualProvisions: BlocksPerYear,
			want:             sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1)),
		},
		{
			annualProvisions: BlocksPerYear * 2,
			want:             sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(2)),
		},
		{
			annualProvisions: (BlocksPerYear * 10) - 1,
			want:             sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(9)),
		},
		{
			annualProvisions: BlocksPerYear / 2,
			want:             sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)),
		},
	}
	for _, tc := range testCases {
		minter.AnnualProvisions = sdk.NewDec(tc.annualProvisions)
		got := minter.CalculateBlockProvision()
		require.True(t, tc.want.IsEqual(got), "want %v got %v", tc.want, got)
	}
}

func BenchmarkCalculateBlockProvision(b *testing.B) {
	b.ReportAllocs()
	minter := DefaultMinter()

	s1 := rand.NewSource(100)
	r1 := rand.New(s1)
	minter.AnnualProvisions = sdk.NewDec(r1.Int63n(1000000))

	for n := 0; n < b.N; n++ {
		minter.CalculateBlockProvision()
	}
}

func BenchmarkCalculateInflationRate(b *testing.B) {
	b.ReportAllocs()
	minter := DefaultMinter()

	for n := 0; n < b.N; n++ {
		ctx := sdk.NewContext(nil, tmproto.Header{Height: int64(n)}, false, nil)
		minter.CalculateInflationRate(ctx)
	}
}

func BenchmarkCalculateAnnualProvisions(b *testing.B) {
	b.ReportAllocs()
	minter := DefaultMinter()
	totalSupply := sdk.NewInt(100000000000000)

	for n := 0; n < b.N; n++ {
		minter.CalculateAnnualProvisions(totalSupply)
	}
}

func Test_yearsSinceGenesis(t *testing.T) {
	type testCase struct {
		name    string
		current time.Time
		want    uint64
	}

	genesis := time.Date(2023, 1, 1, 12, 30, 15, 0, time.UTC) // 2023-01-01T12:30:15Z
	oneDay, err := time.ParseDuration("24h")
	assert.NoError(t, err)
	oneWeek := oneDay * 7
	oneMonth := oneDay * 30
	oneYear := oneDay * 365
	twoYears := oneYear * 2
	tenYears := oneYear * 10
	tenYearsOneMonth := oneYear*10 + oneMonth

	testCases := []testCase{
		{
			name:    "one day after genesis",
			current: genesis.Add(oneDay),
			want:    0,
		},
		{
			name:    "one day before genesis",
			current: genesis.Add(-oneDay),
			want:    0,
		},
		{
			name:    "one week after genesis",
			current: genesis.Add(oneWeek),
			want:    0,
		},
		{
			name:    "one month after genesis",
			current: genesis.Add(oneMonth),
			want:    0,
		},
		{
			name:    "one year after genesis",
			current: genesis.Add(oneYear),
			want:    1,
		},
		{
			name:    "two years after genesis",
			current: genesis.Add(twoYears),
			want:    2,
		},
		{
			name:    "ten years after genesis",
			current: genesis.Add(tenYears),
			want:    10,
		},
		{
			name:    "ten years and one month after genesis",
			current: genesis.Add(tenYearsOneMonth),
			want:    10,
		},
	}

	for _, tc := range testCases {
		got := yearsSinceGenesis(genesis, tc.current)
		assert.Equal(t, tc.want, got, tc.name)
	}
}
