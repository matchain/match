package distribution

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/matchain/match/extensions/basic"
	"github.com/matchain/match/x/evm/statedb"
)

func (extension *Extension) WithdrawDelegatorRewards(
	ctx sdk.Context, origin common.Address, contract *vm.Contract,
	stateDB vm.StateDB, method *abi.Method, args []interface{},
) ([]byte, error) {
	msg, delegator, err := NewMsgWithdrawDelegatorReward(args)
	if err != nil {
		return nil, err
	}

	isContractorDelegator := contract.CallerAddress == delegator
	if !isContractorDelegator && origin != delegator {
		return nil, fmt.Errorf("invalid caller")
	}

	msgSvr := distributionkeeper.NewMsgServerImpl(extension.distributionKeeper)
	res, err := msgSvr.WithdrawDelegatorReward(sdk.WrapSDKContext(ctx), msg)
	if err != nil {
		return nil, err
	}

	if err = extension.EmitWithdrawDelegatorRewardsEvent(ctx, stateDB, delegator, msg.ValidatorAddress, res.Amount); err != nil {
		return nil, err
	}

	if isContractorDelegator {
		stateDB.(*statedb.StateDB).AddBalance(contract.CallerAddress, res.Amount[0].Amount.BigInt())
	}

	return method.Outputs.Pack(basic.NewCoinsResponse(res.Amount))
}
