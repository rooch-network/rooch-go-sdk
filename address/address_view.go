// Copyright (c) RoochNetwork
// SPDX-License-Identifier: Apache-2.0

package address

// AddressView implements the Address interface
type AddressView struct {
	BitcoinAddress BitcoinAddress
	NostrAddress   NostrAddress
	RoochAddress   RoochAddress
}

// NewAddressView creates a new AddressView instance
func NewAddressView(publicKey []byte) (*AddressView, error) {
	return NewAddressViewWithNetwork(publicKey, BitcoinNetworkRegtest)
}

// NewAddressView creates a new AddressView instance
func NewAddressViewWithNetwork(publicKey []byte, network BitcoinNetworkType) (*AddressView, error) {
	bitcoinAddr, err := BitcoinAddressFromPublicKey(publicKey, network)
	if err != nil {
		return nil, err
	}
	nostrAddr, err := NewNostrAddress(publicKey)
	if err != nil {
		return nil, err
	}
	roochAddr, err := bitcoinAddr.GenRoochAddress()
	if err != nil {
		return nil, err
	}

	return &AddressView{
		BitcoinAddress: *bitcoinAddr,
		NostrAddress:   *nostrAddr,
		RoochAddress:   *roochAddr,
	}, nil
}

// ToBytes returns the byte representation of the address
func (av *AddressView) ToBytes() []byte {
	return av.RoochAddress.Bytes()
}

// ToStr returns the string representation of the address
func (av *AddressView) ToStr() string {
	return av.RoochAddress.String()
}
