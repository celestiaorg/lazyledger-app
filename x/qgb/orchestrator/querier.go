package orchestrator

import (
	"context"
	"fmt"
	"time"

	"github.com/tendermint/tendermint/rpc/client/http"

	"github.com/celestiaorg/celestia-app/x/qgb/types"
	tmlog "github.com/tendermint/tendermint/libs/log"
	"google.golang.org/grpc"
)

var _ Querier = &querier{}

type Querier interface {
	//QueryLatestDataCommitmentNonce(ctx context.Context) (uint64, error)
	//QueryLastDataCommitments(ctx context.Context) ([]types.DataCommitment, error)
	QueryDataCommitmentByNonce(ctx context.Context, nonce uint64) (*types.DataCommitment, error)
	QueryDataCommitmentConfirms(ctx context.Context, commit string) ([]types.MsgDataCommitmentConfirm, error)
	QueryDataCommitmentConfirm(
		ctx context.Context,
		endBlock uint64,
		beginBlock uint64,
		address string,
	) (*types.MsgDataCommitmentConfirm, error)
	//QueryLastValset(ctx context.Context) (types.Valset, error)
	QueryTwoThirdsDataCommitmentConfirms(
		ctx context.Context,
		timeout time.Duration,
		dc ExtendedDataCommitment,
	) ([]types.MsgDataCommitmentConfirm, error)
	QueryTwoThirdsValsetConfirms(
		ctx context.Context,
		timeout time.Duration,
		valset types.Valset,
	) ([]types.MsgValsetConfirm, error)
	//QueryLastValsets(ctx context.Context) ([]types.Valset, error)
	QueryValsetConfirm(ctx context.Context, nonce uint64, address string) (*types.MsgValsetConfirm, error)
	QueryValsetByNonce(ctx context.Context, nonce uint64) (*types.Valset, error)
	QueryLastUnbondingHeight(ctx context.Context) (uint64, error)
	QueryHeight(ctx context.Context) (uint64, error)
	QueryLastValsetBeforeNonce(
		ctx context.Context,
		nonce uint64,
	) (*types.Valset, error)
	QueryDataCommitmentConfirmsByExactRange(
		ctx context.Context,
		start uint64,
		end uint64,
	) ([]types.MsgDataCommitmentConfirm, error)
	QueryLatestAttestationNonce(ctx context.Context) (uint64, error)
}

type querier struct {
	qgbRPC        *grpc.ClientConn
	logger        tmlog.Logger
	tendermintRPC *http.HTTP
}

func NewQuerier(qgbRPCAddr, tendermintRPC string, logger tmlog.Logger) (*querier, error) {
	qgbGRPC, err := grpc.Dial(qgbRPCAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	trpc, err := http.New(tendermintRPC, "/websocket")
	if err != nil {
		return nil, err
	}
	err = trpc.Start()
	if err != nil {
		return nil, err
	}

	return &querier{
		qgbRPC:        qgbGRPC,
		logger:        logger,
		tendermintRPC: trpc,
	}, nil
}

// TODO add the other stop methods for other clients
func (q *querier) Stop() {
	err := q.qgbRPC.Close()
	if err != nil {
		q.logger.Error(err.Error())
	}
	err = q.tendermintRPC.Stop()
	if err != nil {
		q.logger.Error(err.Error())
	}
}

func (q *querier) QueryDataCommitmentConfirms(
	ctx context.Context,
	commit string,
) ([]types.MsgDataCommitmentConfirm, error) {
	queryClient := types.NewQueryClient(q.qgbRPC)

	confirmsResp, err := queryClient.DataCommitmentConfirmsByCommitment(
		ctx,
		&types.QueryDataCommitmentConfirmsByCommitmentRequest{
			Commitment: commit,
		},
	)
	if err != nil {
		return nil, err
	}

	return confirmsResp.Confirms, nil
}

func (q *querier) QueryTwoThirdsDataCommitmentConfirms(
	ctx context.Context,
	timeout time.Duration,
	dc ExtendedDataCommitment,
) ([]types.MsgDataCommitmentConfirm, error) {
	valset, err := q.QueryLastValsetBeforeNonce(ctx, dc.Data.Nonce)
	if err != nil {
		return nil, err
	}

	// create a map to easily search for power
	vals := make(map[string]types.BridgeValidator)
	for _, val := range valset.Members {
		vals[val.GetEthereumAddress()] = val
	}

	majThreshHold := valset.TwoThirdsThreshold()

	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case <-time.After(timeout):
			return nil, fmt.Errorf("failure to query for majority validator set confirms: timout %s", timeout)
		default:
			currThreshHold := uint64(0)
			confirms, err := q.QueryDataCommitmentConfirmsByExactRange(ctx, dc.Data.BeginBlock, dc.Data.EndBlock)
			if err != nil {
				return nil, err
			}

			correctConfirms := make([]types.MsgDataCommitmentConfirm, 0)
			for _, dataCommitmentConfirm := range confirms {
				val, has := vals[dataCommitmentConfirm.EthAddress]
				if !has {
					// currently, the orchestrators sign everything even if they didn't exist during a certain valset
					// thus, the relayer finds correct confirms and also incorrect ones. By incorrect, I mean signatures from
					// orchestrators that didn't belong to the valset in question, but they still signed it
					// as part of their catching up mechanism.
					// should be fixed with the new design and https://github.com/celestiaorg/celestia-app/issues/406
					q.logger.Debug(fmt.Sprintf(
						"dataCommitmentConfirm signer not found in stored validator set: address %s nonce %d",
						val.EthereumAddress,
						valset.Nonce,
					))
					continue
				}
				currThreshHold += val.Power
				correctConfirms = append(correctConfirms, dataCommitmentConfirm)
			}

			if currThreshHold >= majThreshHold {
				return correctConfirms, nil
			}
			q.logger.Debug(
				fmt.Sprintf(
					"found DataCommitmentConfirms total power %d number of confirms %d",
					currThreshHold,
					len(confirms),
				),
			)
		}
		// TODO: make the timeout configurable
		time.Sleep(time.Second * 30)
	}
}

