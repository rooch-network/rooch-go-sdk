// Copyright (c) RoochNetwork
// SPDX-License-Identifier: Apache-2.0

package testsuite

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rooch-network/rooch-go-sdk/testsuite/container"

	//container "github.com/rooch-network/rooch-go-sdk/testsuite/container"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	ordNetworkAlias     = "ord"
	bitcoinNetworkAlias = "bitcoind"
)

var defaultCmdAddress string

// TestBox represents the main testing environment
type TestBox struct {
	TmpDir string
	//Network *Network
	//OrdContainer     *StartedOrdContainer
	BitcoinContainer *container.StartedBitcoinContainer
	RoochContainer   *container.StartedRoochContainer // Can be either *StartedRoochContainer or int (PID)
	RoochDir         string
	RoochPort        int
	miningTicker     *time.Ticker
	done             chan bool
}

// NewTestBox creates and initializes a new TestBox instance
func NewTestBox() (*TestBox, error) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "rooch_test_*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	tb := &TestBox{
		TmpDir:   tmpDir,
		RoochDir: filepath.Join(tmpDir, ".rooch_test"),
		done:     make(chan bool),
	}

	// Create rooch directory
	if err := os.MkdirAll(tb.RoochDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create rooch dir: %w", err)
	}

	// Initialize rooch configuration
	if result, err := tb.roochCommand([]string{"init", "--config-dir", tb.RoochDir, "--skip-password"}); err != nil {
		fmt.Printf("roochCommand init exce result: %v\n", result)
		return nil, err
	}

	if result, err := tb.roochCommand([]string{"env", "switch", "--config-dir", tb.RoochDir, "--alias", "local"}); err != nil {
		fmt.Printf("roochCommand switch exce result: %v\n", result)
		return nil, err
	}

	return tb, nil
}

// LoadBitcoinEnv initializes the Bitcoin environment
func (tb *TestBox) LoadBitcoinEnv(customContainer *container.BitcoinContainer, autoMining bool) error {
	var err error
	if customContainer != nil {
		tb.BitcoinContainer = customContainer
		err = customContainer.Start(context.Background())
	} else {
		bitcoinContainer, err := container.NewBitcoinContainer()
		if err != nil {
			return fmt.Errorf("failed to create bitcoin container: %w", err)
		}
		bitcoinContainer.WithHostDataPath(tb.TmpDir)
		//bitcoinContainer.WithNetwork(tb.getNetwork())
		//bitcoinContainer.SetNetworkAliases(bitcoinNetworkAlias)

		tb.BitcoinContainer = bitcoinContainer
		err = bitcoinContainer.Start(context.Background())
	}
	if err != nil {
		return fmt.Errorf("failed to start bitcoin container: %w", err)
	}

	// Wait for container to be ready
	time.Sleep(5 * time.Second)

	if autoMining {
		// Prepare Faucet
		if err := tb.BitcoinContainer.PrepareFaucet(); err != nil {
			return err
		}

		// Start mining blocks periodically
		tb.miningTicker = time.NewTicker(1 * time.Second)
		go func() {
			for {
				select {
				case <-tb.miningTicker.C:
					if tb.BitcoinContainer != nil {
						if err := tb.BitcoinContainer.MineBlock(); err != nil {
							fmt.Printf("Error mining block: %v\n", err)
						}
					}
				case <-tb.done:
					return
				}
			}
		}()
	}

	return nil
}

// LoadRoochEnv initializes the Rooch environment
func (tb *TestBox) LoadRoochEnv(target interface{}, port int) error {
	if port == 0 {
		port = 6767
	}

	// Handle different target types
	switch v := target.(type) {
	case *container.RoochContainer:
		//container, err := v.Start()
		_container, err := v.Start(context.Background())
		if err != nil {
			return err
		}
		tb.RoochContainer = container
		return nil
	case string:
		if v == "local" {
			var err error
			if port == 0 {
				port, err = GetUnusedPort()
				if err != nil {
					return err
				}
			}

			// Generate a random port for metrics
			metricsPort, err := GetUnusedPort()
			if err != nil {
				return err
			}

			cmds := []string{
				"server", "start",
				"-n", "local",
				"-d", "TMP",
				"--port", strconv.Itoa(port),
			}

			if tb.BitcoinContainer != nil {
				cmds = append(cmds,
					"--btc-rpc-url", tb.BitcoinContainer.GetRpcUrl(),
					"--btc-rpc-username", tb.BitcoinContainer.GetRpcUser(),
					"--btc-rpc-password", tb.BitcoinContainer.GetRpcPass(),
					"--btc-sync-block-interval", "1",
				)
			}

			cmds = append(cmds,
				"--traffic-per-second", "1",
				"--traffic-burst-size", "5000",
			)

			envs := []string{fmt.Sprintf("METRICS_HOST_PORT=%d", metricsPort)}
			pid, err := tb.roochAsyncCommand(cmds, fmt.Sprintf("JSON-RPC HTTP Server start listening 0.0.0.0:%d", port), envs)
			if err != nil {
				return err
			}

			tb.RoochContainer = pid
			tb.RoochPort = port
			return nil
		}
	}

	return fmt.Errorf("unsupported target type")
}

