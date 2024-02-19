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
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/matchain/match/x/reward/types"
)

type Keeper struct {
	// Store key required for the EVM Prefix KVStore. It is required by:
	// - storing account's Storage State
	// - storing account's Code
	// - storing transaction Logs
	// - storing Bloom filters by block height. Needed for the Web3 API.
	storeKey         storetypes.StoreKey
	transientKey     storetypes.StoreKey
	cdc              codec.BinaryCodec
	bankKeeper       types.BankKeeper
	stakingKeeper    types.StakingKeeper
	feeCollectorName string
}

func NewKeeper(
	storeKey, transientKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	feeCollectorName string,
) Keeper {
	return Keeper{
		storeKey:         storeKey,
		transientKey:     transientKey,
		cdc:              cdc,
		bankKeeper:       bankKeeper,
		stakingKeeper:    stakingKeeper,
		feeCollectorName: feeCollectorName,
	}
}

func (k Keeper) GetReward(ctx sdk.Context, contract string) types.Reward {
	var reward types.Reward
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(common.HexToAddress(contract).Bytes())
	if bz == nil {
		return reward
	}

	k.cdc.MustUnmarshal(bz, &reward)
	return reward
}

func (k Keeper) SetReward(ctx sdk.Context, reward types.Reward) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&reward)
	if err != nil {
		return err
	}

	contract := common.HexToAddress(reward.ContractAddress)
	store.Set(contract.Bytes(), bz)
	return nil
}
