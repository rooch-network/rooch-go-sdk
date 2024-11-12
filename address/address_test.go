package address

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoochAddress(t *testing.T) {
	t.Run("should create RoochAddress when given a valid hex string", func(t *testing.T) {
		validHex := "0x1234567890abcdef"
		address, err := NewRoochAddress(validHex)

		assert.NoError(t, err)
		assert.NotNil(t, address)

		// Test ToHexAddress
		hexAddr := address.String()
		assert.Equal(t, NormalizeRoochAddress(validHex, true), hexAddr)

		// Test ToBech32Address
		bech32Addr, err := address.ToBech32()
		assert.NoError(t, err)
		assert.NotEmpty(t, bech32Addr)
		assert.True(t, strings.HasPrefix(bech32Addr, RoochBech32Prefix))
	})

	t.Run("should throw error when given an invalid hex string", func(t *testing.T) {
		invalidHex := "0x12345G"
		address, err := NewRoochAddress(invalidHex)

		assert.Error(t, err)
		assert.Nil(t, address)
	})
}
