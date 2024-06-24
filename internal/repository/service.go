package repository

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"gitlab.com/hlf-mipt/basic-asset-core/pkg/model"
)

type AssetsRepository interface {
	Set(ctx contractapi.TransactionContextInterface, item *model.Asset) error
	Get(ctx contractapi.TransactionContextInterface, id string) (*model.Asset, error)
	Delete(ctx contractapi.TransactionContextInterface, id string) error
	Filter(ctx contractapi.TransactionContextInterface, filter map[string]any) ([]model.Asset, error)
	List(ctx contractapi.TransactionContextInterface, query string) ([]model.Asset, error)
}

type Repository struct {
	assets *assetsRepository
}

func NewRepository() *Repository {
	return &Repository{
		assets: newAssetsRepository(),
	}
}

func (svc *Repository) Assets() AssetsRepository {
	return svc.assets
}
