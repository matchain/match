package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: Params{
			Enable:    true,
			Base:      sdk.NewDecWithPrec(2, 1),
			Validator: sdk.NewDecWithPrec(6, 1),
		},
	}
}

func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
