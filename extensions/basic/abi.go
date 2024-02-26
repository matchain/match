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
	"math/big"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
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

func PackNumber(value reflect.Value) []byte {
	switch kind := value.Kind(); kind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return math.U256Bytes(new(big.Int).SetUint64(value.Uint()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return math.U256Bytes(big.NewInt(value.Int()))
	case reflect.Ptr:
		return math.U256Bytes(new(big.Int).Set(value.Interface().(*big.Int)))
	default:
		panic("abi: invalid number type")
	}
}

func genIntType(value int64, size uint) []byte {
	var topic [common.HashLength]byte
	if value < 0 {
		topic = [common.HashLength]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	}
	for i := uint(0); i < size; i++ {
		topic[common.HashLength-i-1] = byte(value >> (i * 8))
	}
	return topic[:]
}

func MakeTopic(value interface{}) (common.Hash, error) {
	var topic common.Hash

	// Try to generate the topic based on simple types
	switch value := value.(type) {
	case common.Hash:
		copy(topic[:], value[:])
	case common.Address:
		copy(topic[common.HashLength-common.AddressLength:], value[:])
	case *big.Int:
		blob := value.Bytes()
		copy(topic[common.HashLength-len(blob):], blob)
	case bool:
		if value {
			topic[common.HashLength-1] = 1
		}
	case int8:
		copy(topic[:], genIntType(int64(value), 1))
	case int16:
		copy(topic[:], genIntType(int64(value), 2))
	case int32:
		copy(topic[:], genIntType(int64(value), 4))
	case int64:
		copy(topic[:], genIntType(value, 8))
	case uint8:
		blob := new(big.Int).SetUint64(uint64(value)).Bytes()
		copy(topic[common.HashLength-len(blob):], blob)
	case uint16:
		blob := new(big.Int).SetUint64(uint64(value)).Bytes()
		copy(topic[common.HashLength-len(blob):], blob)
	case uint32:
		blob := new(big.Int).SetUint64(uint64(value)).Bytes()
		copy(topic[common.HashLength-len(blob):], blob)
	case uint64:
		blob := new(big.Int).SetUint64(value).Bytes()
		copy(topic[common.HashLength-len(blob):], blob)
	case string:
		hash := crypto.Keccak256Hash([]byte(value))
		copy(topic[:], hash[:])
	case []byte:
		hash := crypto.Keccak256Hash(value)
		copy(topic[:], hash[:])

	default:
		val := reflect.ValueOf(value)
		switch {
		// static byte array
		case val.Kind() == reflect.Array && reflect.TypeOf(value).Elem().Kind() == reflect.Uint8:
			reflect.Copy(reflect.ValueOf(topic[:val.Len()]), val)
		default:
			return topic, fmt.Errorf("unsupported indexed type: %T", value)
		}
	}

	return topic, nil
}

type Coin struct {
	Denom  string
	Amount *big.Int
}

func NewCoinsResponse(amount sdk.Coins) []Coin {
	outputs := make([]Coin, len(amount))
	for i, coin := range amount {
		outputs[i] = Coin{
			Denom:  coin.Denom,
			Amount: coin.Amount.BigInt(),
		}
	}
	return outputs
}
