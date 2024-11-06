package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/novifinancial/serde-reflection/serde-generate/runtime/golang/serde"
	"strings"
	"testing"

	"github.com/holiman/uint256"
	"github.com/rooch-network/rooch-go-sdk/types"

	"github.com/blocktree/go-owcrypt"
)

func TestHttpCall(t *testing.T) {
	client := NewRoochClient("http://localhost:9850")
	var result interface{}

	result, err := client.GetNodeInfo(context.Background())
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)

	result, err = client.GetTransactionByHash(context.Background(), "0x0c8cb10681edff02eb100dba665f8df7452fa30307c20d34d462cf653e3bfefa")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)

	result, err = client.GetTransactionInfoByHash(context.Background(), "0x0c8cb10681edff02eb100dba665f8df7452fa30307c20d34d462cf653e3bfefa")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)

	result, err = client.GetTransactionEventByHash(context.Background(), "0x0c8cb10681edff02eb100dba665f8df7452fa30307c20d34d462cf653e3bfefa")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)

	result, err = client.GetBlockByHash(context.Background(), "0x9e635ae64903409378f5146ff89bfea52a61326ffcbf4191fa63cce642cfc2ea")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)

	result, err = client.GetBlockByNumber(context.Background(), 2)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)

	result, err = client.GetBlocksFromNumber(context.Background(), 2, 10)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)

	result, err = client.ListResource(context.Background(), "0xa76b896725a088beafb470fe93251c4d")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)

	result, err = client.GetState(context.Background(), "0xa76b896725a088beafb470fe93251c4d")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)

	result, err = client.GetGasUnitPrice(context.Background())
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)

	call := ContractCall{
		"0x00000000000000000000000000000001::Token::market_cap",
		[]string{"0x00000000000000000000000000000001::STC::STC"},
		[]string{},
	}
	result, err = client.CallContract(context.Background(), call)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)

}

func TestBalance(t *testing.T) {
	client := NewRoochClient("http://localhost:9850")
	var result *ListResource

	result, err := client.ListResource(context.Background(), "0x79f75dc7cb6812760e1afba01dc9380e")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result.GetBalanceOfStc())
}

func TestWsCall(t *testing.T) {
	client := NewRoochClient("ws://localhost:9870")

	//c, err := client.NewTxnSendRecvEventNotifications("0xb75994d55eae88219dc57e7e62a11bc0")
	c, err := client.NewPendingTransactionsNotifications(context.Background())

	if err != nil {
		t.Error(err)
	}

	data := <-c

	fmt.Println(data)

	c1, err := client.NewTxnSendRecvEventNotifications(context.Background(), "0xb75994d55eae88219dc57e7e62a11bc0")

	if err != nil {
		t.Error(err)
	}

	data1 := <-c1

	fmt.Println(data1)

}

func TestSubmitTransaction(t *testing.T) {
	client := NewRoochClient("http://localhost:9850")
	privateKeyString := "7ddee640acc92417aee935daccfa34306b7c2b827a1308711d5b1d9711e1bdac"
	privateKeyBytes, _ := hex.DecodeString(privateKeyString)
	privateKey := types.Ed25519PrivateKey(privateKeyBytes)
	addressArray, _ := types.ToAccountAddress("b75994d55eae88219dc57e7e62a11bc0")
	toaddr, _ := types.ToAccountAddress("ab4039861ca47ec349b64ddb862293bf")

	result, err := client.TransferStc(context.Background(), *addressArray, privateKey, *toaddr, serde.Uint128{
		High: 0,
		Low:  100000,
	})
	if err != nil {
		t.Errorf("%+v\n", err)
	}
	fmt.Println(result)
}

func TestNodeInfo(t *testing.T) {
	client := NewRoochClient("https://test-seed.rooch.network")
	var result interface{}

	result, err := client.GetNodeInfo(context.Background())
	if err != nil {
		t.Error(fmt.Sprintf("%+v", err))
	}
	j, err := json.Marshal(result)
	if err != nil {
		t.Error(fmt.Sprintf("%+v", err))
	}
	fmt.Println(string(j))
}

