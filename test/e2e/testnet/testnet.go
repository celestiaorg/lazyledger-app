//nolint:staticcheck
package testnet

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/celestiaorg/celestia-app/v3/app"
	"github.com/celestiaorg/celestia-app/v3/app/encoding"
	"github.com/celestiaorg/celestia-app/v3/test/util/genesis"
	"github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/celestiaorg/knuu/pkg/preloader"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/rs/zerolog/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type Testnet struct {
	seed      int64
	nodes     []*Node
	genesis   *genesis.Genesis
	keygen    *keyGenerator
	grafana   *GrafanaInfo
	txClients []*TxSim
	knuu      *knuu.Knuu
}

func New(ctx context.Context, name string, seed int64, grafana *GrafanaInfo, chainID string, genesisModifiers ...genesis.Modifier) (*Testnet, error) {
	identifier := fmt.Sprintf("%s_%s", name, time.Now().Format("20060102_150405"))
	kn, err := knuu.New(ctx, knuu.Options{
		Scope:        identifier,
		ProxyEnabled: true,
		// if the tests timeout, pass the timeout option
		// Timeout: 120 * time.Minute,
	})
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("scope", kn.Scope).
		Str("TestName", name).
		Msg("Knuu initialized")

	kn.HandleStopSignal(ctx)

	return &Testnet{
		seed:    seed,
		nodes:   make([]*Node, 0),
		genesis: genesis.NewDefaultGenesis().WithChainID(chainID).WithModifiers(genesisModifiers...),
		keygen:  newKeyGenerator(seed),
		grafana: grafana,
		knuu:    kn,
	}, nil
}

func (t *Testnet) Knuu() *knuu.Knuu {
	return t.knuu
}

func (t *Testnet) NewPreloader() (*preloader.Preloader, error) {
	if t.knuu == nil {
		return nil, errors.New("knuu is not initialized")
	}
	// Since there is one dedicated knuu object for the testnet, each one has its own namespace, and
	// there is one preloader per testnet, can use the same preloader name for all nodes
	return preloader.New("preloader", t.knuu.SystemDependencies)
}

func (t *Testnet) SetConsensusParams(params *tmproto.ConsensusParams) {
	t.genesis.WithConsensusParams(params)
}

func (t *Testnet) SetConsensusMaxBlockSize(size int64) {
	t.genesis.ConsensusParams.Block.MaxBytes = size
}

func (t *Testnet) CreateGenesisNode(ctx context.Context, version string, selfDelegation, upgradeHeightV2 int64, resources Resources, disableBBR bool) error {
	signerKey := t.keygen.Generate(ed25519Type)
	networkKey := t.keygen.Generate(ed25519Type)
	node, err := NewNode(ctx, fmt.Sprintf("val%d", len(t.nodes)), version, 0, selfDelegation, nil, signerKey, networkKey, upgradeHeightV2, resources, t.grafana, t.knuu, disableBBR)
	if err != nil {
		return err
	}
	if err := t.genesis.NewValidator(node.GenesisValidator()); err != nil {
		return err
	}
	t.nodes = append(t.nodes, node)
	return nil
}

func (t *Testnet) CreateGenesisNodes(ctx context.Context, num int, version string, selfDelegation, upgradeHeightV2 int64, resources Resources, disableBBR bool) error {
	for i := 0; i < num; i++ {
		if err := t.CreateGenesisNode(ctx, version, selfDelegation, upgradeHeightV2, resources, disableBBR); err != nil {
			return err
		}
	}
	return nil
}

func (t *Testnet) CreateTxClients(ctx context.Context,
	version string,
	sequences int,
	blobRange string,
	blobPerSequence int,
	resources Resources,
	grpcEndpoints []string,
	upgradeSchedule map[int64]uint64,
) error {
	for i, grpcEndpoint := range grpcEndpoints {
		name := fmt.Sprintf("txsim%d", i)
		err := t.CreateTxClient(ctx, name, version, sequences, blobRange, blobPerSequence, resources, grpcEndpoint, upgradeSchedule)
		if err != nil {
			log.Err(err).Str("name", name).
				Str("grpc endpoint", grpcEndpoint).
				Msg("txsim creation failed")
			return err
		}
		log.Info().
			Str("name", name).
			Str("grpc endpoint", grpcEndpoint).
			Msg("txsim created")
	}
	return nil
}

