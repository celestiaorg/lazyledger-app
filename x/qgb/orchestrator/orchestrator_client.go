package orchestrator

import (
	"context"
	"fmt"
	"github.com/celestiaorg/celestia-app/x/qgb/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmlog "github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/types"
	"sync"
)

var _ AppClient = &orchestratorClient{}

type orchestratorClient struct {
	tendermintRPC       *http.HTTP
	logger              tmlog.Logger
	querier             Querier
	mutex               *sync.Mutex // TODO check if we need the mutex here
	orchestratorAddress string
}

func NewOrchestratorClient(
	logger tmlog.Logger,
	tendermintRPC string,
	querier Querier,
	orchAddr string,
) (AppClient, error) {
	trpc, err := http.New(tendermintRPC, "/websocket")
	if err != nil {
		return nil, err
	}
	err = trpc.Start()
	if err != nil {
		return nil, err
	}

	return &orchestratorClient{
		tendermintRPC:       trpc,
		logger:              logger,
		querier:             querier,
		mutex:               &sync.Mutex{},
		orchestratorAddress: orchAddr,
	}, nil
}

func contains(s []uint64, nonce uint64) bool {
	for _, v := range s {
		if v == nonce {
			return true
		}
	}
	return false
}

func (oc *orchestratorClient) SubscribeValset(ctx context.Context) (<-chan types.Valset, error) {
	valsetsChan := make(chan types.Valset, 100)

	// will change once we have the new design
	go oc.addOldVSAttestations(ctx, valsetsChan) //nolint:errcheck

	results, err := oc.tendermintRPC.Subscribe(
		ctx,
		"valset-changes",
		fmt.Sprintf("%s.%s='%s'", types.EventTypeValsetRequest, sdk.AttributeKeyModule, types.ModuleName),
	)
	if err != nil {
		return nil, err
	}

	nonces := make([]uint64, 10000)

	go func() {
		defer close(valsetsChan)
		for {
			select {
			case <-ctx.Done():
				return
			case <-results:
				valsets, err := oc.querier.QueryLastValsets(ctx)
				if err != nil {
					oc.logger.Error(err.Error())
					continue
				}

				// todo: double check that the first validator set is found
				if len(valsets) < 1 {
					oc.logger.Error("no validator sets found")
					continue
				}

				valset := valsets[0]

				// Checking if we already signed this valset
				resp, err := oc.querier.QueryValsetConfirm(ctx, valset.Nonce, oc.orchestratorAddress)
				if err != nil {
					oc.logger.Error(err.Error())
					continue
				}

				if resp == nil && !contains(nonces, valset.Nonce) {
					valsetsChan <- valset
					nonces = append(nonces, valset.Nonce)
				}
			}
		}
	}()

	return valsetsChan, nil
}

func (oc *orchestratorClient) addOldVSAttestations(ctx context.Context, valsetsChan chan types.Valset) error {
	oc.logger.Info("Started Valsets attestation signature catchup")
	lastUnbondingHeight, err := oc.querier.QueryLastUnbondingHeight(ctx)
	if err != nil {
		oc.logger.Error(err.Error())
		return err
	}

	valsets, err := oc.querier.QueryLastValsets(ctx)
	if err != nil {
		oc.logger.Error(err.Error())
		return err
	}

	// todo: double check that the first validator set is found
	if len(valsets) < 1 {
		oc.logger.Error("no validator sets found")
		return nil
	}
	valsetsChan <- valsets[0]

	previousNonce := valsets[0].Nonce
	for {
		if previousNonce == 1 {
			oc.logger.Info("Finished Valsets attestation signature catchup")
			return nil
		}
		previousNonce = previousNonce - 1
		lastVsConfirm, err := oc.querier.QueryValsetConfirm(ctx, previousNonce, oc.orchestratorAddress)
		if err != nil {
			oc.logger.Error(err.Error())
			return err
		}
		// The valset signed by the orchestrator to get lastVsConfirm
		// Used to get the height that valset waas first introduced
		correspondingVs, err := oc.querier.QueryValsetByNonce(ctx, previousNonce)
		if err != nil {
			oc.logger.Error(err.Error())
			return err
		}
		if int64(correspondingVs.Height) < lastUnbondingHeight {
			// Most likely, we're up to date and don't need to catchup anymore
			oc.logger.Info("Finished Valsets attestation signature catchup")
			return nil
		}
		if lastVsConfirm != nil {
			// in case we have holes in the signatures
			continue
		}

		valsetsChan <- *correspondingVs
	}
}

