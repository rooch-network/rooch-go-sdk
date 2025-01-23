// Copyright (c) RoochNetwork
// SPDX-License-Identifier: Apache-2.0

package container

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"time"

	"github.com/testcontainers/testcontainers-go/wait"
)

const ROOCH_PORT = 6767

type RoochContainer struct {
	testcontainers.Container
	networkName          string
	dataDir              string
	accountDir           string
	port                 int
	ethRpcUrl            string
	btcRpcUrl            string
	btcRpcUsername       string
	btcRpcPassword       string
	btcEndBlockHeight    *int
	btcSyncBlockInterval *int
	hostConfigPath       string
	trafficBurstSize     *int
	trafficPerSecond     *int
	image                string
}

func NewRoochContainer(image string) *RoochContainer {
	if image == "" {
		image = "ghcr.io/rooch-network/rooch:main_debug"
	}

	return &RoochContainer{
		networkName: "local",
		dataDir:     "TMP",
		accountDir:  "/root/.rooch",
		port:        ROOCH_PORT,
		image:       image,
	}
}

// Builder pattern methods
func (r *RoochContainer) WithNetworkName(name string) *RoochContainer {
	r.networkName = name
	return r
}

func (r *RoochContainer) WithDataDir(dir string) *RoochContainer {
	r.dataDir = dir
	return r
}

func (r *RoochContainer) WithPort(port int) *RoochContainer {
	r.port = port
	return r
}

func (r *RoochContainer) WithEthRpcUrl(url string) *RoochContainer {
	r.ethRpcUrl = url
	return r
}

func (r *RoochContainer) WithBtcRpcUrl(url string) *RoochContainer {
	r.btcRpcUrl = url
	return r
}

func (r *RoochContainer) WithBtcRpcUsername(username string) *RoochContainer {
	r.btcRpcUsername = username
	return r
}

func (r *RoochContainer) WithBtcRpcPassword(password string) *RoochContainer {
	r.btcRpcPassword = password
	return r
}

func (r *RoochContainer) WithBtcEndBlockHeight(height int) *RoochContainer {
	r.btcEndBlockHeight = &height
	return r
}

func (r *RoochContainer) WithBtcSyncBlockInterval(interval int) *RoochContainer {
	r.btcSyncBlockInterval = &interval
	return r
}

func (r *RoochContainer) WithHostConfigPath(hostPath string) *RoochContainer {
	r.hostConfigPath = hostPath
	return r
}

func (r *RoochContainer) WithTrafficBurstSize(burstSize int) *RoochContainer {
	r.trafficBurstSize = &burstSize
	return r
}

func (r *RoochContainer) WithTrafficPerSecond(perSecond int) *RoochContainer {
	r.trafficPerSecond = &perSecond
	return r
}

func (r *RoochContainer) InitializeRooch(ctx context.Context) error {
	if r.hostConfigPath == "" {
		return fmt.Errorf("host config path not set. Call WithHostConfigPath() before initializing")
	}

	//req := testcontainers.ContainerRequest{
	//	Image:      r.image,
	//	Cmd:        []string{"init", "--skip-password"},
	//	Mounts:     testcontainers.Mounts(testcontainers.BindMount(r.hostConfigPath, r.accountDir)),
	//	WaitingFor: wait.ForLog("JSON-RPC HTTP Server start listening"),
	//}

	//req := testcontainers.ContainerRequest{
	//	Image: r.image,
	//	Cmd:   []string{"init", "--skip-password"},
	//	Mounts: []testcontainers.ContainerMount{
	//		{
	//			Source: testcontainers.GenericBindMountSource{
	//				HostPath: r.hostConfigPath,
	//			},
	//			Target: r.accountDir,
	//		},
	//	},
	//	WaitingFor: wait.ForLog("JSON-RPC HTTP Server start listening").WithStartupTimeout(10 * time.Second),
	//}

	req := testcontainers.ContainerRequest{
		Image: "ghcr.io/rooch-network/rooch:main_debug",
		Cmd:   []string{"init", "--skip-password"},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      r.hostConfigPath,
				ContainerFilePath: r.accountDir,
				FileMode:          0644,
			},
		},
		WaitingFor: wait.ForLog("JSON-RPC HTTP Server start listening").WithStartupTimeout(10 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return err
	}
	defer container.Terminate(ctx)

	// Switch to local environment
	req.Cmd = []string{"env", "switch", "--alias", "local"}
	container, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return err
	}
	defer container.Terminate(ctx)

	return nil
}