func (q *querier) QueryTwoThirdsValsetConfirms(
	ctx context.Context,
	timeout time.Duration,
	valset types.Valset,
) ([]types.MsgValsetConfirm, error) {
	var currentValset types.Valset
	if valset.Nonce == 1 {
		currentValset = valset
	} else {
		vs, err := q.QueryLastValsetBeforeNonce(ctx, valset.Nonce)
		if err != nil {
			return nil, err
		}
		currentValset = *vs
	}
	// create a map to easily search for power
	vals := make(map[string]types.BridgeValidator)
	for _, val := range currentValset.Members {
		vals[val.GetEthereumAddress()] = val
	}

	majThreshHold := valset.TwoThirdsThreshold()

	for {
		select {
		case <-ctx.Done():
			return nil, nil
		// TODO: remove this extra case, and we can instead rely on the caller to pass a context with a timeout
		case <-time.After(timeout):
			return nil, fmt.Errorf("failure to query for majority validator set confirms: timout %s", timeout)
		default:
			currThreshHold := uint64(0)
			queryClient := types.NewQueryClient(q.qgbRPC)
			confirmsResp, err := queryClient.ValsetConfirmsByNonce(ctx, &types.QueryValsetConfirmsByNonceRequest{
				Nonce: valset.Nonce,
			})
			if err != nil {
				return nil, err
			}

			confirms := make([]types.MsgValsetConfirm, 0)
			for _, valsetConfirm := range confirmsResp.Confirms {
				val, has := vals[valsetConfirm.EthAddress]
				if !has {
					// currently, the orchestrators sign everything even if they didn't exist during a certain valset
					// thus, the relayer finds correct confirms and also incorrect ones. By incorrect, I mean signatures from
					// orchestrators that didn't belong to the valset in question, but they still signed it
					// as part of their catching up mechanism.
					// should be fixed with the new design. and https://github.com/celestiaorg/celestia-app/issues/406
					q.logger.Debug(
						fmt.Sprintf(
							"valSetConfirm signer not found in stored validator set: address %s nonce %d",
							val.EthereumAddress,
							valset.Nonce,
						))
					continue
				}
				currThreshHold += val.Power
				confirms = append(confirms, valsetConfirm)
			}

			if currThreshHold >= majThreshHold {
				return confirms, nil
			}
			q.logger.Debug(
				"foundValsetConfirms",
				fmt.Sprintf(
					"total power %d number of confirms %d",
					currThreshHold,
					len(confirmsResp.Confirms),
				),
			)
		}
		// TODO: make the timeout configurable
		time.Sleep(time.Second * 30)
	}
}

// QueryLastValset TODO change name to reflect the functionality correctly
// TODO make this return a pointer
//func (q *querier) QueryLastValset(ctx context.Context) (types.Valset, error) {
//	queryClient := types.NewQueryClient(q.qgbRPC)
//	lastValsetResp, err := queryClient.LastValsetRequests(ctx, &types.QueryLastValsetRequestsRequest{})
//	if err != nil {
//		return types.Valset{}, err
//	}
//
//	if len(lastValsetResp.Valsets) == 1 {
//		// genesis case
//		return lastValsetResp.Valsets[0], nil
//	}
//
//	if len(lastValsetResp.Valsets) < 2 {
//		return types.Valset{}, errors.New("no validator sets found")
//	}
//
//	valset := lastValsetResp.Valsets[1]
//	return valset, nil
//}

//func (q *querier) QueryLastValsets(ctx context.Context) ([]types.Valset, error) {
//	queryClient := types.NewQueryClient(q.qgbRPC)
//	lastValsetResp, err := queryClient.LastValsetRequests(ctx, &types.QueryLastValsetRequestsRequest{})
//	if err != nil {
//		return nil, err
//	}
//
//	return lastValsetResp.Valsets, nil
//}

