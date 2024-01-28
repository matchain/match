package types

import (
	"fmt"
)

var (
	ParamsKey = []byte("Params")
)

func (p Params) Validate() error {
	base := p.Base
	if base.IsNil() {
		return fmt.Errorf("invalid parameter: nil")
	}
	if base.IsNegative() {
		return fmt.Errorf("value cannot be negative: %T", base)
	}

	validator := p.Validator
	if validator.IsNil() {
		return fmt.Errorf("invalid parameter: nil")
	}
	if validator.IsNegative() {
		return fmt.Errorf("value cannot be negative: %T", validator)
	}

	return nil
}
