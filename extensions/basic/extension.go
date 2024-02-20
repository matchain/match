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
package basic

import (
	"fmt"
	"time"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/matchain/match/x/evm/statedb"
)

type Extension struct {
	abi.ABI
	AuthzKeeper          authzkeeper.Keeper
	ApprovalExpiration   time.Duration
	KvGasConfig          storetypes.GasConfig
	TransientKvGasConfig storetypes.GasConfig
}

func (extension *Extension) RequiredGas(input []byte, isTransaction bool) uint64 {
	argsBz := input[4:]

	if isTransaction {
		return extension.KvGasConfig.WriteCostFlat + (extension.KvGasConfig.WriteCostPerByte * uint64(len(argsBz)))
	}

	return extension.KvGasConfig.ReadCostFlat + (extension.KvGasConfig.ReadCostPerByte * uint64(len(argsBz)))
}

func (extension *Extension) Setup(
	evm *vm.EVM, contract *vm.Contract,
	readOnly bool, isTransaction func(name string) bool,
) (sdk.Context, *statedb.StateDB, *abi.Method, sdk.Gas, []interface{}, error) {
	stateDB, ok := evm.StateDB.(*statedb.StateDB)
	if !ok {
		return sdk.Context{}, nil, nil, 0, nil, fmt.Errorf(ErrNotInEvm)
	}
	ctx := stateDB.GetContext()

	methodId := contract.Input[:4]
	method, err := extension.MethodById(methodId)
	if err != nil {
		return sdk.Context{}, nil, nil, 0, nil, err
	}

	if readOnly && isTransaction(method.Name) {
		return sdk.Context{}, nil, nil, 0, nil, vm.ErrWriteProtection
	}

	argsBz := contract.Input[4:]
	args, err := method.Inputs.Unpack(argsBz)
	if err != nil {
		return sdk.Context{}, nil, nil, 0, nil, err
	}

	initialGas := ctx.GasMeter().GasConsumed()
	defer HandleGasError(ctx, contract, initialGas, &err)()

	ctx = ctx.WithGasMeter(sdk.NewGasMeter(contract.Gas)).
		WithKVGasConfig(extension.KvGasConfig).
		WithTransientKVGasConfig(extension.TransientKvGasConfig)
	ctx.GasMeter().ConsumeGas(initialGas, "creating new gas meter")

	return ctx, stateDB, method, initialGas, args, nil
}

func HandleGasError(ctx sdk.Context, contract *vm.Contract, initialGas sdk.Gas, err *error) func() {
	return func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case sdk.ErrorOutOfGas:
				usedGas := ctx.GasMeter().GasConsumed() - initialGas
				_ = contract.UseGas(usedGas)

				*err = vm.ErrOutOfGas
				ctx = ctx.
					WithKVGasConfig(storetypes.GasConfig{}).
					WithTransientKVGasConfig(storetypes.GasConfig{})
			default:
				panic(r)
			}
		}
	}
}
