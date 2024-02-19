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
	"embed"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/matchain/match/extensions/basic"
)

// Embed abi json file to the executable binary. Needed when importing as dependency.
//
//go:embed abi.json
var f embed.FS

const ExtensionAddress = "0x6980000000000000000000000000000000000000"

type Extension struct {
	basic.Extension
	stakingKeeper stakingkeeper.Keeper
}

// loadAbi loads staking contract for staking module.
func loadAbi() (abi.ABI, error) {
	return basic.LoadAbi(f, "abi.json")
}

func NewExtension(
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
		stakingKeeper: stakingKeeper,
	}, nil
}

// Address returns ethereum address for staking extension.
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

// todo
func (extension *Extension) Run(input []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

// isTransaction returns if given method is valid.
func (*Extension) isTransaction(method string) bool {
	switch method {
	case DelegateMethod,
		UndelegateMethod,
		RedelegateMethod,
		CancelUnbondingDelegationMethod:
		//authorization.ApproveMethod,
		//authorization.RevokeMethod,
		//authorization.IncreaseAllowanceMethod,
		//authorization.DecreaseAllowanceMethod,
		return true
	default:
		return false
	}
}
