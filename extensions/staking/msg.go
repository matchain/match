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
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/matchain/match/extensions/basic"
	"math/big"
)

func NewDelegateMsg(args []interface{}, denom string) (*stakingtypes.MsgDelegate, common.Address, error) {
	delegator, validator, amount, err := parseDelegateMsgArgs(args)
	if err != nil {
		return nil, ZeroAddress, err
	}

	msg := &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(delegator.Bytes()).String(),
		ValidatorAddress: validator,
		Amount: sdk.Coin{
			Denom:  denom,
			Amount: sdk.NewIntFromBigInt(amount),
		},
	}

	err = msg.ValidateBasic()
	if err != nil {
		return nil, ZeroAddress, err
	}

	return msg, delegator, nil
}

func parseDelegateMsgArgs(args []interface{}) (common.Address, string, *big.Int, error) {
	if len(args) != 3 {
		return ZeroAddress, "", nil, fmt.Errorf(basic.ErrInvalidArgs)
	}

	delegator, ok := args[0].(common.Address)
	if !ok || delegator == (ZeroAddress) {
		return ZeroAddress, "", nil, fmt.Errorf(ErrInvalidDelegator)
	}

	validator, ok := args[1].(string)
	if !ok {
		return ZeroAddress, "", nil, fmt.Errorf(ErrInvalidValidator)
	}

	amount, ok := args[2].(*big.Int)
	if !ok {
		return ZeroAddress, "", nil, fmt.Errorf(basic.ErrInvalidAmount)
	}

	return delegator, validator, amount, nil
}
