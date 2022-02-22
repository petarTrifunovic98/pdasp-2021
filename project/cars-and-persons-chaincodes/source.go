package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type CarMalfunction struct {
	Description string
	RepairPrice float32
}

type CarAsset struct {
	ID              string
	Brand           string
	Model           string
	Year            int
	Color           string
	OwnerID         string
	MalfunctionList []CarMalfunction
}

type PersonAsset struct {
	ID                 string
	FirstName          string
	LastName           string
	EmailAddress       string
	AmountOfMoneyOwned float32
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	carAssets := []CarAsset{
		{ID: "car1", Brand: "Opel", Model: "Cascada", Year: 2013, Color: "blue", OwnerID: "person1", MalfunctionList: []CarMalfunction{
			{Description: "Shakey steering wheel", RepairPrice: 50},
			{Description: "Oil leaking", RepairPrice: 75},
		}},
		{ID: "car2", Brand: "Audi", Model: "A4", Year: 2016, Color: "red", OwnerID: "person2", MalfunctionList: []CarMalfunction{
			{Description: "Flat fron left tyre", RepairPrice: 15},
		}},
		{ID: "car3", Brand: "Volvo", Model: "V60", Year: 2014, Color: "green", OwnerID: "person1", MalfunctionList: []CarMalfunction{
			{Description: "Cracked windscreen", RepairPrice: 100},
			{Description: "Loose back wiper", RepairPrice: 5},
		}},
		{ID: "car4", Brand: "Zastava", Model: "Yugo 45", Year: 1985, Color: "yellow", OwnerID: "person1", MalfunctionList: []CarMalfunction{
			{Description: "Broken alternator", RepairPrice: 80},
			{Description: "Broken spark plug", RepairPrice: 70},
			{Description: "Loose exhaust pipe", RepairPrice: 10},
			{Description: "Overheating", RepairPrice: 120},
		}},
		{ID: "car5", Brand: "Mercedes-Benz", Model: "A-class", Year: 2018, Color: "black", OwnerID: "person3", MalfunctionList: []CarMalfunction{}},
		{ID: "car6", Brand: "BMW", Model: "X5", Year: 2018, Color: "white", OwnerID: "person2", MalfunctionList: []CarMalfunction{
			{Description: "Cracked headlight", RepairPrice: 30},
		}},
	}

	personAssets := []PersonAsset{
		{ID: "person1", FirstName: "Petar", LastName: "Trifunovic", EmailAddress: "petar@pdasp.rs", AmountOfMoneyOwned: 332.54},
		{ID: "person2", FirstName: "Marko", LastName: "Markovic", EmailAddress: "marko@pdasp.rs", AmountOfMoneyOwned: 567.4},
		{ID: "person3", FirstName: "Jovana", LastName: "Jovanovic", EmailAddress: "jovana@pdasp.rs", AmountOfMoneyOwned: 143.22},
	}

	for _, carAsset := range carAssets {
		carAssetJSON, err := json.Marshal(carAsset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(carAsset.ID, carAssetJSON)
		if err != nil {
			return fmt.Errorf("failed to put cars to world state. %v", err)
		}
	}

	for _, personAsset := range personAssets {
		personAssetJSON, err := json.Marshal(personAsset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(personAsset.ID, personAssetJSON)
		if err != nil {
			return fmt.Errorf("failed to put persons to world state. %v", err)
		}
	}

	return nil
}

func (s *SmartContract) ReadPersonAsset(ctx contractapi.TransactionContextInterface, id string) (*PersonAsset, error) {
	personAssetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read person from world state: %v", err)
	}
	if personAssetJSON == nil {
		return nil, fmt.Errorf("the person asset %s does not exist", id)
	}

	var personAsset PersonAsset
	err = json.Unmarshal(personAssetJSON, &personAsset)
	if err != nil {
		return nil, err
	}

	return &personAsset, nil
}

func (s *SmartContract) ReadCarAsset(ctx contractapi.TransactionContextInterface, id string) (*CarAsset, error) {
	carAssetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read car from world state: %v", err)
	}
	if carAssetJSON == nil {
		return nil, fmt.Errorf("the car asset %s does not exist", id)
	}

	var carAsset CarAsset
	err = json.Unmarshal(carAssetJSON, &carAsset)
	if err != nil {
		return nil, err
	}

	return &carAsset, nil
}

func (s *SmartContract) TransferCarAsset(ctx contractapi.TransactionContextInterface, id string, newOwnerID string) (string, error) {
	carAsset, err := s.ReadCarAsset(ctx, id)
	if err != nil {
		return "", err
	}

	exists, err := s.PersonAssetExists(ctx, newOwnerID)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("the person %s does not exist", newOwnerID)
	}

	oldOwner := carAsset.OwnerID
	carAsset.OwnerID = newOwnerID

	carAssetJSON, err := json.Marshal(carAsset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, carAssetJSON)
	if err != nil {
		return "", err
	}

	return oldOwner, nil
}

func (s *SmartContract) AddCarMalfunction(ctx contractapi.TransactionContextInterface, id string, description string, repairPrice float32) error {
	carAsset, err := s.ReadCarAsset(ctx, id)
	if err != nil {
		return err
	}

	newMalfunction := CarMalfunction{
		Description: description,
		RepairPrice: repairPrice,
	}

	carAsset.MalfunctionList = append(carAsset.MalfunctionList, newMalfunction)

	carAssetJSON, err := json.Marshal(carAsset)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, carAssetJSON)
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) ChangeCarColor(ctx contractapi.TransactionContextInterface, id string, newColor string) (string, error) {
	carAsset, err := s.ReadCarAsset(ctx, id)
	if err != nil {
		return "", err
	}

	oldColor := carAsset.Color
	carAsset.Color = newColor

	carAssetJSON, err := json.Marshal(carAsset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, carAssetJSON)
	if err != nil {
		return "", err
	}

	return oldColor, nil
}

func (s *SmartContract) RepairCar(ctx contractapi.TransactionContextInterface, id string) error {
	carAsset, err := s.ReadCarAsset(ctx, id)
	if err != nil {
		return err
	}

	personAsset, err := s.ReadPersonAsset(ctx, carAsset.OwnerID)
	if err != nil {
		return err
	}

	repairPriceSum := float32(0)
	for _, carMalfunction := range carAsset.MalfunctionList {
		repairPriceSum += carMalfunction.RepairPrice
		if repairPriceSum > personAsset.AmountOfMoneyOwned {
			return fmt.Errorf("The owner of the car cannot afford to pay the car repair price")
		}
	}

	carAsset.MalfunctionList = []CarMalfunction{}
	personAsset.AmountOfMoneyOwned -= repairPriceSum

	carAssetJSON, err := json.Marshal(carAsset)
	if err != nil {
		return err
	}

	personAssetJSON, err := json.Marshal(personAsset)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, carAssetJSON)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(personAsset.ID, personAssetJSON)
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) PersonAssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	personAssetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read person asset from world state: %v", err)
	}

	return personAssetJSON != nil, nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating cars-and-persons chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting cars-and-persons chaincode: %v", err)
	}
}