// CreateTxClient creates a txsim node and sets it up.
func (t *Testnet) CreateTxClient(
	ctx context.Context,
	// name of the txsim knuu instance
	name string,
	// version of the txsim docker image to be pulled from the registry
	version string,
	// sequences: number of sequences to be run by the txsim
	blobSequences int,
	// blobRange: range of blob sizes to be used by the txsim in bytes
	blobRange string,
	blobPerSequence int,
	// resources: resources to be allocated to the txsim
	resources Resources,
	// grpcEndpoint: grpc endpoint of the node to which the txsim will connect and send transactions
	grpcEndpoint string,
	// upgradeSchedule is a map from height to version, specifying the version to upgrade to at the given height
	upgradeSchedule map[int64]uint64,
) error {
	txsimKeyringDir := filepath.Join(os.TempDir(), name)
	defer os.RemoveAll(txsimKeyringDir)
	cdc := encoding.MakeConfig(app.ModuleEncodingRegisters...).Codec
	txsimKeyring, err := keyring.New(app.Name, keyring.BackendTest, txsimKeyringDir, nil, cdc)
	if err != nil {
		return fmt.Errorf("failed to create keyring: %w", err)
	}

	key, _, err := txsimKeyring.NewMnemonic(name, keyring.English, "", "", hd.Secp256k1)
	if err != nil {
		return fmt.Errorf("failed to create mnemonic: %w", err)
	}
	pk, err := key.GetPubKey()
	if err != nil {
		return fmt.Errorf("failed to get public key: %w", err)
	}
	err = t.genesis.AddAccount(genesis.Account{
		PubKey:  pk,
		Balance: 500000,
		Name:    name,
	})
	if err != nil {
		return fmt.Errorf("failed to add account to genesis: %w", err)
	}
	fmt.Printf("before copy")
	records, err := txsimKeyring.List()
	if err != nil {
		return fmt.Errorf("failed to list keys: %w", err)
	}
	for _, record := range records {
		fmt.Printf("kr record name %v\n", record.Name)
	}

	// Copy the keys from the genesis keyring to the txsim keyring
	for _, node := range t.Nodes() {
		nodeName := node.Name
		fmt.Printf("nodeName %v\n", nodeName)
		armor, err := t.Genesis().Keyring().ExportPrivKeyArmor(nodeName, "")
		if err != nil {
			return fmt.Errorf("failed to export key: %w", err)
		}
		err = txsimKeyring.ImportPrivKey(nodeName, armor, "")
		if err != nil {
			return fmt.Errorf("failed to import key: %w", err)
		}
		fmt.Printf("imported %v into txsimKeyring successfully\n", nodeName)
	}

	fmt.Printf("after copy")
	records, err = txsimKeyring.List()
	if err != nil {
		return fmt.Errorf("failed to list keys: %w", err)
	}
	for _, record := range records {
		fmt.Printf("kr record name %v\n", record.Name)
	}

	txsim, err := CreateTxClient(ctx, name, version, grpcEndpoint, t.seed, blobSequences, blobRange, blobPerSequence, 1, resources, txsimKeyringDir, t.knuu, upgradeSchedule)

	err = txsim.Instance.Build().Commit(ctx)
	if err != nil {
		log.Err(err).
			Str("name", name).
			Msg("error committing txsim")
		return err
	}
	t.txClients = append(t.txClients, txsim)
	return nil
}

func (t *Testnet) StartTxClients(ctx context.Context) error {
	for _, txsim := range t.txClients {
		err := txsim.Instance.Execution().StartAsync(ctx)
		if err != nil {
			log.Err(err).
				Str("name", txsim.Name).
				Msg("txsim failed to start")
			return err
		}
		log.Info().
			Str("name", txsim.Name).
			Msg("txsim started")
	}
	// wait for txsims to start
	for _, txsim := range t.txClients {
		err := txsim.Instance.Execution().WaitInstanceIsRunning(ctx)
		if err != nil {
			return fmt.Errorf("txsim %s failed to run: %w", txsim.Name, err)
		}

	}
	return nil
}

