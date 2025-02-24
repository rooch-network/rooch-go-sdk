package secp256k1

import (
	_ "encoding/base64"
	"github.com/rooch-network/rooch-go-sdk/crypto"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCases = []struct {
	sk string
	pk string
}{
	{
		sk: "roochsecretkey1q9rc3ryrp644d33yy4d2c9mg7wnuuxag7mqs0uq6yp7nmv6yd7usu2j6v3z",
		pk: "Au4i1I9dB6BvAQ+aX8mt4f/wVKjYLhOkD6LEcgB/WBjq",
	},
}

func TestSecp256k1Keypair(t *testing.T) {
	t.Run("Create secp256k1 keypair", func(t *testing.T) {
		kp, _ := GenerateSecp256k1Keypair()
		assert.Equal(t, 33, len(kp.GetPublicKey().ToBytes()))
	})

	t.Run("Export secp256k1 keypair", func(t *testing.T) {
		kp, _ := GenerateSecp256k1Keypair()
		secret := kp.GetSecretKey()
		assert.True(t, strings.HasPrefix(string(secret), crypto.RoochSecretKeyPrefix))
	})

	t.Run("Create secp256k1 keypair from CLI secret key", func(t *testing.T) {
		testKey := "roochsecretkey1q969zv4rhqpuj0nkf2e644yppjf34p6zwr3gq0633qc7n9luzg6w6lycezc"
		expectRoochHexAddress := "0xf892b3fd5fd0e93436ba3dc8d504413769d66901266143d00e49441079243ed0"
		expectRoochBech32Address := "rooch1lzft8l2l6r5ngd468hyd2pzpxa5av6gpyes585qwf9zpq7fy8mgqh9npj5"
		expectNostrAddress := "npub1h54r2zvulk96qjmfnyy83mtry0pp5acnz6uvk637typxtvn90c8s0lrc0g"
		//expectBitcoinAddress := "bcrt1pw9l5h7vepq8cnpugwm848x3at34gg5eq0mamdrjw0krunfjm0zfq65gjzz"

		sk, _ := FromSecp256k1SecretKey(testKey, false)
		addrView, _ := sk.GetSchnorrPublicKey().ToAddress()

		bench32Addr, _ := addrView.RoochAddress.ToBech32()
		assert.Equal(t, expectRoochHexAddress, addrView.RoochAddress.String())
		assert.Equal(t, expectRoochBech32Address, bench32Addr)
		assert.Equal(t, expectNostrAddress, addrView.NostrAddress.ToStr())
		//assert.Equal(t, expectBitcoinAddress, addrView.BitcoinAddress.String())
	})

	t.Run("Create secp256k1 keypair from secret key", func(t *testing.T) {
		// valid secret key is provided by rooch keystore
		sk := testCases[0].sk
		pk := testCases[0].pk

		key, _ := crypto.DecodeRoochSecretKey(sk)
		keypair, _ := FromSecp256k1SecretKey(key.SecretKey, false)
		assert.Equal(t, pk, keypair.GetPublicKey().ToBase64())

		keypair1, _ := FromSecp256k1SecretKey(key.SecretKey, false)
		assert.Equal(t, pk, keypair1.GetPublicKey().ToBase64())
	})

	t.Run("sign", func(t *testing.T) {
		t.Run("should sign data", func(t *testing.T) {
			keypair, _ := NewSecp256k1Keypair(nil)
			message := []byte("hello world")
			signature, err := keypair.Sign(message)
			assert.NoError(t, err)

			isValid, err := keypair.GetPublicKey().Verify(message, signature)
			assert.NoError(t, err)
			assert.True(t, isValid)
		})

		t.Run("Sign data same as rooch cli", func(t *testing.T) {
			// TODO:
		})
	})
}
