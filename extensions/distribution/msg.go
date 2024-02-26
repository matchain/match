package distribution

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/matchain/match/extensions/basic"
)

func NewMsgWithdrawDelegatorReward(args []interface{}) (*distributiontypes.MsgWithdrawDelegatorReward, common.Address, error) {
	// todo: parse arguments
	if len(args) != 2 {
		return nil, basic.ZeroAddress, nil
	}

	delegator, ok := args[0].(common.Address)
	if !ok || delegator == (basic.ZeroAddress) {
		return nil, basic.ZeroAddress, fmt.Errorf(ErrInvalidDelegator)
	}

	validator, _ := args[1].(string)

	msg := &distributiontypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: sdk.AccAddress(delegator.Bytes()).String(),
		ValidatorAddress: validator,
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, basic.ZeroAddress, err
	}

	return msg, delegator, nil
}
