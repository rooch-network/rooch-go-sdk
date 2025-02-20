package api

import (
	"encoding/json"
	"fmt"
	"github.com/rooch-network/rooch-go-sdk/address"
	"github.com/rooch-network/rooch-go-sdk/types"
	"github.com/rooch-network/rooch-go-sdk/utils"
)

// TransactionVariant is the type of transaction, all transactions submitted by this SDK are [TransactionVariantUser]
type TransactionVariant string

const (
	TransactionVariantPending         TransactionVariant = "pending_transaction"          // TransactionVariantPending maps to PendingTransaction
	TransactionVariantUser            TransactionVariant = "user_transaction"             // TransactionVariantUser maps to UserTransaction
	TransactionVariantGenesis         TransactionVariant = "genesis_transaction"          // TransactionVariantGenesis maps to GenesisTransaction
	TransactionVariantBlockMetadata   TransactionVariant = "block_metadata_transaction"   // TransactionVariantBlockMetadata maps to BlockMetadataTransaction
	TransactionVariantBlockEpilogue   TransactionVariant = "block_epilogue_transaction"   // TransactionVariantBlockEpilogue maps to BlockEpilogueTransaction
	TransactionVariantStateCheckpoint TransactionVariant = "state_checkpoint_transaction" // TransactionVariantStateCheckpoint maps to StateCheckpointTransaction
	TransactionVariantValidator       TransactionVariant = "validator_transaction"        // TransactionVariantValidator maps to ValidatorTransaction
	TransactionVariantUnknown         TransactionVariant = "unknown"                      // TransactionVariantUnknown maps to UnknownTransaction for unknown types
)

// CommittedTransaction is an enum type for all possible committed transactions on the blockchain
// This is the same as [Transaction] but with the Success and Version functions always confirmed.
type CommittedTransaction struct {
	Type  TransactionVariant // Type of the transaction
	Inner TransactionImpl    // Inner is the actual transaction
}

// Hash of the transaction for lookup on-chain
func (o *CommittedTransaction) Hash() Hash {
	return o.Inner.TxnHash()
}

// Success of the transaction.  Pending transactions, and genesis may not have a success field.
// If this is the case, it will be nil
func (o *CommittedTransaction) Success() bool {
	return *o.Inner.TxnSuccess()
}

// Version of the transaction on chain, will be nil if it is a PendingTransaction
func (o *CommittedTransaction) Version() uint64 {
	return *o.Inner.TxnVersion()
}