func createTxsimKeyringDir(name string, kr keyring.Keyring) string {
	// Create an account, and store it in a temp directory and add the account as genesis account to
	// the testnet.
	txsimKeyringDir := filepath.Join(os.TempDir(), name)
	log.Info().
		Str("name", name).
		Str("directory", txsimKeyringDir).
		Msg("txsim keyring directory created")

	return txsimKeyringDir
}

// CreateAccount creates an account and adds it to the
// testnet genesis. The account is created with the given name and tokens and
// is persisted in the given txsimKeyringDir.
// If txsimKeyringDir is an empty string, an in-memory keyring is created.
func (t *Testnet) CreateAccount(name string, tokens int64, txsimKeyringDir string) (keyring.Keyring, error) {
	cdc := encoding.MakeConfig(app.ModuleEncodingRegisters...).Codec
	var kr keyring.Keyring
	var err error
	// if no keyring directory is specified, create an in-memory keyring
	if txsimKeyringDir == "" {
		kr = keyring.NewInMemory(cdc)
	} else { // create a keyring with the specified directory
		kr, err = keyring.New(app.Name, keyring.BackendTest, txsimKeyringDir, nil, cdc)
		if err != nil {
			return nil, err
		}
	}
	key, _, err := kr.NewMnemonic(name, keyring.English, "", "", hd.Secp256k1)
	if err != nil {
		return nil, err
	}

	pk, err := key.GetPubKey()
	if err != nil {
		return nil, err
	}
	err = t.genesis.AddAccount(genesis.Account{
		PubKey:  pk,
		Balance: tokens,
		Name:    name,
	})
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("name", name).
		Str("pk", pk.String()).
		Msg("txsim account created and added to genesis")
	return kr, nil
}

func (t *Testnet) CreateNode(ctx context.Context, version string, startHeight, upgradeHeight int64, resources Resources, disableBBR bool) error {
	signerKey := t.keygen.Generate(ed25519Type)
	networkKey := t.keygen.Generate(ed25519Type)
	node, err := NewNode(ctx, fmt.Sprintf("val%d", len(t.nodes)), version, startHeight, 0, nil, signerKey, networkKey, upgradeHeight, resources, t.grafana, t.knuu, disableBBR)
	if err != nil {
		return err
	}
	t.nodes = append(t.nodes, node)
	return nil
}

