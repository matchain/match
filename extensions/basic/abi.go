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
package basic

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func LoadAbi(fs embed.FS, path string) (abi.ABI, error) {
	abiBz, err := fs.ReadFile(path)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to read abi %v", err)
	}

	newAbi, err := abi.JSON(bytes.NewReader(abiBz))
	if err != nil {
		return abi.ABI{}, fmt.Errorf("invalid abi %v", err)
	}

	return newAbi, nil
}
