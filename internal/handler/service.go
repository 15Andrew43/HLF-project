package handler

import (
	"github.com/15Andrew43/HLF-project/internal/repository"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
	repository *repository.Repository
}

func NewSmartContract(r *repository.Repository) *SmartContract {
	return &SmartContract{
		repository: r,
	}
}
