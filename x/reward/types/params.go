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
