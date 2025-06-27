package arbnode

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethdb"

	"github.com/offchainlabs/nitro/arbos/arbostypes"
	"github.com/offchainlabs/nitro/arbutil"
	"github.com/offchainlabs/nitro/staker"
	"github.com/offchainlabs/nitro/validator"
	"github.com/offchainlabs/nitro/validator/server_api"
)

type BlockValidatorAPI struct {
	val *staker.BlockValidator
}

func (a *BlockValidatorAPI) LatestValidated(ctx context.Context) (*staker.GlobalStateValidatedInfo, error) {
	return a.val.ReadLastValidatedInfo()
}

type BlockValidatorDebugAPI struct {
	val *staker.StatelessBlockValidator
}

type ValidateBlockResult struct {
	Valid       bool                    `json:"valid"`
	Latency     string                  `json:"latency"`
	GlobalState validator.GoGlobalState `json:"globalstate"`
}

func (a *BlockValidatorDebugAPI) ValidateMessageNumber(
	ctx context.Context, msgNum hexutil.Uint64, full bool, moduleRootOptional *common.Hash,
) (ValidateBlockResult, error) {
	result := ValidateBlockResult{}

	var moduleRoot common.Hash
	if moduleRootOptional != nil {
		moduleRoot = *moduleRootOptional
	} else {
		moduleRoot = a.val.GetLatestWasmModuleRoot()
	}
	start_time := time.Now()
	valid, gs, err := a.val.ValidateResult(ctx, arbutil.MessageIndex(msgNum), full, moduleRoot)
	result.Latency = fmt.Sprintf("%vms", time.Since(start_time).Milliseconds())
	if gs != nil {
		result.GlobalState = *gs
	}
	result.Valid = valid
	return result, err
}

func (a *BlockValidatorDebugAPI) ValidationInputsAt(ctx context.Context, msgNum hexutil.Uint64, target ethdb.WasmTarget,
) (server_api.InputJSON, error) {
	return a.val.ValidationInputsAt(ctx, arbutil.MessageIndex(msgNum), target)
}

type MaintenanceAPI struct {
	runner *MaintenanceRunner
}

func (a *MaintenanceAPI) SecondsSinceLastMaintenance(ctx context.Context) (int64, error) {
	running, since := a.runner.TimeSinceLastMaintenance()
	if running {
		return 0, errors.New("maintenance currently running")
	}
	return int64(since.Seconds()), nil
}

func (a *MaintenanceAPI) Trigger(ctx context.Context) error {
	return a.runner.Trigger()
}

type LightClientAPI struct {
	db ethdb.Database
}

func (api *LightClientAPI) GetLatestState(ctx context.Context) (*MessageTrackingL2Data, error) {

	latestStateIndexL1, err := GetLatestStateIndex(api.db, L1LatestStateIndexKey)
	if err != nil {
		return nil, err
	}

	latestStateIndexL2, err := GetLatestStateIndex(api.db, L2LatestStateIndexKey)
	if err != nil {
		return nil, err
	}

	latestStateIndex := min(latestStateIndexL1.StateIndex, latestStateIndexL2.StateIndex)

	return GetTrackingDataAt[MessageTrackingL2Data](api.db, latestStateIndex)
}

type StateAtReturnData struct {
	Message       arbostypes.L1IncomingMessage
	L2BlockNumber uint64
	L2BlockHash   common.Hash

	L1TxHash     common.Hash
	DataLocation batchDataLocation
}

func (api *LightClientAPI) GetFullDataAt(ctx context.Context, msgNum uint64) (*StateAtReturnData, error) {
	L1Data, err := GetTrackingDataAt[MessageTrackingL1Data](api.db, arbutil.MessageIndex(msgNum))
	if err != nil {
		return nil, err
	}
	L2Data, err := GetTrackingDataAt[MessageTrackingL2Data](api.db, arbutil.MessageIndex(msgNum))
	if err != nil {
		return nil, err
	}

	return &StateAtReturnData{
		Message:       L1Data.Message,
		L2BlockNumber: L2Data.L2BlockNumber,
		L2BlockHash:   L2Data.L2BlockHash,
		L1TxHash:      L1Data.L1TxHash,
		DataLocation:  L1Data.DataLocation,
	}, nil
}

func (api *LightClientAPI) GetLatestIndexL1(ctx context.Context) (*LatestStateIndex, error) {
	return GetLatestStateIndex(api.db, L1LatestStateIndexKey)
}

func (api *LightClientAPI) GetLatestIndexL2(ctx context.Context) (*LatestStateIndex, error) {
	return GetLatestStateIndex(api.db, L2LatestStateIndexKey)
}

func (api *LightClientAPI) GetL1DataAt(ctx context.Context, msgNum uint64) (*MessageTrackingL1Data, error) {
	return GetTrackingDataAt[MessageTrackingL1Data](api.db, arbutil.MessageIndex(msgNum))
}

func (api *LightClientAPI) GetL2DataAt(ctx context.Context, msgNum uint64) (*MessageTrackingL2Data, error) {
	return GetTrackingDataAt[MessageTrackingL2Data](api.db, arbutil.MessageIndex(msgNum))
}
