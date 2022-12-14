package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"

	"github.com/gorilla/websocket"
)

// var baseUrl = "https://cloudflare-eth.com"
var baseUrl = "https://mainnet.infura.io/v3/164ee0f7f9fa4f94905c5f11e90ec1e9"
var wsUrl = "wss://mainnet.infura.io/ws/v3/164ee0f7f9fa4f94905c5f11e90ec1e9"
var usdtToken = "0xdAC17F958D2ee523a2206206994597C13D831ec7"
var myAddr = "0x64d17712Fc5795e0784bbA04f8dB83Fe16E23f17"

func main() {
	client, err := ethclient.Dial(baseUrl)
	if err != nil {
		panic(err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	val := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_subscribe",
		"params": []interface{}{
			"newPendingTransactions",
		},
	}
	if err := conn.WriteJSON(val); err != nil {
		panic(err)
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			panic(err)
		}

		resp := &PendindResp{}
		resp.Params = &ParamResp{}
		if err := json.Unmarshal(msg, &resp); err != nil {
			panic(err)
		}

		hash := common.HexToHash(resp.Params.Result)
		tx, isPending, err := client.TransactionByHash(context.Background(), hash)
		if err != nil {
			fmt.Println(err)
		}

		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		if isPending {
			//fmt.Println(string(msg))
			fmt.Println("isPending :", isPending, ", hash :", hash.Hex())

			i := new(big.Int)

			if msgTx, err := tx.AsMessage(types.NewEIP155Signer(chainID), i); err == nil {
				fmt.Println(msgTx.To().Hex())
				fmt.Println(msgTx.From().Hex())
				fmt.Println(string(msgTx.Data()))
			}

			//fmt.Println(msg)

		}

		fmt.Println()

	}
}

type PendindResp struct {
	Params *ParamResp `json:"params"`
}

type ParamResp struct {
	Result string `json:"result"`
}

/*
func main() {
  val := map[string]interface{}{
    "jsonrpc": "2.0",
    "method":  "eth_getBalance",
    "params": []string{
      "0xF7931B9b1fFF5Fc63c45577C43DFc0D0dEf16C46",
      "latest",
    },
    "id": 1,
  }

  payloadBuf := new(bytes.Buffer)
  json.NewEncoder(payloadBuf).Encode(val)

  req, err := http.NewRequest(http.MethodPost, baseUrl, payloadBuf)
  if err != nil {
    panic(err)
  }
  req.Header.Add("Content-Type", "application/json; charset=UTF-8")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()

  data := &RespData{}
  if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
    panic(err)
  }

  fbalance := new(big.Float)
  fbalance.SetString(data.Result)
  ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

  fmt.Println(ethValue)
}

type RespData struct {
  Result string json:"result"
}
*/