func (r *RoochContainer) Start(ctx context.Context) (*StartedRoochContainer, error) {
	if r.hostConfigPath == "" {
		return nil, fmt.Errorf("host config path not set. Call WithHostConfigPath() before starting")
	}

	cmd := []string{
		"server",
		"start",
		"-n",
		r.networkName,
		"-d",
		r.dataDir,
		"--port",
		fmt.Sprintf("%d", r.port),
	}

	if r.ethRpcUrl != "" {
		cmd = append(cmd, "--eth-rpc-url", r.ethRpcUrl)
	}
	if r.btcRpcUrl != "" {
		cmd = append(cmd, "--btc-rpc-url", r.btcRpcUrl)
	}
	if r.btcRpcUsername != "" {
		cmd = append(cmd, "--btc-rpc-username", r.btcRpcUsername)
	}
	if r.btcRpcPassword != "" {
		cmd = append(cmd, "--btc-rpc-password", r.btcRpcPassword)
	}
	if r.btcEndBlockHeight != nil {
		cmd = append(cmd, "--btc-end-block-height", fmt.Sprintf("%d", *r.btcEndBlockHeight))
	}
	if r.btcSyncBlockInterval != nil {
		cmd = append(cmd, "--btc-sync-block-interval", fmt.Sprintf("%d", *r.btcSyncBlockInterval))
	}
	if r.trafficPerSecond != nil {
		cmd = append(cmd, "--traffic-per-second", fmt.Sprintf("%d", *r.trafficPerSecond))
	}
	if r.trafficBurstSize != nil {
		cmd = append(cmd, "--traffic-burst-size", fmt.Sprintf("%d", *r.trafficBurstSize))
	}

	//req := testcontainers.ContainerRequest{
	//	Image:        r.image,
	//	ExposedPorts: []string{fmt.Sprintf("%d", r.port)},
	//	Cmd:          cmd,
	//	User:         "root",
	//	Mounts:       testcontainers.Mounts(testcontainers.BindMount(r.hostConfigPath, r.accountDir)),
	//	WaitingFor:   wait.ForLog("JSON-RPC HTTP Server start listening"),
	//}
	req := testcontainers.ContainerRequest{
		Image:        r.image,
		ExposedPorts: []string{fmt.Sprintf("%d", r.port)},
		Cmd:          cmd,
		User:         "root",
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      r.hostConfigPath,
				ContainerFilePath: r.accountDir,
				FileMode:          0644,
			},
		},
		WaitingFor: wait.ForLog("JSON-RPC HTTP Server start listening").WithStartupTimeout(120 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	//containerPort := nat.Port(fmt.Sprintf("%d/tcp", r.port))
	containerPort := nat.Port(fmt.Sprintf("%d/tcp", r.port))
	mappedPort, err := container.MappedPort(ctx, containerPort) // Updated to use containerPort
	if err != nil {
		return nil, err
	}

	return &StartedRoochContainer{
		Container:            container,
		networkName:          r.networkName,
		dataDir:              r.dataDir,
		containerPort:        r.port,
		mappedPort:           mappedPort.Int(),
		ethRpcUrl:            r.ethRpcUrl,
		btcRpcUrl:            r.btcRpcUrl,
		btcRpcUsername:       r.btcRpcUsername,
		btcRpcPassword:       r.btcRpcPassword,
		btcEndBlockHeight:    r.btcEndBlockHeight,
		btcSyncBlockInterval: r.btcSyncBlockInterval,
		trafficBurstSize:     r.trafficBurstSize,
		trafficPerSecond:     r.trafficPerSecond,
	}, nil
}

type StartedRoochContainer struct {
	testcontainers.Container
	networkName          string
	dataDir              string
	containerPort        int
	mappedPort           int
	ethRpcUrl            string
	btcRpcUrl            string
	btcRpcUsername       string
	btcRpcPassword       string
	btcEndBlockHeight    *int
	btcSyncBlockInterval *int
	trafficBurstSize     *int
	trafficPerSecond     *int
}

// Getter methods
func (s *StartedRoochContainer) GetPort() int {
	return s.mappedPort
}

func (s *StartedRoochContainer) GetNetworkName() string {
	return s.networkName
}

func (s *StartedRoochContainer) GetDataDir() string {
	return s.dataDir
}

func (s *StartedRoochContainer) GetEthRpcUrl() string {
	return s.ethRpcUrl
}

func (s *StartedRoochContainer) GetBtcRpcUrl() string {
	return s.btcRpcUrl
}

func (s *StartedRoochContainer) GetBtcRpcUsername() string {
	return s.btcRpcUsername
}

func (s *StartedRoochContainer) GetBtcRpcPassword() string {
	return s.btcRpcPassword
}

func (s *StartedRoochContainer) GetBtcEndBlockHeight() *int {
	return s.btcEndBlockHeight
}

func (s *StartedRoochContainer) GetBtcSyncBlockInterval() *int {
	return s.btcSyncBlockInterval
}

func (s *StartedRoochContainer) GetTrafficBurstSize() *int {
	return s.trafficBurstSize
}

func (s *StartedRoochContainer) GetTrafficPerSecond() *int {
	return s.trafficPerSecond
}

func (s *StartedRoochContainer) GetConnectionAddress(ctx context.Context) (string, error) {
	host, err := s.Host(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", host, s.GetPort()), nil
}
