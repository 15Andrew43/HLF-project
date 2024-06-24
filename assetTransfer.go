package main

import (
	"log"

	"github.com/15Andrew43/HLF-project/internal/handler"
	"github.com/15Andrew43/HLF-project/internal/repository"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	repo := repository.NewRepository()
	assetChaincode, err := contractapi.NewChaincode(handler.NewSmartContract(repo))
	if err != nil {
		log.Printf("error creating asset-transfer-basic chaincode: %v", err)
		return
	}

	if err := assetChaincode.Start(); err != nil {
		log.Printf("error starting asset-transfer-basic chaincode: %v", err)
	}
}
