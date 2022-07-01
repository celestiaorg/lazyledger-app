package orchestrator

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	paytypes "github.com/celestiaorg/celestia-app/x/payment/types"
	"github.com/celestiaorg/celestia-app/x/qgb/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktypestx "github.com/cosmos/cosmos-sdk/types/tx"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	tmlog "github.com/tendermint/tendermint/libs/log"
	corerpctypes "github.com/tendermint/tendermint/rpc/core/types"
	coretypes "github.com/tendermint/tendermint/types"
	"google.golang.org/grpc"
	"math/big"
	"strconv"
	"sync"
	"time"
)

type Orchestrator struct {
	logger tmlog.Logger // maybe use a more general interface

	evmPrivateKey  ecdsa.PrivateKey
	signer         *paytypes.KeyringSigner
	orchEthAddress stakingtypes.EthAddress

	querier     Querier
	broadcaster Broadcaster
	retrier     Retrier
}

func NewOrchestrator(
	logger tmlog.Logger,
	querier Querier,
	broadcaster Broadcaster,
	retrier Retrier,
	signer *paytypes.KeyringSigner,
	evmPrivateKey ecdsa.PrivateKey,
) *Orchestrator {
	orchEthAddr, err := stakingtypes.NewEthAddress(crypto.PubkeyToAddress(evmPrivateKey.PublicKey).Hex())
	if err != nil {
		panic(err)
	}

	return &Orchestrator{
		logger:         logger,
		signer:         signer,
		evmPrivateKey:  evmPrivateKey,
		orchEthAddress: *orchEthAddr,
		querier:        querier,
		broadcaster:    broadcaster,
		retrier:        retrier,
	}
}

func (orch Orchestrator) Start(ctx context.Context) {
	// contains the nonces that will be signed by the orchestrator.
	noncesQueue := make(chan uint64, 100)

	// used to send a signal when the nonces processor wants to notify the nonces enqueuing services to stop.
	signalChan := make(chan struct{})
	defer close(noncesQueue)

	withCancel, cancel := context.WithCancel(ctx)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := orch.enqueueMissingEvents(withCancel, noncesQueue, signalChan)
		if err != nil {
			orch.logger.Error("error enqueing missing attestations", "err", err)
			cancel()
		}
		orch.logger.Error("stopping enqueing missing attestations")
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := orch.startNewEventsListener(withCancel, noncesQueue, signalChan)
		if err != nil {
			orch.logger.Error("error listening to new attestations", "err", err)
			cancel()
		}
		orch.logger.Error("stopping listening to new attestations")
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := orch.processNonces(withCancel, noncesQueue, signalChan)
		if err != nil {
			orch.logger.Error("error processing attestations", "err", err)
			cancel()
		}
		orch.logger.Error("stopping processing attestations")
		wg.Done()
	}()

	// FIXME should we add  another go routine that keep checking if all the attestations
	// were signed every 10min for example?

	select {
	case <-withCancel.Done():
		wg.Wait()
		return
	}
}

func (orch Orchestrator) startNewEventsListener(ctx context.Context, queue chan<- uint64, signalChan <-chan struct{}) error {
	results, err := orch.querier.SubscribeEvents(ctx, "attestation-changes", fmt.Sprintf("%s.%s='%s'", types.EventTypeAttestationRequest, sdk.AttributeKeyModule, types.ModuleName))
	if err != nil {
		panic(err)
	}
	attestationEventName := fmt.Sprintf("%s.%s", types.EventTypeAttestationRequest, types.AttributeKeyNonce)
	orch.logger.Info("listening for new block events...")
	for {
		select {
		case <-signalChan:
			return nil
		case <-ctx.Done():
			return nil
		case result := <-results:
			blockEvent := mustGetEvent(result, coretypes.EventTypeKey)
			isBlock := blockEvent[0] == coretypes.EventNewBlock
			if !isBlock {
				// we only want to handle the attestation when the block is committed
				continue
			}
			attestationEvent := mustGetEvent(result, attestationEventName)
			nonce, err := strconv.Atoi(attestationEvent[0])
			if err != nil {
				return err
			}
			orch.logger.Debug("enqueueing new attestation nonce", "nonce", nonce)
			select {
			case <-signalChan:
				return nil
			case queue <- uint64(nonce):
			}
		}
	}
}