// CleanEnv cleans up all resources
func (tb *TestBox) CleanEnv() {
	if tb.miningTicker != nil {
		tb.miningTicker.Stop()
		tb.done <- true
	}

	if tb.BitcoinContainer != nil {
		tb.BitcoinContainer.Stop()
	}

	switch v := tb.RoochContainer.(type) {
	case int:
		process, err := os.FindProcess(v)
		if err == nil {
			process.Kill()
		}
	case *container.StartedRoochContainer:
		v.Stop()
	}

	os.RemoveAll(tb.TmpDir)
}

// DefaultCmdAddress retrieves the default account address
func (tb *TestBox) DefaultCmdAddress() (string, error) {
	if defaultCmdAddress == "" {
		output, err := tb.roochCommand([]string{"account", "list", "--config-dir", tb.RoochDir, "--json"})
		if err != nil {
			return "", err
		}

		var accounts interface{}
		if err := json.Unmarshal([]byte(output), &accounts); err != nil {
			return "", err
		}

		// Handle different account response formats
		switch v := accounts.(type) {
		case []interface{}:
			for _, acc := range v {
				if account, ok := acc.(map[string]interface{}); ok {
					if active, ok := account["active"].(bool); ok && active {
						if localAccount, ok := account["local_account"].(map[string]interface{}); ok {
							defaultCmdAddress = localAccount["hex_address"].(string)
							break
						}
					}
				}
			}
		case map[string]interface{}:
			if defaultAcc, ok := v["default"].(map[string]interface{}); ok {
				defaultCmdAddress = defaultAcc["hex_address"].(string)
			}
		}

		if defaultCmdAddress == "" {
			return "", fmt.Errorf("no active account address")
		}
	}

	return defaultCmdAddress, nil
}

// GetRoochServerAddress returns the Rooch server address
func (tb *TestBox) GetRoochServerAddress() string {
	switch v := tb.RoochContainer.(type) {
	case *container.StartedRoochContainer:
		return v.GetConnectionAddress()
	default:
		if tb.RoochPort != 0 {
			return fmt.Sprintf("127.0.0.1:%d", tb.RoochPort)
		}
	}
	return "127.0.0.1:6767"
}

// GetFaucetBTC retrieves BTC from the faucet
func (tb *TestBox) GetFaucetBTC(address string, amount float64) (string, error) {
	if tb.BitcoinContainer == nil {
		return "", fmt.Errorf("bitcoin container not started")
	}

	if amount == 0 {
		amount = 0.001
	}

	return tb.BitcoinContainer.GetFaucetBTC(address, amount)
}

// Helper functions

// GetUnusedPort returns an available port number
func GetUnusedPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}

// findRootDir finds the root directory containing the target file
func (tb *TestBox) findRootDir(targetName string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		targetPath := filepath.Join(currentDir, targetName)
		if _, err := os.Stat(targetPath); err == nil {
			return currentDir, nil
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break
		}
		currentDir = parent
	}

	return "", fmt.Errorf("root directory not found")
}

// roochCommand executes a Rooch command synchronously
func (tb *TestBox) roochCommand(args []string, envs ...string) (string, error) {
	root, err := tb.findRootDir("pnpm-workspace.yaml")
	if err != nil {
		return "", err
	}

	roochDir := filepath.Join(root, "target", "debug")
	cmdArgs := []string{filepath.Join(roochDir, "rooch")}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = append(os.Environ(), envs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %w, output: %s", err, string(output))
	}

	return string(output), nil
}

// roochAsyncCommand executes a Rooch command asynchronously
func (tb *TestBox) roochAsyncCommand(args []string, waitFor string, envs []string) (int, error) {
	root, err := tb.findRootDir("pnpm-workspace.yaml")
	if err != nil {
		return 0, err
	}

	roochDir := filepath.Join(root, "target", "debug")
	cmdArgs := []string{filepath.Join(roochDir, "rooch")}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = append(os.Environ(), envs...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}

	if err := cmd.Start(); err != nil {
		return 0, err
	}

	// Read output until we see the waitFor string
	buf := make([]byte, 1024)
	for {
		n, err := stdout.Read(buf)
		if err != nil {
			cmd.Process.Kill()
			return 0, err
		}

		if strings.Contains(string(buf[:n]), waitFor) {
			return cmd.Process.Pid, nil
		}
	}
}