func TestSign(t *testing.T) {
	privateKeyString := "587737ebefb4961d377a3ab2f9ceb37b1fa96eb862dfaf954a4a1a99535dfec0"
	publicKeyString := "32ed52d319694aebc5b52e00836e2f7c7d2c7c7791270ede450d21dbc90cbfa1"

	privateKey, _ := hex.DecodeString(privateKeyString)
	publicKey, _ := hex.DecodeString(publicKeyString)

	publicKeyGen, _ := owcrypt.GenPubkey(privateKey, owcrypt.ECC_CURVE_ED25519_NORMAL)

	message := "Example `personal_sign` message"

	signBytes, _, _ := owcrypt.Signature(privateKey, nil, []byte(message), owcrypt.ECC_CURVE_ED25519_NORMAL)
	signString := hex.EncodeToString(signBytes)

	messageBytes, _ := hex.DecodeString("f7abb31497be2d952de2e1c64e2ce3edae7c4d9f5a522386a38af0c76457301eb75994d55eae88219dc57e7e62a11bc0070000000000000002000000000000000000000000000000010f5472616e73666572536372697074730f706565725f746f5f706565725f76320107000000000000000000000000000000010353544303535443000210b75994d55eae88219dc57e7e62a11bc010a0860100000000000000000000000000809698000000000001000000000000000d3078313a3a5354433a3a5354439e1d000000000000fe\n")

	signBytes, _, _ = owcrypt.Signature(privateKey, nil, messageBytes, owcrypt.ECC_CURVE_ED25519_NORMAL)
	signString = hex.EncodeToString(signBytes)

	fmt.Println(owcrypt.Verify(publicKeyGen, nil, messageBytes, signBytes, owcrypt.ECC_CURVE_ED25519_NORMAL))
	fmt.Println(publicKey)
	fmt.Println(publicKeyGen)
	fmt.Println(signString)
}

func TestDeployContract(t *testing.T) {
	client := NewRoochClient("http://localhost:9850")
	privateKeyString := "7ddee640acc92417aee935daccfa34306b7c2b827a1308711d5b1d9711e1bdac"
	privateKeyBytes, _ := hex.DecodeString(privateKeyString)
	privateKey := types.Ed25519PrivateKey(privateKeyBytes)

	//code,_ := ioutil.ReadFile("~/test/resources/contract/MyCounter.mv")
	//fmt.Println(code)
	code := []byte{161, 28, 235, 11, 2, 0, 0, 0, 9, 1, 0, 4, 2, 4, 4, 3, 8, 25, 5, 33, 12, 7, 45, 78, 8, 123, 32, 10, 155, 1, 5, 12, 160, 1, 81, 13, 241, 1, 2, 0, 0, 1, 1, 0, 2, 12, 0, 0, 3, 0, 1, 0, 0, 4, 2, 1, 0, 0, 5, 0, 1, 0, 0, 6, 2, 1, 0, 1, 8, 0, 4, 0, 1, 6, 12, 0, 1, 12, 1, 7, 8, 0, 1, 5, 9, 77, 121, 67, 111, 117, 110, 116, 101, 114, 6, 83, 105, 103, 110, 101, 114, 7, 67, 111, 117, 110, 116, 101, 114, 4, 105, 110, 99, 114, 12, 105, 110, 99, 114, 95, 99, 111, 117, 110, 116, 101, 114, 4, 105, 110, 105, 116, 12, 105, 110, 105, 116, 95, 99, 111, 117, 110, 116, 101, 114, 5, 118, 97, 108, 117, 101, 10, 97, 100, 100, 114, 101, 115, 115, 95, 111, 102, 248, 175, 3, 221, 8, 222, 73, 216, 30, 78, 253, 158, 36, 192, 57, 204, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 2, 1, 7, 3, 0, 1, 0, 1, 0, 3, 13, 11, 0, 17, 4, 42, 0, 12, 1, 10, 1, 16, 0, 20, 6, 1, 0, 0, 0, 0, 0, 0, 0, 22, 11, 1, 15, 0, 21, 2, 1, 2, 0, 1, 0, 1, 3, 14, 0, 17, 0, 2, 2, 1, 0, 0, 1, 5, 11, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 18, 0, 45, 0, 2, 3, 2, 0, 0, 1, 3, 14, 0, 17, 2, 2, 0, 0, 0}

	toAddr, _ := types.ToAccountAddress("b75994d55eae88219dc57e7e62a11bc0")
	moduleId := types.ModuleId{
		Address: *toAddr,
		Name:    "xxxx",
	}

	scriptFunction := types.ScriptFunction{
		Module:   moduleId,
		Function: "init",
		TyArgs:   []types.TypeTag{},
		Args:     [][]byte{},
	}
	toAddr, _ = types.ToAccountAddress("b75994d55eae88219dc57e7e62a11bc0")
	client.DeployContract(context.Background(), *toAddr, privateKey, scriptFunction, code)
}