// UnmarshalJSON unmarshals the [Transaction] from JSON handling conversion between types
func (o *CommittedTransaction) UnmarshalJSON(b []byte) error {
	type inner struct {
		Type string `json:"type"`
	}
	data := &inner{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	o.Type = TransactionVariant(data.Type)
	switch o.Type {
	case TransactionVariantPending:
		return fmt.Errorf("transaction type is not committed: %s, this is unexpected for the API to return", o.Type)
	case TransactionVariantUser:
		o.Inner = &UserTransaction{}
	default:
		o.Inner = &UnknownTransaction{Type: string(o.Type)}
		o.Type = TransactionVariantUnknown
		return json.Unmarshal(b, &o.Inner.(*UnknownTransaction).Payload)
	}
	return json.Unmarshal(b, o.Inner)
}

// UserTransaction changes the transaction to a [UserTransaction]; however, it will fail if it's not one.
func (o *CommittedTransaction) UserTransaction() (*UserTransaction, error) {
	if o.Type == TransactionVariantUser {
		return o.Inner.(*UserTransaction), nil
	}
	return nil, fmt.Errorf("transaction type is not user: %s", o.Type)
}

// UnknownTransaction changes the transaction to a [UnknownTransaction]; however, it will fail if it's not one.
func (o *CommittedTransaction) UnknownTransaction() (*UnknownTransaction, error) {
	if o.Type == TransactionVariantUnknown {
		return o.Inner.(*UnknownTransaction), nil
	}
	return nil, fmt.Errorf("transaction type is not unknown: %s", o.Type)
}

// Transaction is an enum type for all possible transactions on the blockchain
type Transaction struct {
	Type  TransactionVariant // Type of the transaction
	Inner TransactionImpl    // Inner is the actual transaction
}

// Hash of the transaction for lookup on-chain
func (o *Transaction) Hash() Hash {
	return o.Inner.TxnHash()
}

// Success of the transaction.  Pending transactions, and genesis may not have a success field.
// If this is the case, it will be nil
func (o *Transaction) Success() *bool {
	return o.Inner.TxnSuccess()
}

// Version of the transaction on chain, will be nil if it is a PendingTransaction
func (o *Transaction) Version() *uint64 {
	return o.Inner.TxnVersion()
}

// UnmarshalJSON unmarshals the [Transaction] from JSON handling conversion between types
func (o *Transaction) UnmarshalJSON(b []byte) error {
	type inner struct {
		Type string `json:"type"`
	}
	data := &inner{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	o.Type = TransactionVariant(data.Type)
	switch o.Type {
	case TransactionVariantUser:
		o.Inner = &UserTransaction{}
	default:
		o.Inner = &UnknownTransaction{Type: string(o.Type)}
		o.Type = TransactionVariantUnknown
		return json.Unmarshal(b, &o.Inner.(*UnknownTransaction).Payload)
	}
	return json.Unmarshal(b, o.Inner)
}

// UserTransaction changes the transaction to a [UserTransaction]; however, it will fail if it's not one.
func (o *Transaction) UserTransaction() (*UserTransaction, error) {
	if o.Type == TransactionVariantUser {
		return o.Inner.(*UserTransaction), nil
	}
	return nil, fmt.Errorf("transaction type is not user: %s", o.Type)
}

// TransactionImpl is an interface for all transactions
type TransactionImpl interface {
	// TxnSuccess tells us if the transaction is a success.  It will be nil if the transaction is not committed.
	TxnSuccess() *bool

	// TxnHash gives us the hash of the transaction.
	TxnHash() Hash

	// TxnVersion gives us the ledger version of the transaction. It will be nil if the transaction is not committed.
	TxnVersion() *uint64
}

// UnknownTransaction is a transaction type that is not recognized by the SDK
type UnknownTransaction struct {
	Type    string         // Type is the type of the unknown transaction
	Payload map[string]any // Payload is the raw JSON payload
}

// TxnSuccess tells us if the transaction is a success.  It will be nil if the transaction is not committed.
func (u *UnknownTransaction) TxnSuccess() *bool {
	success := u.Payload["success"]
	if success == nil {
		return nil
	}
	successBool := success.(bool)
	return &successBool
}

// TxnHash gives us the hash of the transaction.
func (u *UnknownTransaction) TxnHash() Hash {
	return u.Payload["hash"].(string)
}

// TxnVersion gives us the ledger version of the transaction. It will be nil if the transaction is not committed.
func (u *UnknownTransaction) TxnVersion() *uint64 {
	versionStr := u.Payload["version"].(string)
	num, err := utils.StrToUint64(versionStr)
	if err != nil {
		return nil
	} else {
		return &num
	}
}

// UserTransaction is a user submitted transaction as an entry function, script, or more.
//
// These transactions are the only transactions submitted by users to the blockchain.
type UserTransaction struct {
	Version                 uint64                  // Version of the transaction, starts at 0 and increments per transaction.
	Hash                    Hash                    // Hash of the transaction, it is a SHA3-256 hash in hexadecimal format with a leading 0x.
	AccumulatorRootHash     Hash                    // AccumulatorRootHash of the transaction.
	StateChangeHash         Hash                    // StateChangeHash of the transaction.
	EventRootHash           Hash                    // EventRootHash of the transaction.
	GasUsed                 uint64                  // GasUsed by the transaction, will be in gas units.
	Success                 bool                    // Success of the transaction.
	VmStatus                string                  // VmStatus of the transaction, this will contain the error if any.
	Changes                 []*WriteSetChange       // Changes to the ledger from the transaction, should never be empty.
	Events                  []*Event                // Events emitted by the transaction, may be empty.
	Sender                  *address.AccountAddress // Sender of the transaction, will never be nil.
	SequenceNumber          uint64                  // SequenceNumber of the transaction, starts at 0 and increments per transaction submitted by the sender.
	MaxGasAmount            uint64                  // MaxGasAmount of the transaction, this is the max amount of gas units that the user is willing to pay.
	GasUnitPrice            uint64                  // GasUnitPrice of the transaction, this is the multiplier per unit of gas to tokens.
	ExpirationTimestampSecs uint64                  // ExpirationTimestampSecs of the transaction, this is the Unix timestamp in seconds when the transaction expires.
	Payload                 *TransactionPayload     // Payload of the transaction, this is the actual transaction data.
	Signature               *Signature              // Signature is the AccountAuthenticator of the sender.
	Timestamp               uint64                  // Timestamp is the Unix timestamp in microseconds when the block of the transaction was committed.
	StateCheckpointHash     Hash                    // StateCheckpointHash of the transaction. Optional, and will be "" if not set.
}

// TxnHash gives us the hash of the transaction.
func (o *UserTransaction) TxnHash() Hash {
	return o.Hash
}

// TxnSuccess tells us if the transaction is a success.  It will never be nil.
func (o *UserTransaction) TxnSuccess() *bool {
	return &o.Success
}

// TxnVersion gives us the ledger version of the transaction. It will never be nil.
func (o *UserTransaction) TxnVersion() *uint64 {
	return &o.Version
}

// UnmarshalJSON unmarshals the [UserTransaction] from JSON handling conversion between types
func (o *UserTransaction) UnmarshalJSON(b []byte) error {
	type inner struct {
		Version                 types.U64               `json:"version"`
		Hash                    Hash                    `json:"hash"`
		AccumulatorRootHash     Hash                    `json:"accumulator_root_hash"`
		StateChangeHash         Hash                    `json:"state_change_hash"`
		EventRootHash           Hash                    `json:"event_root_hash"`
		GasUsed                 types.U64               `json:"gas_used"`
		Success                 bool                    `json:"success"`
		VmStatus                string                  `json:"vm_status"`
		Changes                 []*WriteSetChange       `json:"changes"`
		Events                  []*Event                `json:"events"`
		Sender                  *address.AccountAddress `json:"sender"`
		SequenceNumber          types.U64               `json:"sequence_number"`
		MaxGasAmount            types.U64               `json:"max_gas_amount"`
		GasUnitPrice            types.U64               `json:"gas_unit_price"`
		ExpirationTimestampSecs types.U64               `json:"expiration_timestamp_secs"`
		Payload                 *TransactionPayload     `json:"payload"`
		Signature               *Signature              `json:"signature"`
		Timestamp               types.U64               `json:"timestamp"`
		StateCheckpointHash     Hash                    `json:"state_checkpoint_hash"` // Optional
	}
	data := &inner{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	o.Version = data.Version.ToUint64()
	o.Hash = data.Hash
	o.AccumulatorRootHash = data.AccumulatorRootHash
	o.StateChangeHash = data.StateChangeHash
	o.EventRootHash = data.EventRootHash
	o.GasUsed = data.GasUsed.ToUint64()
	o.Success = data.Success
	o.VmStatus = data.VmStatus
	o.Changes = data.Changes
	o.Events = data.Events
	o.Sender = data.Sender
	o.SequenceNumber = data.SequenceNumber.ToUint64()
	o.MaxGasAmount = data.MaxGasAmount.ToUint64()
	o.GasUnitPrice = data.GasUnitPrice.ToUint64()
	o.ExpirationTimestampSecs = data.ExpirationTimestampSecs.ToUint64()
	o.Payload = data.Payload
	o.Signature = data.Signature
	o.Timestamp = data.Timestamp.ToUint64()
	o.StateCheckpointHash = data.StateCheckpointHash
	return nil
}

// SubmitTransactionResponse is the response from submitting a transaction to the blockchain, it is the same
// as a [PendingTransaction]
type SubmitTransactionResponse = PendingTransaction

// BatchSubmitTransactionResponse is the response from submitting a batch of transactions to the blockchain
type BatchSubmitTransactionResponse struct {
	// TransactionFailures is the list of transactions that failed to submit, if it is empty, all were successful
	TransactionFailures []BatchSubmitTransactionFailure `json:"transaction_failures"`
}
