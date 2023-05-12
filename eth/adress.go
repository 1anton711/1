package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	var i int
	for {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}
		address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
		if strings.HasPrefix(address, "0x") {
			count := 0
			for _, c := range address[2:] {
				if c == rune(address[2]) && c == rune(address[3]) && c == rune(address[4]) {

					count++
					if count == 4 {
						file, err := os.OpenFile("adresses.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						if err != nil {
							log.Fatal(err)
						}
						defer file.Close()
						privateBytes := crypto.FromECDSA(privateKey)
						publicBytes := crypto.FromECDSAPub(publicKeyECDSA)
						line := fmt.Sprintf("Address: %s\nPrivate key: %s\nPublic key: %s\n\n", address, hexutil.Encode(privateBytes), hexutil.Encode(publicBytes))
						if _, err := file.WriteString(line); err != nil {
							log.Fatal(err)
						}
						fmt.Println("Address found:", address)
						i++
						if i == 10 {
							return
						}
					}
				} else {
					count = 0
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
