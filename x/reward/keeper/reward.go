package keeper

import (
	"errors"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"

	matchtypes "github.com/matchain/match/types"
	"github.com/matchain/match/x/reward/types"
)

func (k Keeper) Send(ctx sdk.Context, msg core.Message, usedGas uint64, denom string) error {
	params := k.GetParams(ctx)
	if !params.Enable {
		return nil
	}

	txFee := sdk.NewIntFromUint64(usedGas).Mul(sdk.NewIntFromBigInt(msg.GasPrice()))
	base := params.Base
	baseFee := base.MulInt(txFee).TruncateInt()
	baseFees := sdk.Coins{{
		Denom:  denom,
		Amount: baseFee,
	}}

	contract := msg.To().String()
	reward := k.GetReward(ctx, contract)
	developer := sdk.AccAddress(common.HexToAddress(reward.DeployerAddress).Bytes())

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.feeCollectorName, developer, baseFees); err != nil {
		return errorsmod.Wrapf(
			err,
			"fee collector account failed to distribute developer fees (%s) to developer address %s. contract %s",
			baseFees, developer, contract,
		)
	}

	return nil
}

func (k Keeper) Burn(ctx sdk.Context, msg core.Message, usedGas uint64, denom string) error {
	params := k.GetParams(ctx)
	if !params.Enable {
		return nil
	}

	burn := sdk.NewDec(1).Sub(params.Base).Sub(params.Validator)
	if burn.IsNegative() {
		return errorsmod.Wrapf(errors.New("burn ratio can't be negative"), "ratio %d", burn)
	}

	txFee := sdk.NewIntFromUint64(usedGas).Mul(sdk.NewIntFromBigInt(msg.GasPrice()))
	burntFee := burn.MulInt(txFee).TruncateInt()
	burntFees := sdk.Coins{{
		Denom:  denom,
		Amount: burntFee,
	}}

	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, burntFees)
	if err != nil {
		return errorsmod.Wrapf(
			err,
			"fee collector account failed to send fees (%s) to module %s.",
			burntFees, types.ModuleName,
		)
	}

	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, burntFees)
	if err != nil {
		return errorsmod.Wrapf(
			err,
			"reward account failed to burn fees (%s).",
			burntFees,
		)
	}

	return nil
}

func (k *Keeper) MintReward(ctx sdk.Context) error {
	baseWei := big.NewInt(10)
	baseRewardBigInt := baseWei.Exp(baseWei, big.NewInt(matchtypes.BaseDenomUnit), nil)
	baseRewardCoin := sdk.NewCoin(matchtypes.Match, sdk.NewIntFromBigInt(baseRewardBigInt))
	baseRewardCoins := sdk.NewCoins(baseRewardCoin)
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, baseRewardCoins); err != nil {
		return errorsmod.Wrapf(err, "failed to mint block reward (%s)", baseRewardCoin.String())
	}

	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, baseRewardCoins); err != nil {
		return errorsmod.Wrapf(err, "failed to send block reward (%s) to (%s)", baseRewardCoin.String(), k.feeCollectorName)
	}

	return nil
}
