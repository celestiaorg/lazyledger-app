package app_test

import (
	"fmt"
	"testing"

	"github.com/celestiaorg/celestia-app/v3/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
)

func TestExportAppStateAndValidators(t *testing.T) {
	t.Run("should return exported app for version 1", func(t *testing.T) {
		forZeroHeight := true
		jailAllowedAddrs := []string{}
		testApp, _ := SetupTestAppWithUpgradeHeight(t, 3)

		exported, err := testApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs)
		require.NoError(t, err)
		assert.NotNil(t, exported)
		assert.Equal(t, uint64(1), exported.ConsensusParams.Version.AppVersion)
	})
	t.Run("should return exported app for version 2", func(t *testing.T) {
		forZeroHeight := false
		jailAllowedAddrs := []string{}

		testApp, _ := SetupTestAppWithUpgradeHeight(t, 3)
		upgradeToV2(t, testApp)
		ctx, err := testApp.CreateQueryContext(testApp.LastBlockHeight(), false)
		require.NoError(t, err)
		version := testApp.GetAppVersionFromParamStore(ctx)
		fmt.Println(version)
		exported, err := testApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs)
		require.NoError(t, err)
		assert.NotNil(t, exported)
		assert.Equal(t, uint64(2), exported.ConsensusParams.Version.AppVersion)
	})
}

func upgradeToV2(t *testing.T, testApp *app.App) {
	testApp.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
		Height:  2,
		Version: tmversion.Consensus{App: 1},
	}})
	// Upgrade from v1 -> v2
	testApp.EndBlock(abci.RequestEndBlock{Height: 2})
	testApp.Commit()
	require.EqualValues(t, 2, testApp.AppVersion())
	testApp.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
		Height:  3,
		Version: tmversion.Consensus{App: 2},
	}})
	testApp.EndBlock(abci.RequestEndBlock{Height: 3})
	testApp.Commit()
	require.EqualValues(t, 3, testApp.LastBlockHeight())
	require.Equal(t, uint64(2), testApp.AppVersion())
}
