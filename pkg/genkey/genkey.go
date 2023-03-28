// Package genkey performs the generation of random key.
package genkey

import (
	"crypto/rand"
	"math/big"
)

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenRandKey performs a random key generation specified size.
func GenRandKey(sizeKey int) string {
	var (
		key         string
		restrictInt = int64(len(chars))
	)

	for {
		// return a random key if length string equals size key.
		if len(key) == sizeKey {
			return key
		}

		// if err is not nil to return default key.
		n, err := rand.Int(rand.Reader, big.NewInt(restrictInt))
		if err != nil {
			return chars[:sizeKey]
		}

		key += string(chars[n.Int64()])
	}
}
