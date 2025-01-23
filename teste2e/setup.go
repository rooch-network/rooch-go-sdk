// Copyright (c) RoochNetwork
// SPDX-License-Identifier: Apache-2.0

package test_e2e

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/rooch-network/rooch-go-sdk/address"
	"github.com/rooch-network/rooch-go-sdk/client"
	//"github.com/rooch-network/rooch-go-sdk/keypairs"
	//"github.com/rooch-network/rooch-go-sdk/types"
	//"github.com/rooch-network/rooch-go-sdk/testsuite/container"
)

var DefaultNodeURL string

func init() {
	// Similar to the TypeScript version, getting URL from env or default
	DefaultNodeURL = getEnvOrDefault("VITE_FULLNODE_URL", client.GetRoochNodeURL("localnet"))
}

type TestBox struct {
	container.RoochContainer
	client  *client.RoochClient
	keypair *keypair.Secp256k1Keypair
}

// Setup creates a new TestBox instance with a generated keypair
func Setup() *TestBox {
	kp := keypair.GenerateSecp256k1Keypair()
	return &TestBox{
		keypair: kp,
		client:  client.NewRoochClient(DefaultNodeURL),
	}
}

// LoadRoochEnv initializes the Rooch environment
func (tb *TestBox) LoadRoochEnv(target interface{}, port int) error {
	if port == 0 {
		port = 6768
	}

	// Call parent LoadRoochEnv (implementation depends on RoochContainer)
	if err := tb.RoochContainer.LoadRoochEnv(target, port); err != nil {
		return err
	}

	roochServerAddr := tb.GetRoochServerAddress()
	tb.client = client.NewRoochClient(fmt.Sprintf("http://%s", roochServerAddr))
	return nil
}

// GetClient returns the RoochClient instance
func (tb *TestBox) GetClient() *client.RoochClient {
	return tb.client
}

// Address returns the Rooch address associated with the keypair
func (tb *TestBox) Address() address.RoochAddress {
	return tb.keypair.GetRoochAddress()
}

// SignAndExecuteTransaction signs and executes a transaction
func (tb *TestBox) SignAndExecuteTransaction(tx *transaction.Transaction) (bool, error) {
	result, err := tb.client.SignAndExecuteTransaction(&client.SignAndExecuteTransactionRequest{
		Transaction: tx,
		Signer:      tb.keypair,
	})
	if err != nil {
		return false, err
	}

	return result.ExecutionInfo.Status.Type == "executed", nil
}

type PublishPackageOptions struct {
	NamedAddresses string
}

// PublishPackage publishes a Move package
func (tb *TestBox) PublishPackage(packagePath string, box *TestBox, options *PublishPackageOptions) (bool, error) {
	if options == nil {
		options = &PublishPackageOptions{
			NamedAddresses: "rooch_examples=default",
		}
	}

	namedAddresses := strings.ReplaceAll(options.NamedAddresses, "default", box.Address().ToHexAddress())

	if err := tb.RoochCommand(fmt.Sprintf("move build -p %s --named-addresses %s --install-dir %s --json",
		packagePath, namedAddresses, tb.GetTmpDir())); err != nil {
		return false, err
	}

	// Read package file
	fileBytes, err := ioutil.ReadFile(tb.GetTmpDir() + "/package.rpd")
	if err != nil {
		return false, err
	}

	// Create and execute transaction
	tx := transaction.NewTransaction()
	tx.CallFunction(&transaction.FunctionCall{
		Target: "0x2::module_store::publish_package_entry",
		Args:   []interface{}{fileBytes},
	})

	return box.SignAndExecuteTransaction(tx)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