func TestDryRunRaw(t *testing.T) {
	client := NewRoochClient("http://localhost:9850")
	context := context.Background()

	sender, _ := types.ToAccountAddress("b75994d55eae88219dc57e7e62a11bc0")
	senderPk, _ := HexStringToBytes("a173e69a0e9f87be8179181c3174c7fb1b00eee0955eeece5eb7918bb9a43dbc")
	receiver, _ := types.ToAccountAddress("ab4039861ca47ec349b64ddb862293bf")

	coreAddress, err := hex.DecodeString("00000000000000000000000000000001")
	if err != nil {
		t.Error(err, "decode core address failed")
	}

	var addressArray [16]byte

	copy(addressArray[:], coreAddress[:16])
	coinType := types.StructTag{
		Address: types.AccountAddress(addressArray),
		Module:  types.Identifier("STC"),
		Name:    types.Identifier("STC"),
	}
	payload := Encode_peer_to_peer_v2_script_function(&types.TypeTag__Struct{Value: coinType}, *receiver, serde.Uint128{
		High: 0,
		Low:  10000,
	})

	price, err := client.GetGasUnitPrice(context)
	if err != nil {
		t.Errorf("%+v", err)
	}

	state, err := client.GetState(context, "0x"+hex.EncodeToString(sender[:]))

	if err != nil {
		t.Errorf("%+v", err)
	}

	rawUserTransaction, err := client.BuildRawUserTransaction(context, *sender, payload, price, DefaultMaxGasAmount, state.SequenceNumber)
	if err != nil {
		t.Errorf("%+v", err)
	}

	result, err := client.DryRunRaw(context, *rawUserTransaction, senderPk)
	if err != nil {
		t.Errorf("%+v", err)
	}

	fmt.Println(result)
}

func TestEstimateGas(t *testing.T) {
	//client := NewRoochClient("http://localhost:9850")
	//chainId := 254 //local dev chainId
	client := NewRoochClient("https://test-seed.rooch.network")
	chainId := 251 // test chainId
	context := context.Background()

	sender := "0x569ab535990a17ac9afd1bc57faec683"
	senderPubKey, _ := HexStringToBytes("0xe4cb4052dc3398f3794918f5650fdefb0a5272c4d51220fbf9538ca2c379b00b")
	receiver := "0x17d882a26d86ccb0eedae1bd3db4f47c"
	_ = senderPubKey

	price, err := client.GetGasUnitPrice(context)
	if err != nil {
		t.Errorf("%+v", err)
	}

	// seqNumber, err := client.GetAccountSequenceNumber(context, sender)
	// fmt.Printf("current account seqNumber: %d\n", seqNumber)
	// _ = seqNumber
	// //seqNumber = 1 //note: if seqNumber is incorrent, dry run will fail.
	// if err != nil {
	// 	t.Errorf("%+v", err)
	// }

	code := "0x01::TransferScripts::peer_to_peer_v2"
	typeArgs := []string{"0x01::STC::STC"}
	args := []string{receiver, "1u128"}
	result, err := client.EstimateGas(context, chainId, price, DefaultMaxGasAmount, sender, senderPubKey, nil, code, typeArgs, args)
	if err != nil {
		t.Errorf("%+v", err)
	}

	fmt.Println(result)
}

func TestGetAccountSequenceNumber(t *testing.T) {
	client := NewRoochClient("https://test-seed.rooch.network")
	seqNumber, err := client.GetAccountSequenceNumber(context.Background(), "0x569ab535990a17ac9afd1bc57faec683")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("current account seqNumber: %d\n", seqNumber)
}

func TestGetEvents(t *testing.T) {
	client := NewRoochClient("https://test-seed.rooch.network")
	var to uint64 = 208141
	events, err := client.GetEvents(context.Background(), &EventFilter{
		FromBlock: to - 1,
		ToBlock:   &to,
		TypeTags:  []string{"0xe498d62f5d1f469d2f72eb3e9dc8f230::CrossChainManager::CrossChainEvent"},
	})
	if err != nil {
		t.Error(err)
	}
	j, _ := json.Marshal(events)
	fmt.Println(string(j))
}

