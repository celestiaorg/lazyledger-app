package testnode

import (
	"context"
	"strings"
	"sync"

	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	srvgrpc "github.com/cosmos/cosmos-sdk/server/grpc"
	srvtypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/rpc/client/local"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// mut is used to ensure that only one testnode is running at a time.
// mut was added to prevent simultaneous test executions from impacting each other and causing flaky failures.
var mut sync.Mutex

// StartNode starts the tendermint node along with a local core rpc client. The
// rpc is returned via the client.Context. The function returned should be
// called during cleanup to teardown the node, core client, along with canceling
// the internal context.Context in the returned Context.
func StartNode(tmNode *node.Node, cctx Context) (Context, func() error, error) {
	// locking the mutex not to be able to spin up multiple nodes at the same time.
	mut.Lock()
	if err := tmNode.Start(); err != nil {
		// unlocking the mutex after cleanup finishes.
		mut.Unlock()
		return cctx, func() error { return nil }, err
	}

	coreClient := local.New(tmNode)

	cctx.Context = cctx.WithClient(coreClient)
	goCtx, cancel := context.WithCancel(context.Background())
	cctx.rootCtx = goCtx
	cleanup := func() error {
		// unlocking the mutex after cleanup finishes.
		defer mut.Unlock()
		cancel()
		err := tmNode.Stop()
		if err != nil {
			return err
		}
		tmNode.Wait()
		return nil
	}

	return cctx, cleanup, nil
}

// StartGRPCServer starts the grpc server using the provided application and
// config. A grpc client connection to that server is also added to the client
// context. The returned function should be used to shutdown the server.
func StartGRPCServer(app srvtypes.Application, appCfg *srvconfig.Config, cctx Context) (Context, func() error, error) {
	emptycleanup := func() error { return nil }
	// Add the tx service in the gRPC router.
	app.RegisterTxService(cctx.Context)

	// Add the tendermint queries service in the gRPC router.
	app.RegisterTendermintService(cctx.Context)

	grpcSrv, err := srvgrpc.StartGRPCServer(cctx.Context, app, appCfg.GRPC)
	if err != nil {
		return Context{}, emptycleanup, err
	}

	nodeGRPCAddr := strings.Replace(appCfg.GRPC.Address, "0.0.0.0", "localhost", 1)
	conn, err := grpc.Dial(nodeGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return Context{}, emptycleanup, err
	}

	cctx.Context = cctx.WithGRPCClient(conn)

	return cctx, func() error {
		grpcSrv.Stop()
		return nil
	}, nil
}

// DefaultAppConfig wraps the default config described in the server
func DefaultAppConfig() *srvconfig.Config {
	return srvconfig.DefaultConfig()
}
