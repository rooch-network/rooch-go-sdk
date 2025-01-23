module github.com/rooch-network/rooch-go-sdk

go 1.22.0

//toolchain go1.23.3
toolchain go1.22.7

require (
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/ethereum/go-ethereum v1.14.12
	github.com/stretchr/testify v1.10.0
	github.com/testcontainers/testcontainers-go v0.35.0
	github.com/tyler-smith/go-bip32 v1.0.0
	github.com/tyler-smith/go-bip39 v1.1.0
	golang.org/x/crypto v0.32.0
)

require (
	github.com/btcsuite/btcd/btcec/v2 v2.3.4
	//github.com/btcsuite/btcd/btcutil v1.1.5 // indirect
	github.com/btcsuite/btcd/btcutil v1.1.6
	github.com/btcsuite/btcd/chaincfg/chainhash v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/crypto/blake256 v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	dario.cat/mergo v1.0.1 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/FactomProject/basen v0.0.0-20150613233007-fe3947df716e // indirect
	github.com/FactomProject/btcutilecc v0.0.0-20130527213604-d3a63a5752ec // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/containerd/containerd v1.7.25 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/containerd/platforms v0.2.1 // indirect
	github.com/cpuguy83/dockercfg v0.3.2 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/docker v27.5.0+incompatible // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/lufia/plan9stats v0.0.0-20240909124753-873cd0166683 // indirect
	github.com/magiconair/properties v1.8.9 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/patternmatcher v0.6.0 // indirect
	github.com/moby/sys/sequential v0.6.0 // indirect
	github.com/moby/sys/user v0.3.0 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/shirou/gopsutil/v3 v3.24.5 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.9.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.58.0 // indirect
	go.opentelemetry.io/otel v1.33.0 // indirect
	go.opentelemetry.io/otel/metric v1.33.0 // indirect
	go.opentelemetry.io/otel/trace v1.33.0 // indirect
)

//require (
//	github.com/cucumber/godog v0.14.1
//	github.com/ethereum/go-ethereum v1.14.5
//	github.com/hasura/go-graphql-client v0.12.1
//	github.com/hdevalence/ed25519consensus v0.2.0
//	github.com/stretchr/testify v1.9.0
//	golang.org/x/crypto v0.24.0
//)
//
//require (
//	filippo.io/edwards25519 v1.1.0 // indirect
//	github.com/btcsuite/btcd/btcec/v2 v2.3.3 // indirect
//	github.com/cucumber/gherkin/go/v26 v26.2.0 // indirect
//	github.com/cucumber/messages/go/v21 v21.0.1 // indirect
//	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
//	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
//	github.com/gofrs/uuid v4.3.1+incompatible // indirect
//	github.com/google/uuid v1.6.0 // indirect
//	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
//	github.com/hashicorp/go-memdb v1.3.4 // indirect
//	github.com/hashicorp/golang-lru v0.5.4 // indirect
//	github.com/holiman/uint256 v1.2.4 // indirect
//	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
//	github.com/rogpeppe/go-internal v1.12.0 // indirect
//	github.com/spf13/pflag v1.0.5 // indirect
//	golang.org/x/sys v0.21.0 // indirect
//	gopkg.in/yaml.v3 v3.0.1 // indirect
//	nhooyr.io/websocket v1.8.11 // indirect
//)
