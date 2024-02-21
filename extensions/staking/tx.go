// Copyright 2024 Match Foundation
// This file is part of the Match Network packages.
//
// Match is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Match packages are distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Match packages. If not, see https://github.com/matchain/match/blob/main/LICENSE
package staking

import (
	"bytes"
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/matchain/match/extensions/basic"
	"github.com/matchain/match/x/evm/statedb"
)

func (extension *Extension) Delegate(
	ctx sdk.Context,
	origin common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	msg, delegator, err := NewDelegateMsg(args, extension.stakingKeeper.BondDenom(ctx))
	if err != nil {
		return nil, err
	}

	extension.Logger(ctx).Info(
		"tx called",
		"method", method.Name,
		"args", fmt.Sprintf(
			"{ delegator_address: %s, validator_address: %s, amount: %s }",
			delegator,
			msg.ValidatorAddress,
			msg.Amount.Amount,
		),
	)

	var (
		//stakeAuthz *stakingtypes.StakeAuthorization
		//expiration *time.Time

		isCallerOrigin    = contract.CallerAddress == origin
		isCallerDelegator = contract.CallerAddress == delegator
	)

	if isCallerDelegator {
		delegator = origin
	} else if origin != delegator {
		return nil, fmt.Errorf(ErrDifferentOriginFromDelegator, origin.String(), delegator.String())
	}

	if !isCallerOrigin {
	}

	msgSvr := stakingkeeper.NewMsgServerImpl(extension.stakingKeeper)
	if _, err := msgSvr.Delegate(sdk.WrapSDKContext(ctx), msg); err != nil {
		return nil, err
	}

	if !isCallerOrigin {
		//if err := extension.UpdateStakingAuthorization(ctx, contract.CallerAddress, delegator, stakeAuthz, expiration, DelegateMsg, msg); err != nil {
		//	return nil, err
		//}
	}

	err = extension.EmitDelegateEvent(ctx, stateDB, msg, delegator)
	if err != nil {
		return nil, err
	}

	if isCallerDelegator {
		stateDB.(*statedb.StateDB).SubBalance(contract.CallerAddress, msg.Amount.Amount.BigInt())
	}

	return method.Outputs.Pack(true)
}

func (extension *Extension) EmitDelegateEvent(
	ctx sdk.Context,
	stateDB vm.StateDB,
	msg *stakingtypes.MsgDelegate,
	delegator common.Address,
) error {
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return err
	}

	// Get the validator to estimate the new shares delegated
	// NOTE: At this point the validator has already been checked, so no need to check again
	validator, _ := extension.stakingKeeper.GetValidator(ctx, valAddr)

	// Get only the new shares based on the delegation amount
	newShares, err := validator.SharesFromTokens(msg.Amount.Amount)
	if err != nil {
		return err
	}

	// Prepare the event topics
	event := extension.ABI.Events[EventTypeDelegate]
	topics, err := extension.createStakingEthTxTopics(3, event, delegator, msg.ValidatorAddress)
	if err != nil {
		return err
	}

	// Prepare the event data
	var b bytes.Buffer
	b.Write(basic.PackNumber(reflect.ValueOf(msg.Amount.Amount.BigInt())))
	b.Write(basic.PackNumber(reflect.ValueOf(newShares.BigInt())))

	stateDB.AddLog(&ethtypes.Log{
		Address:     extension.Address(),
		Topics:      topics,
		Data:        b.Bytes(),
		BlockNumber: uint64(ctx.BlockHeight()),
	})

	return nil
}

func (extension *Extension) createStakingEthTxTopics(
	topicsLen uint64,
	event abi.Event,
	delegator common.Address,
	validator string,
) ([]common.Hash, error) {
	topics := make([]common.Hash, topicsLen)
	topics[0] = event.ID

	var err error
	topics[1], err = basic.MakeTopic(delegator)
	if err != nil {
		return nil, err
	}

	topics[2], err = basic.MakeTopic(validator)
	if err != nil {
		return nil, err
	}

	return topics, nil
}