func (orch Orchestrator) enqueueMissingEvents(ctx context.Context, queue chan<- uint64, signalChan <-chan struct{}) error {
	latestNonce, err := orch.querier.QueryLatestAttestationNonce(ctx)
	if err != nil {
		return err
	}

	lastUnbondingHeight, err := orch.querier.QueryLastUnbondingHeight(ctx)
	if err != nil {
		return err
	}

	orch.logger.Info("syncing missing nonces", "latest_nonce", latestNonce, "last_unbonding_height", lastUnbondingHeight)

	for i := lastUnbondingHeight; i < latestNonce; i++ {
		select {
		case <-signalChan:
			return nil
		case <-ctx.Done():
			return nil
		default:
			orch.logger.Debug("enqueueing missing attestation nonce", "nonce", latestNonce-i)
			select {
			case <-signalChan:
				return nil
			case queue <- latestNonce - i:
			}
		}
	}
	orch.logger.Info("finished syncing missing nonces", "latest_nonce", latestNonce, "last_unbonding_height", lastUnbondingHeight)
	return nil
}

func (orch Orchestrator) processNonces(ctx context.Context, noncesQueue <-chan uint64, signalChan chan<- struct{}) error {
	for {
		select {
		case <-ctx.Done():
			close(signalChan)
			return nil
		case nonce := <-noncesQueue:
			orch.logger.Debug("processing nonce", "nonce", nonce)
			if err := orch.Process(ctx, nonce); err != nil {
				orch.logger.Error("failed to process nonce, retrying...", "nonce", nonce, "err", err)
				if err := orch.retrier.Retry(nonce, ctx, orch.Process); err != nil {
					close(signalChan)
					return err
				}
			}
		}
	}
}

func (orch Orchestrator) Process(ctx context.Context, nonce uint64) error {
	att, err := orch.querier.QueryAttestationByNonce(ctx, nonce)
	if err != nil {
		return err
	}
	switch att.Type() {
	case types.ValsetRequestType:
		vs, ok := att.(*types.Valset)
		if !ok {
			return errors.Wrap(types.ErrAttestationNotValsetRequest, strconv.FormatUint(nonce, 10))
		}
		resp, err := orch.querier.QueryValsetConfirm(ctx, nonce, orch.signer.GetSignerInfo().GetAddress().String())
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("valset %d", nonce))
		}
		if resp != nil {
			orch.logger.Debug("already signed valset", "nonce", nonce, "signature", resp.Signature)
			return nil
		}
		err = orch.processValsetEvent(ctx, *vs)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("valset %d", nonce))
		}
		return nil

	case types.DataCommitmentRequestType:
		dc, ok := att.(*types.DataCommitment)
		if !ok {
			return errors.Wrap(types.ErrAttestationNotDataCommitmentRequest, strconv.FormatUint(nonce, 10))
		}
		resp, err := orch.querier.QueryDataCommitmentConfirm(ctx, dc.EndBlock, dc.BeginBlock, orch.signer.GetSignerInfo().GetAddress().String())
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("data commitment %d", nonce))
		}
		if resp != nil {
			orch.logger.Debug("already signed data commitment", "nonce", nonce, "begin_block", resp.BeginBlock, "end_block", resp.EndBlock, "commitment", resp.Commitment, "signature", resp.Signature)
			return nil
		}
		err = orch.processDataCommitmentEvent(ctx, *dc)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("data commitment %d", nonce))
		}
		return nil

	default:
		return errors.Wrap(ErrUnknownAttestationType, strconv.FormatUint(nonce, 10))
	}
}

func (orch Orchestrator) processValsetEvent(ctx context.Context, valset types.Valset) error {
	signBytes, err := valset.SignBytes(types.BridgeId)
	if err != nil {
		return err
	}
	signature, err := types.NewEthereumSignature(signBytes.Bytes(), &orch.evmPrivateKey)
	if err != nil {
		return err
	}

	// create and send the valset hash
	msg := types.NewMsgValsetConfirm(
		valset.Nonce,
		orch.orchEthAddress,
		orch.signer.GetSignerInfo().GetAddress(),
		ethcmn.Bytes2Hex(signature),
	)
	hash, err := orch.broadcaster.BroadcastTx(ctx, msg)
	if err != nil {
		return err
	}
	orch.logger.Info("signed Valset", "nonce", msg.Nonce, "tx_hash", hash)
	return nil
}

func (orch Orchestrator) processDataCommitmentEvent(
	ctx context.Context,
	dc types.DataCommitment,
) error {
	commitment, err := orch.querier.QueryCommitment(
		ctx,
		fmt.Sprintf("block.height >= %d AND block.height <= %d",
			dc.BeginBlock,
			dc.EndBlock,
		),
	)
	if err != nil {
		return err
	}
	dataRootHash := types.DataCommitmentTupleRootSignBytes(types.BridgeId, big.NewInt(int64(dc.Nonce)), commitment)
	dcSig, err := types.NewEthereumSignature(dataRootHash.Bytes(), &orch.evmPrivateKey)
	if err != nil {
		return err
	}

	msg := types.NewMsgDataCommitmentConfirm(
		commitment.String(),
		ethcmn.Bytes2Hex(dcSig),
		orch.signer.GetSignerInfo().GetAddress(),
		orch.orchEthAddress,
		dc.BeginBlock,
		dc.EndBlock,
		dc.Nonce,
	)
	hash, err := orch.broadcaster.BroadcastTx(ctx, msg)
	if err != nil {
		return err
	}
	orch.logger.Info("signed commitment", "nonce", msg.Nonce, "begin_block", msg.BeginBlock, "end_block", msg.EndBlock, "commitment", commitment, "tx_hash", hash)
	return nil
}

