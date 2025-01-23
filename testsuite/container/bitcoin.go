// Copyright (c) RoochNetwork
// SPDX-License-Identifier: Apache-2.0

package container

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultImage = "lncm/bitcoind:v25.1"
)

var bitcoinPorts = []string{"18443", "18444", "28333", "28332"}

type BitcoinContainer struct {
	testcontainers.Container
	rpcBind      string
	rpcUser      string
	rpcPass      string
	hostDataPath string
}

type BitcoinContainerOption func(*BitcoinContainer)

// NewBitcoinContainer creates a new Bitcoin container with the given options
func NewBitcoinContainer(opts ...BitcoinContainerOption) (*BitcoinContainer, error) {
	c := &BitcoinContainer{
		rpcBind: "0.0.0.0:18443",
		rpcUser: "roochuser",
		rpcPass: "roochpass",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// WithHostDataPath sets the host data path for Bitcoin
func (c *BitcoinContainer) WithHostDataPath(hostPath string) {
	//return func(c *BitcoinContainer) {
	c.hostDataPath = filepath.Join(hostPath, "bitcoin")
	os.MkdirAll(c.hostDataPath, 0755)
}

// WithRpcBind sets the RPC bind address
func (c *BitcoinContainer) WithRpcBind(rpcBind string) {
	//return func(c *BitcoinContainer) {
	c.rpcBind = rpcBind
	//}
}

// WithRpcUser sets the RPC username
func (c *BitcoinContainer) WithRpcUser(rpcUser string) {
	//return func(c *BitcoinContainer) {
	c.rpcUser = rpcUser
	//}
}

// WithRpcPass sets the RPC password
func (c *BitcoinContainer) WithRpcPass(rpcPass string) {
	//return func(c *BitcoinContainer) {
	c.rpcPass = rpcPass
	//}
}

func (c *BitcoinContainer) GetRpcUrl() string{
	c.
}


func (c *BitcoinContainer) generateRpcauth() (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, salt)
	h.Write([]byte(c.rpcPass))
	passwordHmac := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s:%x$%s", c.rpcUser, salt, passwordHmac), nil
}

// Start starts the Bitcoin container
func (c *BitcoinContainer) Start(ctx context.Context) (*StartedRoochContainer, error) {
	if c.hostDataPath == "" {
		return fmt.Errorf("Bitcoin host config path not set. Call WithHostDataPath() before initializing")
	}

	rpcauth, err := c.generateRpcauth()
	if err != nil {
		return err
	}

	req := testcontainers.ContainerRequest{
		Image:        defaultImage,
		ExposedPorts: bitcoinPorts,
		Env: map[string]string{
			"RPC_BIND": c.rpcBind,
			"RPC_USER": c.rpcUser,
			"RPC_PASS": c.rpcPass,
			"RPC_AUTH": rpcauth,
		},
		User: "root",
		Cmd: []string{
			"-chain=regtest",
			"-txindex=1",
			"-fallbackfee=0.00001",
			"-zmqpubrawblock=tcp://0.0.0.0:28332",
			"-zmqpubrawtx=tcp://0.0.0.0:28333",
			"-rpcallowip=0.0.0.0/0",
			fmt.Sprintf("-rpcbind=%s", c.rpcBind),
			fmt.Sprintf("-rpcauth=%s", rpcauth),
		},
		//BindMounts: map[string]string{
		//	c.hostDataPath: "/data/.bitcoin",
		//},
		Mounts: testcontainers.Mounts(
			testcontainers.VolumeMount(c.hostDataPath, "/data/.bitcoin"),
		),
		WaitingFor: wait.ForLog("txindex thread start").WithStartupTimeout(120 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return err
	}

	c.Container = container
	return nil
}

// ExecuteRpcCommand executes a Bitcoin RPC command
func (c *BitcoinContainer) ExecuteRpcCommand(ctx context.Context, command string, params ...interface{}) (string, error) {
	strParams := make([]string, len(params))
	for i, param := range params {
		strParams[i] = fmt.Sprintf("%v", param)
	}
	return c.ExecuteRpcCommandRaw(ctx, []string{}, command, strParams...)
}

// ExecuteRpcCommandRaw executes a Bitcoin RPC command with raw parameters
func (c *BitcoinContainer) ExecuteRpcCommandRaw(ctx context.Context, opts []string, command string, params ...string) (string, error) {
	cmd := append([]string{"bitcoin-cli", "-regtest"}, opts...)
	cmd = append(cmd, command)
	cmd = append(cmd, params...)

	exitCode, reader, err := c.Container.Exec(ctx, cmd)
	if err != nil {
		return "", err
	}

	if exitCode != 0 {
		return "", fmt.Errorf("executeRpcCommand failed with exit code %d for command: %s", exitCode, command)
	}

	// Read the output from the io.Reader
	output, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read command output: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// PrepareFaucet prepares the Bitcoin faucet
func (c *BitcoinContainer) PrepareFaucet(ctx context.Context) (string, error) {
	if _, err := c.ExecuteRpcCommand(ctx, "createwallet", "faucet_wallet"); err != nil {
		return "", err
	}

	preminedAddress, err := c.ExecuteRpcCommandRaw(ctx, []string{"-rpcwallet=faucet_wallet"}, "getnewaddress")
	if err != nil {
		return "", err
	}

	_, err = c.ExecuteRpcCommandRaw(ctx, []string{"-rpcwallet=faucet_wallet"}, "generatetoaddress", "101", preminedAddress)
	if err != nil {
		return "", err
	}

	return preminedAddress, nil
}

// GetFaucetBTC sends Bitcoin from the faucet to the specified address
func (c *BitcoinContainer) GetFaucetBTC(ctx context.Context, address string, amount float64, preminedAddress string) (string, error) {
	if preminedAddress == "" {
		return "", fmt.Errorf("failed to generate pre-mined address")
	}

	txid, err := c.ExecuteRpcCommandRaw(ctx, []string{"-rpcwallet=faucet_wallet"}, "sendtoaddress", address, fmt.Sprintf("%f", amount))
	if err != nil {
		return "", err
	}

	_, err = c.ExecuteRpcCommandRaw(ctx, []string{"-rpcwallet=faucet_wallet"}, "generatetoaddress", "1", preminedAddress)
	if err != nil {
		return "", err
	}

	return txid, nil
}

// MineBlock mines a new block
func (c *BitcoinContainer) MineBlock(ctx context.Context, preminedAddress string) error {
	if preminedAddress == "" {
		return fmt.Errorf("failed to generate pre-mined address")
	}

	_, err := c.ExecuteRpcCommandRaw(ctx, []string{"-rpcwallet=faucet_wallet"}, "generatetoaddress", "1", preminedAddress)
	return err
}