func (q *querier) QueryLastValsetBeforeNonce(ctx context.Context, nonce uint64) (*types.Valset, error) {
	queryClient := types.NewQueryClient(q.qgbRPC)
	resp, err := queryClient.LastValsetBeforeNonce(
		ctx,
		&types.QueryLastValsetBeforeNonceRequest{Nonce: nonce},
	)
	if err != nil {
		return nil, err
	}

	return resp.Valset, nil
}

func (q *querier) QueryValsetConfirm(
	ctx context.Context,
	nonce uint64,
	address string,
) (*types.MsgValsetConfirm, error) {
	queryClient := types.NewQueryClient(q.qgbRPC)
	// machi bdarora this nonce is for a valset confirm
	resp, err := queryClient.ValsetConfirm(ctx, &types.QueryValsetConfirmRequest{Nonce: nonce, Address: address})
	if err != nil {
		return nil, err
	}

	return resp.Confirm, nil
}

func (q *querier) QueryHeight(ctx context.Context) (uint64, error) {
	resp, err := q.tendermintRPC.Status(ctx)
	if err != nil {
		return 0, err
	}

	return uint64(resp.SyncInfo.LatestBlockHeight), nil
}

func (q *querier) QueryLastUnbondingHeight(ctx context.Context) (uint64, error) {
	queryClient := types.NewQueryClient(q.qgbRPC)
	resp, err := queryClient.LastUnbondingHeight(ctx, &types.QueryLastUnbondingHeightRequest{})
	if err != nil {
		return 0, err
	}

	return resp.Height, nil
}

func (q *querier) QueryDataCommitmentConfirm(
	ctx context.Context,
	endBlock uint64,
	beginBlock uint64,
	address string,
) (*types.MsgDataCommitmentConfirm, error) {
	queryClient := types.NewQueryClient(q.qgbRPC)

	confirmsResp, err := queryClient.DataCommitmentConfirm(
		ctx,
		&types.QueryDataCommitmentConfirmRequest{
			EndBlock:   endBlock,
			BeginBlock: beginBlock,
			Address:    address,
		},
	)
	if err != nil {
		return nil, err
	}

	return confirmsResp.Confirm, nil
}

func (q *querier) QueryDataCommitmentConfirmsByExactRange(
	ctx context.Context,
	start uint64,
	end uint64,
) ([]types.MsgDataCommitmentConfirm, error) {
	queryClient := types.NewQueryClient(q.qgbRPC)
	confirmsResp, err := queryClient.DataCommitmentConfirmsByExactRange(
		ctx,
		&types.QueryDataCommitmentConfirmsByExactRangeRequest{
			BeginBlock: start,
			EndBlock:   end,
		},
	)
	if err != nil {
		return nil, err
	}
	return confirmsResp.Confirms, nil
}

func (q *querier) QueryDataCommitmentByNonce(ctx context.Context, nonce uint64) (*types.DataCommitment, error) {
	queryClient := types.NewQueryClient(q.qgbRPC)

	dc, err := queryClient.AttestationRequestByNonce(
		ctx,
		&types.QueryAttestationRequestByNonceRequest{Nonce: nonce},
	)
	if err != nil {
		return nil, err
	}

	encCfg := MakeEncodingConfig()

	var unmarshalledAttestation types.AttestationRequestI
	err = encCfg.InterfaceRegistry.UnpackAny(dc.Attestation, &unmarshalledAttestation)
	if err != nil {
		return nil, err
	}

	if unmarshalledAttestation.Type() != types.DataCommitmentRequestType {
		return nil, types.ErrAttestationNotDataCommitmentRequest
	}

	dcc, ok := unmarshalledAttestation.(*types.DataCommitment)
	if !ok {
		return nil, types.ErrAttestationNotDataCommitmentRequest
	}

	return dcc, nil
}

func (q *querier) QueryValsetByNonce(ctx context.Context, nonce uint64) (*types.Valset, error) {
	queryClient := types.NewQueryClient(q.qgbRPC)

	dc, err := queryClient.AttestationRequestByNonce(
		ctx,
		&types.QueryAttestationRequestByNonceRequest{Nonce: nonce},
	)
	if err != nil {
		return nil, err
	}

	value, ok := dc.Attestation.GetCachedValue().(types.Valset)
	if !ok {
		return nil, types.ErrAttestationNotValsetRequest
	}

	return &value, nil
}

func (q *querier) QueryLatestAttestationNonce(ctx context.Context) (uint64, error) {
	queryClient := types.NewQueryClient(q.qgbRPC)

	resp, err := queryClient.LatestAttestationNonce(
		ctx,
		&types.QueryLatestAttestationNonceRequest{},
	)
	if err != nil {
		return 0, err
	}

	return resp.Nonce, nil
}
