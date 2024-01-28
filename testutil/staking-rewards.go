// Copyright 2022 Match Foundation
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

package testutil

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// GetTotalDelegationRewards returns the total delegation rewards that are currently
// outstanding for the given address.
func GetTotalDelegationRewards(ctx sdk.Context, distributionKeeper distributionkeeper.Keeper, addr sdk.AccAddress) (sdk.DecCoins, error) {
	resp, err := distributionKeeper.DelegationTotalRewards(
		ctx,
		&distributiontypes.QueryDelegationTotalRewardsRequest{
			DelegatorAddress: addr.String(),
		},
	)
	if err != nil {
		return nil, err
	}

	return resp.Total, nil
}
