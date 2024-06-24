package handler

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

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
	item.Owner = toOwner(clientID, mspID)
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
	mspID, _, err := GetClient(ctx)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return nil, err
	}

	asset, err := s.repository.Assets().Get(ctx, id)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return nil, err
	}

	if !s.isSameMSP(asset.Owner, mspID) {
		err = errors.New("access denied: only participants of the same MSP can read this asset")
		log.Printf("err: %v", err.Error())
		return nil, err
	}

	return asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, item *model.Asset) error {
	mspID, clientID, err := GetClient(ctx)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return err
	}

	asset, err := s.ReadAsset(ctx, item.ID)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return err
	}

	if !s.isOwner(asset.Owner, clientID, mspID) {
		err = errors.New("access denied: only the creator can update their asset")
		log.Printf("err: %v", err.Error())
		return err
	}

	err = s.repository.Assets().Set(ctx, item)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return err
	}

	return nil
}

// DeleteAsset deletes a given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	mspID, clientID, err := GetClient(ctx)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return err
	}

	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return err
	}

	if !s.isOwner(asset.Owner, clientID, mspID) {
		err = errors.New("access denied: only the creator can delete their asset")
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
	mspID, clientID, err := GetClient(ctx)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return "", err
	}

	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return "", err
	}

	if !s.isOwner(asset.Owner, clientID, mspID) {
		err = errors.New("access denied: only the creator can transfer their asset")
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
	mspID, _, err := GetClient(ctx)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return nil, err
	}

	assets, err := s.repository.Assets().Filter(ctx, map[string]any{
		"type": fmt.Sprint(reflect.TypeOf(new(model.Asset))),
	})
	if err != nil {
		log.Printf("err: %v", err.Error())
		return nil, err
	}

	filteredAssets := make([]model.Asset, 0)
	for _, asset := range assets {
		if s.isSameMSP(asset.Owner, mspID) {
			filteredAssets = append(filteredAssets, asset)
		}
	}

	return filteredAssets, nil
}

func (s *SmartContract) QueryAssets(ctx contractapi.TransactionContextInterface, query string) ([]model.Asset, error) {
	mspID, _, err := GetClient(ctx)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return nil, err
	}

	assets, err := s.repository.Assets().List(ctx, query)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return nil, err
	}

	filteredAssets := make([]model.Asset, 0)
	for _, asset := range assets {
		if s.isSameMSP(asset.Owner, mspID) {
			filteredAssets = append(filteredAssets, asset)
		}
	}

	return filteredAssets, nil
}

func (s *SmartContract) isAccessGranted(ctx contractapi.TransactionContextInterface, item *model.Asset) (bool, error) {
	mspID, clientID, err := GetClient(ctx)
	if err != nil {
		log.Printf("err: %v", err.Error())
		return false, err
	}
	return s.isOwner(item.Owner, clientID, mspID), nil
}

func (s *SmartContract) isOwner(owner, clientID, mspID string) bool {
	return owner == toOwner(clientID, mspID)
}

func (s *SmartContract) isSameMSP(owner, mspID string) bool {
	return strings.HasSuffix(owner, "@"+mspID)
}