func (t *Testnet) Setup(ctx context.Context, configOpts ...Option) error {
	genesis, err := t.genesis.Export()
	if err != nil {
		return err
	}

	for _, node := range t.nodes {
		// nodes are initialized with the addresses of all other
		// nodes in their addressbook
		peers := make([]string, 0, len(t.nodes)-1)
		for _, peer := range t.nodes {
			if peer.Name != node.Name {
				peers = append(peers, peer.AddressP2P(ctx, true))
			}
		}

		err := node.Init(ctx, genesis, peers, configOpts...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Testnet) RPCEndpoints() []string {
	rpcEndpoints := make([]string, len(t.nodes))
	for idx, node := range t.nodes {
		rpcEndpoints[idx] = node.AddressRPC()
	}
	return rpcEndpoints
}

// FIXME: This does not work currently with the reverse proxy
// func (t *Testnet) GRPCEndpoints() []string {
// 	grpcEndpoints := make([]string, len(t.nodes))
// 	for idx, node := range t.nodes {
// 		grpcEndpoints[idx] = node.AddressGRPC()
// 	}
// 	return grpcEndpoints
// }

// RemoteGRPCEndpoints retrieves the gRPC endpoint addresses of the
// testnet's validator nodes.
func (t *Testnet) RemoteGRPCEndpoints(ctx context.Context) ([]string, error) {
	grpcEndpoints := make([]string, len(t.nodes))
	for idx, node := range t.nodes {
		grpcEP, err := node.RemoteAddressGRPC(ctx)
		if err != nil {
			return nil, err
		}
		grpcEndpoints[idx] = grpcEP
	}
	return grpcEndpoints, nil
}

func (t *Testnet) GetGenesisValidators() []genesis.Validator {
	validators := make([]genesis.Validator, len(t.nodes))
	for i, node := range t.nodes {
		validators[i] = node.GenesisValidator()
	}
	return validators
}

// RemoteRPCEndpoints retrieves the RPC endpoint addresses of the testnet's
// validator nodes.
func (t *Testnet) RemoteRPCEndpoints(ctx context.Context) ([]string, error) {
	rpcEndpoints := make([]string, len(t.nodes))
	for idx, node := range t.nodes {
		grpcEP, err := node.RemoteAddressRPC(ctx)
		if err != nil {
			return nil, err
		}
		rpcEndpoints[idx] = grpcEP
	}
	return rpcEndpoints, nil
}

// WaitToSync waits for the started nodes to sync with the network and move
// past the genesis block.
func (t *Testnet) WaitToSync(ctx context.Context) error {
	genesisNodes := make([]*Node, 0)
	for _, node := range t.nodes {
		if node.StartHeight == 0 {
			genesisNodes = append(genesisNodes, node)
		}
	}
	for _, node := range genesisNodes {
		log.Info().Str("name", node.Name).Msg(
			"waiting for node to sync")
		client, err := node.Client()
		if err != nil {
			return fmt.Errorf("failed to initialize client for node %s: %w", node.Name, err)
		}
		for i := 0; i < 10; i++ {
			resp, err := client.Status(ctx)
			if err == nil {
				if resp.SyncInfo.LatestBlockHeight > 0 {
					log.Info().Int("attempts", i).Str("name", node.Name).Msg(
						"node has synced")
					break
				}
			} else {
				err = errors.New("error getting status")
			}
			if i == 9 {
				return fmt.Errorf("failed to start node %s: %w", node.Name, err)
			}
			log.Info().Str("name", node.Name).Int("attempt", i).Msg(
				"node is not synced yet, waiting...")
			time.Sleep(time.Duration(i) * time.Second)
		}
	}
	return nil
}

// StartNodes starts the testnet nodes and setup proxies.
// It does not wait for the nodes to produce blocks.
// For that, use WaitToSync.
func (t *Testnet) StartNodes(ctx context.Context) error {
	genesisNodes := make([]*Node, 0)
	for _, node := range t.nodes {
		if node.StartHeight == 0 {
			genesisNodes = append(genesisNodes, node)
		}
	}
	// start genesis nodes asynchronously
	for _, node := range genesisNodes {
		err := node.StartAsync(ctx)
		if err != nil {
			return fmt.Errorf("node %s failed to start: %w", node.Name, err)
		}
	}
	log.Info().Msg("create endpoint proxies for genesis nodes")
	// wait for instances to be running
	for _, node := range genesisNodes {
		err := node.WaitUntilStartedAndCreateProxy(ctx)
		if err != nil {
			log.Err(err).Str("name", node.Name).Str("version",
				node.Version).Msg("failed to start and forward ports")
			return fmt.Errorf("node %s failed to start: %w", node.Name, err)
		}
		log.Info().Str("name", node.Name).Str("version",
			node.Version).Msg("started and ports forwarded")
	}
	return nil
}

func (t *Testnet) Start(ctx context.Context) error {
	// start nodes and setup proxies
	err := t.StartNodes(ctx)
	if err != nil {
		return err
	}
	// wait for nodes to sync
	log.Info().Msg("waiting for genesis nodes to sync")
	err = t.WaitToSync(ctx)
	if err != nil {
		return err
	}

	return t.StartTxClients(ctx)
}

func (t *Testnet) Cleanup(ctx context.Context) {
	if err := t.knuu.CleanUp(ctx); err != nil {
		log.Err(err).Msg("failed to cleanup knuu")
	}
}

func (t *Testnet) Node(i int) *Node {
	return t.nodes[i]
}

func (t *Testnet) Nodes() []*Node {
	return t.nodes
}

func (t *Testnet) Genesis() *genesis.Genesis {
	return t.genesis
}
