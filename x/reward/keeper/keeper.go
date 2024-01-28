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
