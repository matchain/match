package distribution

import (
	"bytes"
	"embed"
	"reflect"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/matchain/match/extensions/basic"
)

var _ vm.PrecompiledContract = &Extension{}

const ExtensionAddress = "0x6980000000000000000000000000000000000001"

// Embed abi json file to the executable binary. Needed when importing as dependency.
//
//go:embed abi.json
var f embed.FS

type Extension struct {
	basic.Extension
	distributionKeeper distributionkeeper.Keeper
	stakingKeeper      stakingkeeper.Keeper
}

func (extension *Extension) Address() common.Address {
	return common.HexToAddress(ExtensionAddress)
}

func (extension *Extension) RequiredGas(input []byte) uint64 {
	methodId := input[:4]

	method, err := extension.MethodById(methodId)
	if err != nil {
		return 0
	}

	return extension.Extension.RequiredGas(input, extension.isTransaction(method.Name))
}

func (extension *Extension) Run(evm *vm.EVM, contract *vm.Contract, readonly bool) (bz []byte, err error) {
	ctx, stateDB, method, initialGas, args, err := extension.Setup(evm, contract, readonly, extension.isTransaction)
	if err != nil {
		return nil, err
	}

	defer basic.HandleGasError(ctx, contract, initialGas, &err)()

	switch method.Name {
	case WithdrawDelegatorRewards:
		bz, err = extension.WithdrawDelegatorRewards(ctx, evm.Origin, contract, stateDB, method, args)
	}

	if err != nil {
		return nil, err
	}

	cost := ctx.GasMeter().GasConsumed() - initialGas

	if !contract.UseGas(cost) {
		return nil, vm.ErrOutOfGas
	}

	return bz, nil
}

func (extension *Extension) isTransaction(methodId string) bool {
	switch methodId {
	case ClaimRewards,
		SetWithdrawAddress,
		WithdrawDelegatorRewards,
		WithdrawValidatorCommission:
		return true
	default:
		return false
	}
}

func (extension *Extension) EmitWithdrawDelegatorRewardsEvent(ctx sdk.Context, stateDB vm.StateDB, delegator common.Address, validator string, coins sdk.Coins) error {
	event := extension.ABI.Events[EventTypeWithdrawDelegatorRewards]
	topics := make([]common.Hash, 3)

	// The first topic is always the signature of the event.
	topics[0] = event.ID

	var err error
	topics[1], err = basic.MakeTopic(delegator)
	if err != nil {
		return err
	}

	topics[2], err = basic.MakeTopic(validator)
	if err != nil {
		return err
	}

	// Prepare the event data
	var b bytes.Buffer
	b.Write(basic.PackNumber(reflect.ValueOf(coins[0].Amount.BigInt())))

	stateDB.AddLog(&ethtypes.Log{
		Address:     extension.Address(),
		Topics:      topics,
		Data:        b.Bytes(),
		BlockNumber: uint64(ctx.BlockHeight()),
	})

	return nil
}

// loadAbi loads staking contract for staking module.
func loadAbi() (abi.ABI, error) {
	return basic.LoadAbi(f, "abi.json")
}

func NewExtension(
	dk distributionkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
	authzKeeper authzkeeper.Keeper,
) (*Extension, error) {
	abi, err := loadAbi()
	if err != nil {
		return nil, err
	}

	return &Extension{
		Extension: basic.Extension{
			ABI:                  abi,
			AuthzKeeper:          authzKeeper,
			ApprovalExpiration:   basic.DefaultExpirationDuration,
			KvGasConfig:          storetypes.KVGasConfig(),
			TransientKvGasConfig: storetypes.TransientGasConfig(),
		},
		distributionKeeper: dk,
		stakingKeeper:      stakingKeeper,
	}, nil
}