var _ Broadcaster = &broadcaster{}

type Broadcaster interface {
	BroadcastTx(ctx context.Context, msg sdk.Msg) (string, error)
}

type broadcaster struct {
	mutex   *sync.Mutex
	signer  *paytypes.KeyringSigner
	qgbGrpc *grpc.ClientConn
}

func NewBroadcaster(qgbGrpcAddr string, signer *paytypes.KeyringSigner) (Broadcaster, error) {
	qgbGrpc, err := grpc.Dial(qgbGrpcAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &broadcaster{
		mutex:   &sync.Mutex{}, // investigate if this is needed
		signer:  signer,
		qgbGrpc: qgbGrpc,
	}, nil
}

func (bc *broadcaster) BroadcastTx(ctx context.Context, msg sdk.Msg) (string, error) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	err := bc.signer.QueryAccountNumber(ctx, bc.qgbGrpc)
	if err != nil {
		return "", err
	}

	builder := bc.signer.NewTxBuilder()
	// TODO make gas limit configurable
	builder.SetGasLimit(9999999999999)
	// TODO: update this api
	// via https://github.com/celestiaorg/celestia-app/pull/187/commits/37f96d9af30011736a3e6048bbb35bad6f5b795c
	tx, err := bc.signer.BuildSignedTx(builder, msg)
	if err != nil {
		return "", err
	}

	rawTx, err := bc.signer.EncodeTx(tx)
	if err != nil {
		return "", err
	}

	// FIXME sdktypestx.BroadcastMode_BROADCAST_MODE_BLOCK waits for a block to be minted containing
	// the transaction to continue. This makes the orchestrator slow to catchup.
	// It would be better to just send the transaction. Then, another job would keep an eye
	// if the transaction was included. If not, retry it. But this would mean we should increment ourselves
	// the sequence number after each broadcasted transaction.
	// TODO  check if we can move this outside of the paytypes
	resp, err := paytypes.BroadcastTx(ctx, bc.qgbGrpc, sdktypestx.BroadcastMode_BROADCAST_MODE_BLOCK, rawTx)
	if err != nil {
		return "", err
	}

	if resp.TxResponse.Code != 0 {
		return "", errors.Wrap(ErrFailedBroadcast, resp.TxResponse.RawLog)
	}

	return resp.TxResponse.TxHash, nil
}

var _ Retrier = &retrier{}

type retrier struct {
	logger        tmlog.Logger
	retriesNumber int
}

func NewRetrier(logger tmlog.Logger, retriesNumber int) *retrier {
	return &retrier{
		logger:        logger,
		retriesNumber: retriesNumber,
	}
}

type Retrier interface {
	Retry(nonce uint64, ctx context.Context, retryMethod func(context.Context, uint64) error) error
	RetryThenFail(nonce uint64, ctx context.Context, retryMethod func(context.Context, uint64) error)
}

func (r retrier) Retry(nonce uint64, ctx context.Context, retryMethod func(context.Context, uint64) error) error {
	var err error
	for i := 0; i <= r.retriesNumber; i++ {
		// We can implement some exponential backoff in here
		select {
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(10 * time.Second)
			r.logger.Info("retrying", "nonce", nonce, "retry_number", i, "retries_left", r.retriesNumber-i)
			err = retryMethod(ctx, nonce)
			if err == nil {
				r.logger.Info("nonce processing succeeded", "nonce", nonce, "retries_number", i)
				return nil
			}
			r.logger.Error("failed to process nonce", "nonce", nonce, "retry", i, "err", err)
		}
	}
	return err
}

func (r retrier) RetryThenFail(nonce uint64, ctx context.Context, retryMethod func(context.Context, uint64) error) {
	err := r.Retry(nonce, ctx, retryMethod)
	if err != nil {
		panic(err)
	}
}

// mustGetEvent takes a corerpctypes.ResultEvent and checks whether it has
// the provided eventName. If not, it panics.
func mustGetEvent(result corerpctypes.ResultEvent, eventName string) []string {
	ev := result.Events[eventName]
	if ev == nil || len(ev) == 0 {
		panic(errors.Wrap(
			types.ErrEmpty,
			fmt.Sprintf(
				"%s not found in event %s",
				coretypes.EventTypeKey,
				result.Events,
			),
		))
	}
	return ev
}
