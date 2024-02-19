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
	"time"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
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