func TestGetTransactionProof(t *testing.T) {
	client := NewRoochClient("https://test-seed.rooch.network")
	blockHash := "0x37e8dd4f432a1c3a7b6dfcaa90ebf2aafa0287678ffe4b8ad2373c5b48ffb20c"
	var txGlobalIdx uint64 = 231188
	var eventIndex int = 1
	p, err := client.GetTransactionProof(context.Background(), blockHash, txGlobalIdx, &eventIndex)
	if err != nil {
		t.FailNow()
	}
	fmt.Println("--------------- transaction proof -----------------")
	j, err := json.Marshal(p)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(string(j))
	ebs, err := hexToBytes(p.EventProof.Event)
	if err != nil {
		t.FailNow()
	}
	ev0, err := EventToContractEventV0(ebs)
	if err != nil {
		t.FailNow()
	}
	//fmt.Println(ev0.Key)
	fmt.Println("----------------- event TypeTag ------------------")
	tts, err := getEventTypeTagString(ev0.TypeTag)
	if err != nil { //|| s != "0xe498d62f5d1f469d2f72eb3e9dc8f230::CrossChainManager::CrossChainEvent" {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(tts)
}

func getEventTypeTagString(tt types.TypeTag) (string, error) {
	switch tt := tt.(type) {
	case *types.TypeTag__Struct:
		return "0x" + hex.EncodeToString(tt.Value.Address[:]) + "::" + string(tt.Value.Module) + "::" + string(tt.Value.Name), nil
	default:
		return "", fmt.Errorf("unknown TypeTag type")
	}
}

func TestHeaderByNumber(t *testing.T) {
	client := NewRoochClient("http://localhost:9850")
	hdr, err := client.HeaderByNumber(context.Background(), 1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(hdr)
}

func TestGetEpochResource(t *testing.T) {
	//stateroot := "0x3fda84965ac0b1a9bc70e10e73b56df9a3c5245a639762fa1e4712f4e123a61d"
	//fmt.Println(stateroot)
	client := NewRoochClient("https://test-seed.rooch.network")
	//epochRes, err := client.GetEpochResource(context.Background(), &stateroot)
	epochRes, err := client.GetEpochResourceByHeight(context.Background(), 11)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(epochRes)
	bs, err := json.Marshal(epochRes)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(bs))
}

func TestHeaderWithDifficultyInfoByNumber(t *testing.T) {
	client := NewRoochClient("https://test-seed.rooch.network")
	h, err := client.HeaderWithDifficultyInfoByNumber(context.Background(), 535440)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(h)
	fmt.Println(h.BlockTimeTarget)
	fmt.Println(h.BlockDifficutyWindow)
	bs, err := json.Marshal(h)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(string(bs))
	h2 := &BlockHeaderWithDifficultyInfo{}
	json.Unmarshal(bs, h2)
	fmt.Println(h2)
	fmt.Println(h2.BlockTimeTarget)
	fmt.Println(h2.BlockDifficutyWindow)
	fmt.Println(h2.BlockInfo)
}

func TestGetBlockHeaderAndBlockInfoByNumber(t *testing.T) {
	client := NewRoochClient("https://test-seed.rooch.network")
	hdrAndInfo, err := client.GetBlockHeaderAndBlockInfoByNumber(context.Background(), 175991)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	j, _ := json.Marshal(hdrAndInfo)
	fmt.Println(string(j))

	thdr, _ := hdrAndInfo.BlockHeader.ToTypesHeader()
	tinfo, _ := hdrAndInfo.BlockInfo.ToTypesBlockInfo()
	tHdrAndInfo := types.BlockHeaderAndBlockInfo{
		BlockHeader: *thdr,
		BlockInfo:   *tinfo,
	}
	//fmt.Println(tHdrAndInfo)
	bs, _ := tHdrAndInfo.BcsSerialize()

	tHdrAndInfo2, err := types.BcsDeserializeBlockHeaderAndBlockInfo(bs)
	if err != nil {
		t.FailNow()
	}
	//fmt.Println(tHdrAndInfo2)
	bs2, _ := tHdrAndInfo2.BcsSerialize()

	if !bytes.Equal(bs, bs2) {
		t.FailNow()
	}

	fmt.Println("-------------- difficulty info. ---------------")
	difficulty := new(uint256.Int).SetBytes(thdr.Difficulty[:])
	totalDifficulty := new(uint256.Int).SetBytes(tinfo.TotalDifficulty[:])
	parent, err := client.GetBlockHeaderAndBlockInfoByNumber(context.Background(), thdr.Number-1)
	if err != nil {
		t.FailNow()
	}
	ptinfo, err := parent.BlockInfo.ToTypesBlockInfo()
	if err != nil {
		t.FailNow()
	}
	parentTotalDifficulty := new(uint256.Int).SetBytes(ptinfo.TotalDifficulty[:])
	fmt.Printf("difficulty: %d, totalDifficulty: %d, parentTotalDifficulty: %d\n", difficulty, totalDifficulty, parentTotalDifficulty)
	if totalDifficulty.Cmp(new(uint256.Int).Add(parentTotalDifficulty, difficulty)) != 0 {
		t.FailNow()
	} else {
		//fmt.Println("ok")
	}
}

func TestGetBlockInfoByNumber(t *testing.T) {
	client := NewRoochClient("https://test-seed.rooch.network")
	h, err := client.GetBlockInfoByNumber(context.Background(), 175991)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	j, _ := json.Marshal(h)
	fmt.Println(string(j))
}

func TestGetTransactionInfoByHash(t *testing.T) {
	stcclient := NewRoochClient("https://test-seed.rooch.network")
	txInfo, err := stcclient.GetTransactionInfoByHash(context.Background(), "0xb199fbbf9c7aeef9a0257f0d496ecd0f11ded014526965b4a294c66041272cae")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(strings.EqualFold("\"Executed\"", string(txInfo.Status)))
	fmt.Println(string(txInfo.Status))
}
