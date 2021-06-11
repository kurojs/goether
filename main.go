package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/k0kubun/pp"
)

func main() {
	rpcClient, err := rpc.DialContext(context.Background(), os.Getenv("ETH_URL_RPC"))
	if err != nil {
		log.Fatal(err)
	}

	params := []interface{}{
		"newPendingTransactions",
	}

	txHashCh := make(chan string)
	sub, err := rpcClient.EthSubscribe(context.Background(), txHashCh, params...)
	if err != nil {
		log.Fatal(err)
	}

	client, err := ethclient.Dial(os.Getenv("ETH_URL"))
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case txHash := <-txHashCh:
			tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(txHash))
			if err != nil {
				fmt.Print(err)
			}
			pp.Println(tx)
		}
	}
}
