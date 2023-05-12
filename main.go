package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	numAccounts := 10                       // set the number of accounts to generate
	file, err := os.Create("addresses.txt") // create the output file
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for i := 0; i < numAccounts; i++ {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}

		privateKeyBytes := crypto.FromECDSA(privateKey)
		fmt.Println("Private Key:", hexutil.Encode(privateKeyBytes))

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}

		publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
		fmt.Println("Public Key:", hexutil.Encode(publicKeyBytes))

		address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
		fmt.Println("Address:", address)

		// write the address to the output file
		_, err = file.WriteString(address + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}