func (oc *orchestratorClient) SubscribeDataCommitment(ctx context.Context) (<-chan ExtendedDataCommitment, error) {
	dataCommitments := make(chan ExtendedDataCommitment, 100)

	// will change once we have the new design
	go oc.addOldDCAttestations(ctx, dataCommitments) //nolint:errcheck

	// queryClient := types.NewQueryClient(orchestratorClient.qgbRPC)

	// resp, err := queryClient.Params(ctx, &types.QueryParamsRequest{})
	// if err != nil {
	// 	return nil, err
	// }

	// params := resp.Params
	q := coretypes.EventQueryNewBlockHeader.String()
	results, err := oc.tendermintRPC.Subscribe(ctx, "height", q)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(dataCommitments)

		for {
			select {
			case <-ctx.Done():
				return
			case ev := <-results:
				eventDataHeader := ev.Data.(coretypes.EventDataNewBlockHeader)
				height := eventDataHeader.Header.Height
				// todo: refactor to ensure that no ranges of blocks are missed if the
				// parameters are changed
				if height%int64(types.DataCommitmentWindow) != 0 {
					continue
				}

				// TODO: calculate start height some other way that can handle changes
				// in the data window param
				startHeight := height - int64(types.DataCommitmentWindow)
				endHeight := height

				// create and send the data commitment
				dcResp, err := oc.tendermintRPC.DataCommitment(
					ctx,
					fmt.Sprintf("block.height >= %d AND block.height <= %d",
						startHeight,
						endHeight,
					),
				)
				if err != nil {
					oc.logger.Error(err.Error())
					continue
				}

				// TODO: store the nonce in the state somewhere, so that we don't have
				// to assume what the nonce is
				nonce := uint64(height) / types.DataCommitmentWindow

				dataCommitments <- ExtendedDataCommitment{
					Commitment: dcResp.DataCommitment,
					Start:      startHeight,
					End:        endHeight,
					Nonce:      nonce,
				}
			}
		}
	}()

	return dataCommitments, nil
}

func (oc *orchestratorClient) addOldDCAttestations(
	ctx context.Context,
	dataCommitmentsChan chan ExtendedDataCommitment,
) error {
	oc.logger.Info("Started Data Commitments attestation signature catchup")
	lastUnbondingHeight, err := oc.querier.QueryLastUnbondingHeight(ctx)
	if err != nil {
		oc.logger.Error(err.Error())
		return err
	}

	currentHeight, err := oc.querier.QueryHeight(ctx)
	if err != nil {
		oc.logger.Error(err.Error())
		return err
	}

	var previousBeginBlock int64
	var previousEndBlock int64

	if currentHeight%types.DataCommitmentWindow == 0 {
		previousBeginBlock = currentHeight
	} else {
		// to have a correct range
		previousBeginBlock = currentHeight - currentHeight%types.DataCommitmentWindow
	}

	for {
		// Will be refactored when we have data commitment requests
		previousEndBlock = previousBeginBlock
		previousBeginBlock = previousEndBlock - int64(types.DataCommitmentWindow)

		if previousBeginBlock == 0 {
			oc.logger.Info("Finished Data Commitments attestation signature catchup")
			return nil
		}

		previousCommitment, err := oc.tendermintRPC.DataCommitment(
			ctx,
			fmt.Sprintf("block.height >= %d AND block.height <= %d",
				previousBeginBlock,
				previousEndBlock,
			),
		)
		if err != nil {
			oc.logger.Error(err.Error())
			continue
		}

		existingConfirm, err := oc.querier.QueryDataCommitmentConfirm(
			ctx,
			previousCommitment.DataCommitment.String(),
			oc.orchestratorAddress,
		)
		if err != nil {
			oc.logger.Error(err.Error())
			continue
		}

		if previousEndBlock < lastUnbondingHeight {
			// Most likely, we're up to date and don't need to catchup anymore
			oc.logger.Info("Finished Data Commitments attestation signature catchup")
			return nil
		}
		if existingConfirm != nil {
			// In case we have holes in the signatures
			continue
		}
		previousNonce := uint64(previousEndBlock) / types.DataCommitmentWindow
		dataCommitmentsChan <- ExtendedDataCommitment{
			Commitment: previousCommitment.DataCommitment,
			Start:      previousBeginBlock,
			End:        previousEndBlock,
			Nonce:      previousNonce,
		}
	}
}
