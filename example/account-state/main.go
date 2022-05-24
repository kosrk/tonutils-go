package main

import (
	"context"
	"fmt"
	"log"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
)

func main() {
	client := liteclient.NewClient()

	// connect to mainnet lite server
	err := client.Connect(context.Background(), "135.181.140.212:13206", "K0t3+IWLOXHYMvMcrGZDPs+pn58a17LFbnXoQkKc2xw=")
	if err != nil {
		log.Fatalln("connection err: ", err.Error())
		return
	}

	// initialize ton api lite connection wrapper
	api := ton.NewAPIClient(client)

	// we need fresh block info to run get methods
	b, err := api.GetBlockInfo(context.Background())
	if err != nil {
		log.Fatalln("get block err:", err.Error())
		return
	}

	addr := address.MustParseAddr("EQAoUyP1KBBRvTVAUxlAI_9mmSH05guWrNZ5PfmVFL7zs2b6")

	res, err := api.GetAccount(context.Background(), b, addr)
	if err != nil {
		log.Fatalln("get account err:", err.Error())
		return
	}

	fmt.Printf("Status: %s\n", res.State.Status)
	fmt.Printf("Balance: %s TON\n", res.State.Balance.TON())
	if res.Data != nil {
		fmt.Printf("Data: %s\n", res.Data.Dump())
	}

	// take last tx info from account info
	lastHash := res.LastTxHash
	lastLt := res.LastTxLT

	fmt.Printf("\nTransactions:\n")
	for {
		// last transaction has 0 prev lt
		if lastLt == 0 {
			break
		}

		// load transactions in batches with size 15
		list, err := api.ListTransactions(context.Background(), addr, 15, lastLt, lastHash)
		if err != nil {
			log.Printf("send err: %s", err.Error())
			return
		}

		// oldest = first in list
		for _, t := range list {
			fmt.Println(t.String())
		}

		// set previous info from the oldest transaction in list
		lastHash = list[0].PrevTxHash
		lastLt = list[0].PrevTxLT
	}
}
