package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMinter returns a new Minter object.
func NewMinter(inflationRate sdk.Dec, annualProvisions sdk.Dec) Minter {
	return Minter{
		InflationRate:    inflationRate,
		AnnualProvisions: annualProvisions,
	}
}

// DefaultMinter returns a Minter object with default values.
func DefaultMinter() Minter {
	return NewMinter(initalInflationRate, sdk.NewDec(0))
}

// ValidateMinter returns an error if the provided minter is invalid.
func ValidateMinter(minter Minter) error {
	if minter.InflationRate.IsNegative() {
		return fmt.Errorf("minter inflation %v should be positive", minter.InflationRate.String())
	}
	if minter.AnnualProvisions.IsNegative() {
		return fmt.Errorf("minter annual provisions %v should be positive", minter.AnnualProvisions.String())
	}
	return nil
}

// CalculateInflationRate returns the inflation rate for the current year depending on
// the current block height in context. The inflation rate is expected to
// decrease every year according to the schedule specified in the README.
func (m Minter) CalculateInflationRate(ctx sdk.Context) sdk.Dec {
	year := uint64(ctx.BlockHeader().Height) / BlocksPerYear
	inflationRate := initalInflationRate.Mul(sdk.OneDec().Sub(disinflationRate).Power(year))

	if inflationRate.LT(targetInflationRate) {
		return targetInflationRate
	}
	return inflationRate
}

// CalculateAnnualProvisions returns the annual provisions (i.e. the total
// number of tokens that should be minted due to inflation for the current
// year).
func (m Minter) CalculateAnnualProvisions(totalSupply math.Int) sdk.Dec {
	return m.InflationRate.MulInt(totalSupply)
}

// CalculateBlockProvision returns the provisions for a block (i.e. the total number of
// coins that should be minted due to inflation for the current block).
func (m Minter) CalculateBlockProvision() sdk.Coin {
	blockProvision := m.AnnualProvisions.QuoInt(blocksPerYear)
	return sdk.NewCoin(sdk.DefaultBondDenom, blockProvision.TruncateInt())
}
