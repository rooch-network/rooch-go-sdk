// Copyright (c) RoochNetwork
// SPDX-License-Identifier: Apache-2.0

package test_e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

// TestBox equivalent structure
type TestBox struct {
	keypair          *KeyPair
	bitcoinContainer *Container
	ordContainer     *Container
	client           *RoochClient
	// Add other necessary fields
}

func (tb *TestBox) setup() *TestBox {
	// Initialize TestBox
	return tb
}

func (tb *TestBox) loadBitcoinEnv() error {
	// Bitcoin environment setup
	return nil
}

func (tb *TestBox) loadORDEnv() error {
	// ORD environment setup
	return nil
}

func (tb *TestBox) loadRoochEnv(env string, index int) error {
	// Rooch environment setup
	return nil
}

func (tb *TestBox) cleanEnv() {
	// Cleanup environment
}

func (tb *TestBox) delay(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

func TestBitcoinAssets(t *testing.T) {
	var testBox *TestBox

	// Setup
	t.Run("Setup", func(t *testing.T) {
		testBox = new(TestBox).setup()
		err := testBox.loadBitcoinEnv()
		if err != nil {
			t.Fatal(err)
		}
		err = testBox.loadORDEnv()
		if err != nil {
			t.Fatal(err)
		}
		err = testBox.loadRoochEnv("local", 0)
		if err != nil {
			t.Fatal(err)
		}
	})

	// Cleanup after all tests
	defer testBox.cleanEnv()

	// Test UTXO query
	t.Run("Query UTXO", func(t *testing.T) {
		addr := testBox.keypair.GetSchnorrPublicKey().ToAddress().BitcoinAddress.ToStr()

		result, err := testBox.bitcoinContainer.ExecuteRpcCommandRaw([]string{}, "generatetoaddress", []string{"50", addr})
		if err != nil {
			t.Fatal(err)
		}
		if result == nil {
			t.Fatal("Expected result to be defined")
		}

		// Wait for rooch indexer
		testBox.delay(10)

		utxos, err := testBox.client.QueryUTXO(QueryFilter{
			Owner: addr,
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(utxos.Data) == 0 {
			t.Fatal("Expected utxos length to be greater than 0")
		}
	})

	// Test Inscriptions query
	t.Run("Query Inscriptions", func(t *testing.T) {
		// Init wallet
		result, err := testBox.ordContainer.ExecCmd("wallet create")
		if err != nil || result.ExitCode != 0 {
			t.Fatal("Failed to create wallet")
		}

		result, err = testBox.ordContainer.ExecCmd("wallet receive")
		if err != nil || result.ExitCode != 0 {
			t.Fatal("Failed to receive wallet")
		}

		var response struct {
			Addresses []string `json:"addresses"`
		}
		if err := json.Unmarshal([]byte(result.Output), &response); err != nil {
			t.Fatal(err)
		}
		addr := response.Addresses[0]

		// Mint UTXO
		result, err = testBox.bitcoinContainer.ExecuteRpcCommandRaw([]string{}, "generatetoaddress", []string{"101", addr})
		if err != nil {
			t.Fatal(err)
		}

		// Wait for ord sync and index
		testBox.delay(10)

		// Check wallet balance
		result, err = testBox.ordContainer.ExecCmd("wallet balance")
		if err != nil || result.ExitCode != 0 {
			t.Fatal("Failed to get wallet balance")
		}

		var balanceResponse struct {
			Total int64 `json:"total"`
		}
		if err := json.Unmarshal([]byte(result.Output), &balanceResponse); err != nil {
			t.Fatal(err)
		}
		if balanceResponse.Total != 5000000000 {
			t.Fatalf("Expected balance 5000000000, got %d", balanceResponse.Total)
		}

		// Create inscription
		inscriptionContent := `{"p":"brc-20","op":"mint","tick":"Rooch","amt":"1"}`
		filePath := fmt.Sprintf("/%s/hello.txt", testBox.ordContainer.GetHostDataPath())
		if err := os.WriteFile(filePath, []byte(inscriptionContent), 0644); err != nil {
			t.Fatal(err)
		}

		result, err = testBox.ordContainer.ExecCmd(fmt.Sprintf("wallet inscribe --fee-rate 1 --file /data/hello.txt --destination %s", addr))
		if err != nil || result.ExitCode != 0 {
			t.Fatal("Failed to inscribe")
		}

		// Generate more blocks
		result, err = testBox.bitcoinContainer.ExecuteRpcCommandRaw([]string{}, "generatetoaddress", []string{"1", addr})
		if err != nil {
			t.Fatal(err)
		}

		// Wait for rooch indexer
		testBox.delay(10)

		// Query UTXOs
		utxos, err := testBox.client.QueryUTXO(QueryFilter{
			Owner: addr,
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(utxos.Data) == 0 {
			t.Fatal("Expected utxos length to be greater than 0")
		}

		// Query Inscriptions
		inscriptions, err := testBox.client.QueryInscriptions(QueryFilter{
			Owner: addr,
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(inscriptions.Data) == 0 {
			t.Fatal("Expected inscriptions length to be greater than 0")
		}
	})
}
