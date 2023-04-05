package genkey_test

import (
	"fmt"
	"testing"

	"github.com/alaleks/geospace/pkg/genkey"
)

func TestCreateByBig(t *testing.T) {
	tests := []int{8, 64, 256}

	// testing byte size
	for _, tt := range tests {
		t.Run(fmt.Sprintf("gen key size %d", tt), func(t *testing.T) {
			key := genkey.CreateByBig(tt)
			if len(key) != tt {
				t.Errorf("an incorrect key size was generated: %d but should be %d", tt, len(key))
			}
		})
	}

	// testing uniq key
	for _, tt := range tests {
		t.Run(fmt.Sprintf("gen key size %d", tt), func(t *testing.T) {
			key1 := genkey.CreateByBig(tt)
			key2 := genkey.CreateByBig(tt)
			key3 := genkey.CreateByBig(tt)
			if key1 == key2 || key1 == key3 || key2 == key3 {
				t.Errorf("keys match: key1: %s , key2: %s, key3: %s",
					key1, key2, key3)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	tests := []int{8, 64, 256}

	// testing byte size
	for _, tt := range tests {
		t.Run(fmt.Sprintf("gen key size %d", tt), func(t *testing.T) {
			key := genkey.Create(tt)
			if len([]byte(key)) != tt {
				t.Errorf("an incorrect key size was generated: %d but should be %d", tt, len(key))
			}
		})
	}

	// testing uniq key
	for _, tt := range tests {
		t.Run(fmt.Sprintf("gen key size %d", tt), func(t *testing.T) {
			key1 := genkey.Create(tt)
			key2 := genkey.Create(tt)
			key3 := genkey.Create(tt)
			if key1 == key2 || key1 == key3 || key2 == key3 {
				t.Errorf("keys match: key1: %s , key2: %s, key3: %s",
					key1, key2, key3)
			}
		})
	}
}

func TestCreateWithoutCrypto(t *testing.T) {
	tests := []int{8, 64, 256}

	// testing byte size
	for _, tt := range tests {
		t.Run(fmt.Sprintf("gen key size %d", tt), func(t *testing.T) {
			key := genkey.CreateWithoutCrypto(tt)
			if len([]byte(key)) != tt {
				t.Errorf("an incorrect key size was generated: %d but should be %d", tt, len(key))
			}
		})
	}

	// testing uniq key
	for _, tt := range tests {
		t.Run(fmt.Sprintf("gen key size %d", tt), func(t *testing.T) {
			key1 := genkey.CreateWithoutCrypto(tt)
			key2 := genkey.CreateWithoutCrypto(tt)
			key3 := genkey.CreateWithoutCrypto(tt)
			if key1 == key2 || key1 == key3 || key2 == key3 {
				t.Errorf("keys match: key1: %s , key2: %s, key3: %s",
					key1, key2, key3)
			}
		})
	}
}

func BenchmarkCreateByBig(b *testing.B) {
	tests := []int{8, 64, 256}
	for _, tt := range tests {
		b.ResetTimer()

		b.Run(fmt.Sprintf("gen key size %d", tt), func(b *testing.B) {
			_ = genkey.CreateByBig(tt)
		})
	}
}

func BenchmarkCreate(b *testing.B) {
	tests := []int{8, 64, 256}
	for _, tt := range tests {
		b.ResetTimer()

		b.Run(fmt.Sprintf("gen key size %d", tt), func(b *testing.B) {
			_ = genkey.Create(tt)
		})
	}
}

func BenchmarkCreateWithoutCrypto(b *testing.B) {
	tests := []int{8, 64, 256}
	for _, tt := range tests {
		b.ResetTimer()

		b.Run(fmt.Sprintf("gen key size %d", tt), func(b *testing.B) {
			_ = genkey.CreateWithoutCrypto(tt)
		})
	}
}
