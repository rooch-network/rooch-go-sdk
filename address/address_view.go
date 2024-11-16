// Copyright (c) RoochNetwork
// SPDX-License-Identifier: Apache-2.0

package address

// AddressView implements the Address interface
type AddressView struct {
	bitcoinAddress BitcoinAddress
	nostrAddress   NostrAddress
	roochAddress   RoochAddress
}

// NewAddressView creates a new AddressView instance
func NewAddressView(publicKey []byte, network BitcoinNetworkType) (*AddressView, error) {
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
		bitcoinAddress: *bitcoinAddr,
		nostrAddress:   *nostrAddr,
		roochAddress:   *roochAddr,
	}, nil
}

// ToBytes returns the byte representation of the address
func (av *AddressView) ToBytes() []byte {
	return av.roochAddress.Bytes()
}

// ToStr returns the string representation of the address
func (av *AddressView) ToStr() string {
	return av.roochAddress.String()
}
