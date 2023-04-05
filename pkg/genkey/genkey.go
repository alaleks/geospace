// Package genkey performs the generation of random key.
package genkey

import (
	"crypto/rand"
	"math/big"
	rd "math/rand"
	"time"
)

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// CreateWithoutCrypto performs a random key generation specified size from chars using math/rand.
// This function dont use the cryptographic rand.
func CreateWithoutCrypto(sizeKey int) string {
	rd.Seed(time.Now().UnixNano())
	b := make([]byte, sizeKey)
	for i := range b {
		b[i] = chars[rd.Intn(len(chars))]
	}
	return string(b)
}

// CreateByBig performs a random key generation specified size from chars using package big.
func CreateByBig(sizeKey int) string {
	symbols := big.NewInt(int64(len(chars)))
	states := big.NewInt(0)
	states.Exp(symbols, big.NewInt(int64(sizeKey)), nil)
	r, err := rand.Int(rand.Reader, states)
	if err != nil {
		return CreateWithoutCrypto(sizeKey)
	}

	bytes := make([]byte, sizeKey)
	r2 := big.NewInt(0)
	symbol := big.NewInt(0)
	for i := range bytes {
		r2.DivMod(r, symbols, symbol)
		r, r2 = r2, r
		bytes[i] = chars[symbol.Int64()]
	}
	return string(bytes)
}

// Create performs a random key generation specified size.
func Create(sizeKey int) string {
	length := len(chars)
	b := make([]byte, sizeKey)

	_, err := rand.Read(b)
	if err != nil {
		return CreateWithoutCrypto(sizeKey)
	}

	for i := 0; i < sizeKey; i++ {
		b[i] = chars[int(b[i])%length]
	}
	return string(b)
}
