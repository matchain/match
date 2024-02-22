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
package keeper

import (
	"fmt"
	"maps"

	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	stakingextension "github.com/matchain/match/extensions/staking"
)

func AvailableExtensions(
	stakingKeeper stakingkeeper.Keeper,
	distributionKeeper distributionkeeper.Keeper,
	authzKeeper authzkeeper.Keeper,
) map[common.Address]vm.PrecompiledContract {
	extensions := maps.Clone(vm.PrecompiledContractsBerlin)

	stakingExtension, err := stakingextension.NewExtension(stakingKeeper, authzKeeper)
	if err != nil {
		panic(fmt.Errorf("failed to load staking extension %v", err))
	}

	extensions[stakingExtension.Address()] = stakingExtension

	return extensions
}

// WithExtensions sets extensions for mapping between
// cosmos modules and contracts.
func (k *Keeper) WithExtensions(extensions map[common.Address]vm.PrecompiledContract) *Keeper {
	if k.extensions != nil {
		panic("contracts for extensions already set")
	}

	if len(extensions) == 0 {
		panic("empty extensions")
	}

	k.extensions = extensions
	return k
}

// Extensions returns available extensions.
func (k *Keeper) Extensions() map[common.Address]vm.PrecompiledContract {
	return k.extensions
}
