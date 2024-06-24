package handler

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/15Andrew43/HLF-project/internal/repository"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"gitlab.com/hlf-mipt/basic-asset-core/pkg/model"
)

func toOwner(clientID, mspID string) string {
	return fmt.Sprintf("%v@%v", clientID, mspID)
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, item *model.Asset) error {
	exists, err := s.AssetExists(ctx, item.ID)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return err
	}
	if exists {
		err = fmt.Errorf("the asset %s already exists", item.ID)
		log.Printf("err: %v", err.Error())
		return err
	}

	mspID, clientID, err := GetClient(ctx)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return err
	}
	item.Owner = toOwner(mspID, clientID)
	item.Type = fmt.Sprint(reflect.TypeOf(*item))

	err = s.repository.Assets().Set(ctx, item)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return err
	}

	return nil
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*model.Asset, error) {
	asset, err := s.repository.Assets().Get(ctx, id)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return nil, err
	}

	granted, err := s.isAccessGranted(ctx, asset)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return nil, err
	}

	if !granted {
		err = errors.New("forbidden")
		log.Printf("err: %v", err.Error())
		return nil, err
	}

	return asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, item *model.Asset) error {
	asset, err := s.ReadAsset(ctx, item.ID)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return err
	}

	err = s.repository.Assets().Set(ctx, asset)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return err
	}

	return nil
}

// DeleteAsset deletes a given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	_, err := s.ReadAsset(ctx, id)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return err
	}

	err = s.repository.Assets().Delete(ctx, id)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return err
	}

	return nil
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	_, err := s.repository.Assets().Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.NotExistsErr) {
			return false, nil
		}
		log.Printf("err: %v", err.Error())
		return false, err
	}

	return true, nil
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) (string, error) {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return "", err
	}

	oldOwner := asset.Owner
	asset.Owner = newOwner

	err = s.repository.Assets().Set(ctx, asset)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return "", err
	}

	return oldOwner, nil
}

func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]model.Asset, error) {
	assets, err := s.repository.Assets().Filter(ctx, map[string]any{
		"type": fmt.Sprint(reflect.TypeOf(new(model.Asset))),
	})
	if err != nil {
		log.Printf("err: %v", err.Error())
		return nil, err
	}

	return assets, nil
}

func (s *SmartContract) QueryAssets(ctx contractapi.TransactionContextInterface, query string) ([]model.Asset, error) {
	assets, err := s.repository.Assets().List(ctx, query)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return nil, err
	}

	result := make([]model.Asset, 0)
	for _, i := range assets {
		granted, err := s.isAccessGranted(ctx, &i)
		if err != nil {
			log.Printf("err: %v", err.Error())
			return nil, err
		}

		if granted {
			result = append(result, i)
		}
	}

	return result, nil
}

func (s *SmartContract) isAccessGranted(ctx contractapi.TransactionContextInterface, item *model.Asset) (bool, error) {
	mspID, clientID, err := GetClient(ctx)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return false, err
	}
	owner := toOwner(mspID, clientID)
	if item.Owner != owner {
		return false, nil
	}

	return true, nil
}
