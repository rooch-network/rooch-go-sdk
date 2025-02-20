# rooch-go-sdk

# Rooch Go SDK

Rooch Go SDK implementation.

## 用户指南

### 引入依赖

``` 
import "github.com/rooch-network/rooch-go-sdk"
```

### 节点配置

TODO

默认的 http 端口是 6677 ，websocket 端口是 9870。

## Examples

### 签名并发送交易

#### 转账

```
	client := NewRoochClient("http://localhost:9850")
	privateKeyString := "7ddee640acc92417aee935daccfa34306b7c2b827a1308711d5b1d9711e1bdac"
	privateKeyBytes, _ := hex.DecodeString(privateKeyString)
	privateKey := types.Ed25519PrivateKey(privateKeyBytes)
	addressArray := types.ToAccountAddress("b75994d55eae88219dc57e7e62a11bc0")

	result, err := client.Transfer(addressArray, privateKey, types.ToAccountAddress("ab4039861ca47ec349b64ddb862293bf"), serde.Uint128{
		High: 0,
		Low:  100000,
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result)

```

#### 部署合约

```
	client := NewRoochClient("http://localhost:9850")
	privateKeyString := "7ddee640acc92417aee935daccfa34306b7c2b827a1308711d5b1d9711e1bdac"
	privateKeyBytes, _ := hex.DecodeString(privateKeyString)
	privateKey := types.Ed25519PrivateKey(privateKeyBytes)

	//code,_ := ioutil.ReadFile("~/test/resources/contract/MyCounter.mv") 
	//fmt.Println(code)
	code := []byte{161, 28, 235, 11, 2, 0, 0, 0, 9, 1, 0, 4, 2, 4, 4, 3, 8, 25, 5, 33, 12, 7, 45, 78, 8, 123, 32, 10, 155, 1, 5, 12, 160, 1, 81, 13, 241, 1, 2, 0, 0, 1, 1, 0, 2, 12, 0, 0, 3, 0, 1, 0, 0, 4, 2, 1, 0, 0, 5, 0, 1, 0, 0, 6, 2, 1, 0, 1, 8, 0, 4, 0, 1, 6, 12, 0, 1, 12, 1, 7, 8, 0, 1, 5, 9, 77, 121, 67, 111, 117, 110, 116, 101, 114, 6, 83, 105, 103, 110, 101, 114, 7, 67, 111, 117, 110, 116, 101, 114, 4, 105, 110, 99, 114, 12, 105, 110, 99, 114, 95, 99, 111, 117, 110, 116, 101, 114, 4, 105, 110, 105, 116, 12, 105, 110, 105, 116, 95, 99, 111, 117, 110, 116, 101, 114, 5, 118, 97, 108, 117, 101, 10, 97, 100, 100, 114, 101, 115, 115, 95, 111, 102, 248, 175, 3, 221, 8, 222, 73, 216, 30, 78, 253, 158, 36, 192, 57, 204, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 2, 1, 7, 3, 0, 1, 0, 1, 0, 3, 13, 11, 0, 17, 4, 42, 0, 12, 1, 10, 1, 16, 0, 20, 6, 1, 0, 0, 0, 0, 0, 0, 0, 22, 11, 1, 15, 0, 21, 2, 1, 2, 0, 1, 0, 1, 3, 14, 0, 17, 0, 2, 2, 1, 0, 0, 1, 5, 11, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 18, 0, 45, 0, 2, 3, 2, 0, 0, 1, 3, 14, 0, 17, 2, 2, 0, 0, 0}

	moduleId := types.ModuleId{
		types.ToAccountAddress("b75994d55eae88219dc57e7e62a11bc0"),
		"xxxx",
	}

	scriptFunction := types.ScriptFunction{
		moduleId,
		"init",
		[]types.TypeTag{},
		[][]byte{},
	}
	client.DeployContract(types.ToAccountAddress("b75994d55eae88219dc57e7e62a11bc0"), privateKey, scriptFunction, code)
```

### 根据地址查询链上最新状态或者资源

#### 查询资源

```
	client := NewRoochClient("http://localhost:9850")

	result, err := client.ListResource("0xa76b896725a088beafb470fe93251c4d")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)

	result, err = client.GetState("0xa76b896725a088beafb470fe93251c4d")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)

```

#### 查询最新状态

```
	client := NewRoochClient("http://localhost:9850")

	result, err := client.GetState("0xa76b896725a088beafb470fe93251c4d")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)
```

#### 查询txn

```
	client := NewRoochClient("http://localhost:9850")
	
	var result interface{}

	result, err = client.GetTransactionByHash("0x0c8cb10681edff02eb100dba665f8df7452fa30307c20d34d462cf653e3bfefa")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)

	result, err = client.GetTransactionInfoByHash("0x0c8cb10681edff02eb100dba665f8df7452fa30307c20d34d462cf653e3bfefa")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)

```

## License

rooch-go-sdk is licensed as [Apache 2.0](./LICENSE).
