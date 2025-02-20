package address

import (
//"testing"

// "github.com/stretchr/testify/assert"
)

//type TestCase struct {
//	BtcAddr   string
//	RoochAddr string
//	HexAddr   string
//}
//
//var testCases = []TestCase{
//	{
//		BtcAddr:   "18cBEMRxXHqzWWCxZNtU91F5sbUNKhL5PX",
//		RoochAddr: "rooch1gxterelcypsyvh8cc9kg73dtnyct822ykx8pmu383qruzt4r93jshtc9fj",
//		HexAddr:   "0x419791e7f82060465cf8c16c8f45ab9930b3a944b18e1df2278807c12ea32c65",
//	},
//	{
//		BtcAddr:   "bc1q262qeyyhdakrje5qaux8m2a3r4z8sw8vu5mysh",
//		RoochAddr: "rooch10lnft7hhq37vl0y97lwvkmzqt48fk76y0z88rfcu8zg6qm8qegfqx0rq2h",
//		HexAddr:   "0x7fe695faf7047ccfbc85f7dccb6c405d4e9b7b44788e71a71c3891a06ce0ca12",
//	},
//}
//
//func TestBitcoinAddress(t *testing.T) {
//	t.Run("New address with ed25519 keypair", func(t *testing.T) {
//		kp := keypairs.GenerateEd25519Keypair()
//		address := kp.GetPublicKey().ToAddress()
//		assert.NotNil(t, address)
//	})
//
//	t.Run("From address", func(t *testing.T) {
//		btcAddr := testCases[0].BtcAddr
//		addr, err := NewBitcoinAddress(btcAddr)
//		assert.NoError(t, err)
//		assert.NotNil(t, addr)
//	})
//
//	t.Run("To rooch address", func(t *testing.T) {
//		for _, item := range testCases {
//			addr, err := NewBitcoinAddress(item.BtcAddr)
//			assert.NoError(t, err)
//
//			roochAddr := addr.GenRoochAddress()
//
//			genRoochAddr := roochAddr.ToBech32Address()
//			genRoochHexAddr := roochAddr.ToHexAddress()
//
//			assert.True(t, IsValidAddress(genRoochAddr))
//			assert.True(t, IsValidAddress(genRoochHexAddr))
//			assert.Equal(t, item.RoochAddr, genRoochAddr)
//			assert.Equal(t, item.HexAddr, genRoochHexAddr)
//		}
//	})
//
//	t.Run("From hex address", func(t *testing.T) {
//		hexAddr := "020145966003624094dae2deeb30815eedd38f96c45c3fdb1261f5d697fc4137e0de"
//		expectBTCAddr := "bc1pgktxqqmzgz2d4ck7avcgzhhd6w8ed3zu8ld3yc0466tlcsfhur0qj3y0wm"
//
//		btcAddr, err := NewBitcoinAddress(hexAddr)
//		assert.NoError(t, err)
//
//		assert.Equal(t, expectBTCAddr, btcAddr.ToStr())
//	})
//}
