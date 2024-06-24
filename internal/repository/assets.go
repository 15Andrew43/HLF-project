package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"gitlab.com/hlf-mipt/basic-asset-core/pkg/model"
)

var NotExistsErr = errors.New("the asset does not exist")

type assetsRepository struct {
}

func newAssetsRepository() *assetsRepository {
	return &assetsRepository{}
}

func (svc *assetsRepository) Set(ctx contractapi.TransactionContextInterface, item *model.Asset) error {
	assetJSON, err := json.Marshal(*item)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(item.ID, assetJSON)
	if err != nil {
		return err
	}

	return nil
}

func (svc *assetsRepository) Get(ctx contractapi.TransactionContextInterface, id string) (*model.Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		err = fmt.Errorf("failed to read from world state: %v", err)
		return nil, err
	}
	if assetJSON == nil {
		err = NotExistsErr
		return nil, err
	}

	var asset model.Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (svc *assetsRepository) Delete(ctx contractapi.TransactionContextInterface, id string) error {
	err := ctx.GetStub().DelState(id)
	if err != nil {
		err = fmt.Errorf("failed to delete from world state: %v", err)
		return nil
	}
	return nil
}

func (svc *assetsRepository) Filter(ctx contractapi.TransactionContextInterface, filter map[string]any) ([]model.Asset, error) {
	f := make([]string, 0)
	for key, value := range filter {
		if fmt.Sprint(reflect.TypeOf(value)) == "string" {
			f = append(f, fmt.Sprintf(`"%v":"%v"`, key, value))
		} else {
			f = append(f, fmt.Sprintf(`"%v":%v`, key, value))
		}
	}

	query := fmt.Sprintf(`"selector": {%v}"`, strings.Join(f, ","))
	return svc.List(ctx, query)
}

func (svc *assetsRepository) List(ctx contractapi.TransactionContextInterface, query string) ([]model.Asset, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []model.Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset model.Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	return assets, nil
}
