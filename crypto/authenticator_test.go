package crypto

import (
	"bytes"
	"crypto/ed25519"
	"github.com/rooch-network/rooch-go-sdk/bcs"
	"github.com/rooch-network/rooch-go-sdk/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthenticationKey_FromPublicKey(t *testing.T) {
	// Ed25519
	privateKey, err := GenerateEd25519PrivateKey()
	assert.NoError(t, err)
	publicKey := privateKey.PubKey()

	authKey := AuthenticationKey{}
	authKey.FromPublicKey(publicKey)

	hash := utils.Sha3256Hash([][]byte{
		publicKey.Bytes(),
		{Ed25519Scheme},
	})

	assert.Equal(t, hash[:], authKey[:])
}

func Test_AuthenticationKeySerialization(t *testing.T) {
	bytesWithLength := []byte{
		32,
		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
	}
	bytes := []byte{
		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
	}
	authKey := AuthenticationKey(bytes)
	serialized, err := bcs.Serialize(&authKey)
	assert.NoError(t, err)
	assert.Equal(t, bytesWithLength, serialized)

	newAuthKey := AuthenticationKey{}
	err = bcs.Deserialize(&newAuthKey, serialized)
	assert.NoError(t, err)
	assert.Equal(t, authKey, newAuthKey)
}

func Test_AuthenticatorSerialization(t *testing.T) {
	msg := []byte{0x01, 0x02}
	privateKey, err := GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	authenticator, err := privateKey.Sign(msg)
	assert.NoError(t, err)

	serialized, err := bcs.Serialize(authenticator)
	assert.NoError(t, err)
	assert.Equal(t, uint8(AccountAuthenticatorEd25519), serialized[0])
	assert.Len(t, serialized, 1+(1+ed25519.PublicKeySize)+(1+ed25519.SignatureSize))

	newAuthenticator := &AccountAuthenticator{}
	err = bcs.Deserialize(newAuthenticator, serialized)
	assert.NoError(t, err)
	assert.Equal(t, authenticator.Variant, newAuthenticator.Variant)
	assert.Equal(t, authenticator.Auth, newAuthenticator.Auth)
}

func Test_AuthenticatorVerification(t *testing.T) {
	msg := []byte{0x01, 0x02}
	privateKey, err := GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	authenticator, err := privateKey.Sign(msg)
	assert.NoError(t, err)

	assert.True(t, authenticator.Verify(msg))
}

func Test_InvalidAuthenticatorDeserialization(t *testing.T) {
	serialized := []byte{0xFF}
	newAuthenticator := &AccountAuthenticator{}
	err := bcs.Deserialize(newAuthenticator, serialized)
	assert.Error(t, err)
	serialized = []byte{0x4F}
	newAuthenticator = &AccountAuthenticator{}
	err = bcs.Deserialize(newAuthenticator, serialized)
	assert.Error(t, err)
}

func Test_InvalidAuthenticationKeyDeserialization(t *testing.T) {
	serialized := []byte{0xFF}
	newAuthkey := AuthenticationKey{}
	err := bcs.Deserialize(&newAuthkey, serialized)
	assert.Error(t, err)
}

func TestBitcoinSignMessage(t *testing.T) {
	t.Run("should correctly construct with valid txData and messageInfo", func(t *testing.T) {
		txData := []byte{1, 2, 3, 4}
		messageInfo := "Test Message Info"
		bitcoinSignMessage := NewBitcoinSignMessage(txData, messageInfo)

		if bitcoinSignMessage.MessagePrefix != "\u0018Bitcoin Signed Message:\n" {
			t.Errorf("Expected message prefix to be '\\u0018Bitcoin Signed Message:\\n', got %s", bitcoinSignMessage.MessagePrefix)
		}

		expectedMessageInfo := "Rooch Transaction:\nTest Message Info\n"
		if bitcoinSignMessage.MessageInfo != expectedMessageInfo {
			t.Errorf("Expected message info to be '%s', got %s", expectedMessageInfo, bitcoinSignMessage.MessageInfo)
		}

		if !bytes.Equal(bitcoinSignMessage.TxHash, txData) {
			t.Errorf("Expected txHash to equal input txData")
		}
	})

	t.Run("should correctly generate raw message string", func(t *testing.T) {
		txData := []byte{1, 2, 3, 4}
		messageInfo := "Test Message Info"
		bitcoinSignMessage := NewBitcoinSignMessage(txData, messageInfo)

		if !bytes.Equal(bitcoinSignMessage.TxHash, txData) {
			t.Errorf("Expected txHash to equal input txData")
		}

		expected := "Rooch Transaction:\nTest Message Info\n01020304"
		if bitcoinSignMessage.Raw() != expected {
			t.Errorf("Expected raw message to be '%s', got %s", expected, bitcoinSignMessage.Raw())
		}
	})

	t.Run("should handle empty messageInfo gracefully", func(t *testing.T) {
		txData := []byte{}
		messageInfo := ""
		bitcoinSignMessage := NewBitcoinSignMessage(txData, messageInfo)

		expectedMessageInfo := "Rooch Transaction:\n"
		if bitcoinSignMessage.MessageInfo != expectedMessageInfo {
			t.Errorf("Expected message info to be '%s', got %s", expectedMessageInfo, bitcoinSignMessage.MessageInfo)
		}

		if bitcoinSignMessage.Raw() != expectedMessageInfo {
			t.Errorf("Expected raw message to be '%s', got %s", expectedMessageInfo, bitcoinSignMessage.Raw())
		}
	})

	t.Run("should correctly encode message with valid txHash and messageInfo", func(t *testing.T) {
		txData := []byte{0x01, 0x02, 0x03, 0x04}
		messageInfo := "Example message info"
		bitcoinSignMessage := NewBitcoinSignMessage(txData, messageInfo)
		encodedData := bitcoinSignMessage.Encode()

		if len(encodedData) == 0 {
			t.Error("Expected encoded data to not be empty")
		}

		if len(encodedData) > 255 {
			t.Error("Expected encoded data length to be less than or equal to 255")
		}
	})
}
