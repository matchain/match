package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (k *Keeper) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {}

func (k *Keeper) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	err := k.MintReward(ctx)
	if err != nil {
		panic(err)
	}

	return []abci.ValidatorUpdate{}
}
