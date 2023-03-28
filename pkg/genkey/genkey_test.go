package genkey_test

import (
	"fmt"
	"testing"

	"github.com/alaleks/geospace/pkg/genkey"
)

func TestGenRandKey(t *testing.T) {
	tests := []int{0, 8, 48, 256}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("gen key size %d", tt), func(t *testing.T) {
			key := genkey.GenRandKey(tt)
			if len(key) != tt {
				t.Errorf("an incorrect key size was generated: %d but should be %d", tt, len(key))
			}
		})
	}
}

func BenchmarkGenRandKey(b *testing.B) {
	tests := []int{8, 48, 64}
	for _, tt := range tests {
		b.ResetTimer()

		b.Run(fmt.Sprintf("gen key size %d", tt), func(b *testing.B) {
			_ = genkey.GenRandKey(tt)
		})
	}
}
