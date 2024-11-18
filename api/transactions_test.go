package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransaction_UserTransaction(t *testing.T) {
	testJson := `{
  "version": "1010733903",
  "hash": "0xae3f1f751c6cacd61f46054a5e9e39ca9f094802875befbc54ceecbcdf6eff69",
  "state_change_hash": "0x3e8340786d2085a2160fa368c380ed412d4a5a3c5ccad692092c4bc0074fde3e",
  "event_root_hash": "0xe6e2ae41a57d9ab1c7dc58851d7beb4d5be43797ba7225d3e2a3b69c35fe7c2d",
  "state_checkpoint_hash": null,
  "gas_used": "5",
  "success": true,
  "vm_status": "Executed successfully",
  "accumulator_root_hash": "0xf9fdaddf6051311cb54e3756a343faa346f1c9137370762f6eef8e375a7031bb",
  "changes": [
  {
    "address": "0x2932a152328163661f0ae591911270d0edfe0a765beb48a270b9b8a70e766572",
    "state_key_hash": "0xb59b40ac86b159eee2c76ff2eb121b91aa8638ef806d08ed6e061bd60c9b134d",
    "data": {
      "type": "0x1::object::ObjectCore",
      "data": {
        "allow_ungated_transfer": true,
        "guid_creation_num": "1125899906842626",
        "owner": "0x8038df5e61a19a5f86ad01f4389736b08250dad1b4aa864afc4fc639a2581ca8",
        "transfer_events": {
          "counter": "1",
          "guid": {
            "id": {
              "addr": "0x2932a152328163661f0ae591911270d0edfe0a765beb48a270b9b8a70e766572",
              "creation_num": "1125899906842624"
            }
          }
        }
      }
    },
    "type": "write_resource"
  },
  {
    "address": "0x2932a152328163661f0ae591911270d0edfe0a765beb48a270b9b8a70e766572",
    "state_key_hash": "0xb59b40ac86b159eee2c76ff2eb121b91aa8638ef806d08ed6e061bd60c9b134d",
    "data": {
      "type": "0x3::gas_coin::RGas",
      "data": {
        "burn_ref": {
          "vec": [
          {
            "inner": {
              "vec": [
              {
                "self": "0x2932a152328163661f0ae591911270d0edfe0a765beb48a270b9b8a70e766572"
              }
              ]
            },
            "self": {
              "vec": []
            }
          }
          ]
        },
        "mutator_ref": {
          "vec": [
          {
            "self": "0x2932a152328163661f0ae591911270d0edfe0a765beb48a270b9b8a70e766572"
          }
          ]
        },
        "property_mutator_ref": {
          "self": "0x2932a152328163661f0ae591911270d0edfe0a765beb48a270b9b8a70e766572"
        },
        "transfer_ref": {
          "vec": [
          {
            "self": "0x2932a152328163661f0ae591911270d0edfe0a765beb48a270b9b8a70e766572"
          }
          ]
        }
      }
    },
    "type": "write_resource"
  },
  {
    "address": "0x2932a152328163661f0ae591911270d0edfe0a765beb48a270b9b8a70e766572",
    "state_key_hash": "0xb59b40ac86b159eee2c76ff2eb121b91aa8638ef806d08ed6e061bd60c9b134d",
    "data": {
      "type": "0x4::property_map::PropertyMap",
      "data": {
        "inner": {
	      "data": []
        }
      }
    },
    "type": "write_resource"
  },
  {
    "address": "0x2932a152328163661f0ae591911270d0edfe0a765beb48a270b9b8a70e766572",
    "state_key_hash": "0xb59b40ac86b159eee2c76ff2eb121b91aa8638ef806d08ed6e061bd60c9b134d",
    "data": {
      "type": "0x4::token::Token",
      "data": {
        "collection": {
          "inner": "0x778adb39026a14009cf5aa93eb53d81299e40c7a8dbcdbf7b490cbc29749d259"
        },
        "description": "This is BLACK FLAG ARMY NFT",
        "index": "0",
        "mutation_events": {
          "counter": "0",
          "guid": {
            "id": {
              "addr": "0x2932a152328163661f0ae591911270d0edfe0a765beb48a270b9b8a70e766572",
              "creation_num": "1125899906842625"
            }
          }
        },
        "name": "",
        "uri": "https://bafybeierhssqdg7fv64xkkjuvsq4bikj2yfmuxm4dvb6jxb2un4yw37ohi.ipfs.w3s.link/68.webp"
      }
    },
    "type": "write_resource"
  }
  ],
  "sender": "0xa46c6c7a65d605685e23055a6a906fb7284ba87849cbeb579d5c07424938241e",
  "sequence_number": "242217",
  "max_gas_amount": "2018",
  "gas_unit_price": "100",
  "expiration_timestamp_secs": "1719968695",
  "payload": {
    "function": "0x1::object::transfer",
    "type_arguments": [
      "0x4::token::Token"
    ],
    "arguments": [
    {
      "inner": "0x2932a152328163661f0ae591911270d0edfe0a765beb48a270b9b8a70e766572"
    },
    "0x8038df5e61a19a5f86ad01f4389736b08250dad1b4aa864afc4fc639a2581ca8"
    ],
    "type": "entry_function_payload"
  },
  "signature": {
    "public_key": "0x5e10e3db4e3c700142b9a3e18c40038db5903f2dedfe41d09aca74a8c68565d6",
    "signature": "0xa95686dab2c93cf1720e300b929e3656cc6cdc3a8389dc12bb9bd5a17ae3af975bee9d618f080266e3a60f1e2968220a83d773e2b3902edfe54127ed0a7b290b",
    "type": "ed25519_signature"
  },
  "events": [
  {
    "guid": {
      "creation_number": "1125899906842624",
      "account_address": "0x2932a152328163661f0ae591911270d0edfe0a765beb48a270b9b8a70e766572"
    },
    "sequence_number": "0",
    "type": "0x1::object::TransferEvent",
    "data": {
      "from": "0xa46c6c7a65d605685e23055a6a906fb7284ba87849cbeb579d5c07424938241e",
      "object": "0x2932a152328163661f0ae591911270d0edfe0a765beb48a270b9b8a70e766572",
      "to": "0x8038df5e61a19a5f86ad01f4389736b08250dad1b4aa864afc4fc639a2581ca8"
    }
  },
  {
    "guid": {
      "creation_number": "0",
      "account_address": "0x0"
    },
    "sequence_number": "0",
    "type": "0x1::transaction_fee::FeeStatement",
    "data": {
      "execution_gas_units": "3",
      "io_gas_units": "2",
      "storage_fee_octas": "0",
      "storage_fee_refund_octas": "0",
      "total_charge_gas_units": "5"
    }
  }
  ],
  "timestamp": "1719965096135309",
  "type": "user_transaction"
}`
	data := &Transaction{}
	err := json.Unmarshal([]byte(testJson), &data)
	assert.NoError(t, err)
	assert.Equal(t, TransactionVariantUser, data.Type)
	data2 := &CommittedTransaction{}
	err = json.Unmarshal([]byte(testJson), &data2)
	assert.NoError(t, err)
	assert.Equal(t, TransactionVariantUser, data2.Type)

	txn, err := data.UserTransaction()
	assert.NoError(t, err)
	txn2, err := data2.UserTransaction()
	assert.NoError(t, err)
	assert.Equal(t, txn, txn2)

	assert.Equal(t, uint64(1010733903), txn.Version)
	assert.Equal(t, uint64(1719965096135309), txn.Timestamp)
	assert.Equal(t, uint64(242217), txn.SequenceNumber)
	assert.Equal(t, uint64(100), txn.GasUnitPrice)
	assert.Equal(t, uint64(2018), txn.MaxGasAmount)
	assert.Equal(t, uint64(1719968695), txn.ExpirationTimestampSecs)

	// TODO: test some more

	// Check functions
	assert.Equal(t, *data.Version(), data2.Version())
	assert.Equal(t, data.Hash(), data2.Hash())
	assert.Equal(t, *data.Success(), data2.Success())
}
