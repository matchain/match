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
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/matchain/match/x/reward/types"
)

type Hooks struct {
	k Keeper
}

func (h Hooks) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	return h.k.PostTxProcessing(ctx, msg, receipt)
}

func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

func (k Keeper) PostTxProcessing(
	ctx sdk.Context,
	msg core.Message,
	receipt *ethtypes.Receipt,
) error {
	params := k.GetParams(ctx)
	if !params.Enable {
		return nil
	}

	contract := msg.To()
	if contract == nil {
		reward := types.Reward{
			ContractAddress: receipt.ContractAddress.String(),
			DeployerAddress: msg.From().String(),
		}
		err := k.SetReward(ctx, reward)
		if err != nil {
			return errorsmod.Wrapf(
				err,
				"failed to set reward for contract %s and deployer %s",
				reward.ContractAddress, reward.DeployerAddress,
			)
		}

		return nil
	}

	//acct := k.evmKeeper.GetAccountWithoutBalance(ctx, *contract)
	//txFee := sdk.NewIntFromUint64(receipt.GasUsed).Mul(sdk.NewIntFromBigInt(msg.GasPrice()))
	//evmDenom := k.evmKeeper.GetParams(ctx).EvmDenom
	//ratio := sdk.NewDec(1)
	//
	//burntFee := params.Base.MulInt(txFee).TruncateInt()
	//burntFees := sdk.Coins{{
	//	Denom:  evmDenom,
	//	Amount: burntFee,
	//}}
	//if err := k.bankKeeper.BurnCoins(ctx, k.feeCollectorName, burntFees); err != nil {
	//	return errorsmod.Wrapf(
	//		err,
	//		"fee collector account failed to burn fees (%s). contract %s",
	//		burntFees, contract,
	//	)
	//}
	//ratio = ratio.Sub(params.Base)
	//
	//if acct.IsContract() {
	//	reward := k.GetReward(ctx, contract.String())
	//	developer := sdk.AccAddress(common.HexToAddress(reward.DeployerAddress).Bytes())
	//	developerFee := params.Base.MulInt(txFee).TruncateInt()
	//	developerFees := sdk.Coins{{
	//		Denom:  evmDenom,
	//		Amount: developerFee,
	//	}}
	//	err := k.bankKeeper.SendCoinsFromModuleToAccount(
	//		ctx,
	//		k.feeCollectorName,
	//		developer,
	//		developerFees)
	//	if err != nil {
	//		return errorsmod.Wrapf(
	//			err,
	//			"fee collector account failed to distribute developer fees (%s) to withdraw address %s. contract %s",
	//			developerFees, developer, contract,
	//		)
	//	}
	//
	//	ratio = ratio.Sub(params.Base)
	//}

	return nil
}
