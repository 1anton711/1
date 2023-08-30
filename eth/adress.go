package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/97d8b38f52c741968acb32ec13c31a31")
	if err != nil {
		log.Fatal(err)
	}

	var privateKeyBatch []*ecdsa.PrivateKey
	for i := 0; i < 1000; i++ {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}
		privateKeyBatch = append(privateKeyBatch, privateKey)
	}

	publicKeyBatch := make([]*ecdsa.PublicKey, len(privateKeyBatch))
	for i, privateKey := range privateKeyBatch {
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("error casting public key to ECDSA")
		}
		publicKeyBatch[i] = publicKeyECDSA
	}

	value := big.NewFloat(1)
	etherValue, _ := new(big.Float).Mul(value, big.NewFloat(1e18)).Int(nil)
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x7B5C3eEC47d52D8d17dF51D2d3a4Cf5f957aD1D0")
	var data []byte

	//Gas U

	var wg sync.WaitGroup
	for {
		start := time.Now()
		for i, publicKeyECDSA := range publicKeyBatch {
			fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
			fmt.Printf("%d. fromAddress: %s\n", i+1, fromAddress.Hex())

			wg.Add(1)
			go func() {
				defer wg.Done()

				balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("balance: %s\n", balance.String())

				if balance.Cmp(big.NewInt(0)) > 0 {
					nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
					if err != nil {
						log.Fatal(err)
					}

					tx := types.NewTransaction(nonce, toAddress, etherValue, gasLimit, gasPrice, data)

					chainID, err := client.NetworkID(context.Background())
					if err != nil {
						log.Fatal(err)
					}

					signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyBatch[i])
					if err != nil {
						log.Fatal(err)
					}

					err = client.SendTransaction(context.Background(), signedTx)
					if err != nil {
						log.Fatal(err)
					}

					fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
				}
			}()
		}

		wg.Wait()
		fmt.Printf("time elapsed: %v\n", time.Since(start))
		time.Sleep(1 * time.Millisecond)
	}
}
