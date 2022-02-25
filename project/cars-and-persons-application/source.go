/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	mspID         = "Org1MSP"
	cryptoPath    = "../test-network/organizations/peerOrganizations/org1.example.com"
	certPath      = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/cert.pem"
	keyPath       = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"
	tlsCertPath   = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint  = "localhost:7051"
	gatewayPeer   = "peer0.org1.example.com"
	channelName   = "mychannel"
	chaincodeName = "basic"
)

var now = time.Now()
var assetId = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

func main() {
	log.Println("============ application-golang starts ============")

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gateway.Close()

	network := gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	fmt.Println("Initializing ledger")
	//initLedger(contract)

	var option int

loop:
	for {
		fmt.Println("Choose an option by entering a number:")
		fmt.Println("0 - Initialize ledger")
		fmt.Println("1 - Read person asset")
		fmt.Println("2 - Read car asset")
		fmt.Println("3 - Get cars by color")
		fmt.Println("4 - Get cars by color and owner")
		fmt.Println("5 - Transfer car to another owner")
		fmt.Println("6 - Add car malfunction")
		fmt.Println("7 - Change car color")
		fmt.Println("8 - Repair car")
		fmt.Println("9 - Exit")

		fmt.Scanf("%d", &option)

		switch option {
		case 0:
			fmt.Println("Initializing ledger...")
			initLedger(contract)

		case 1:
			fmt.Printf("Enter person ID: ")
			var personID string
			fmt.Scanf("%s", &personID)
			readPersonAsset(contract, personID)

		case 2:
			fmt.Printf("Enter car ID: ")
			var carID string
			fmt.Scanf("%s", &carID)
			readCarAsset(contract, carID)

		case 3:
			fmt.Printf("Enter car color: ")
			var color string
			fmt.Scanf("%s", &color)
			getCarsByColor(contract, color)

		case 4:
			fmt.Printf("Enter car color: ")
			var color string
			fmt.Scanf("%s", &color)

			fmt.Printf("Enter car owner: ")
			var ownerID string
			fmt.Scanf("%s", &ownerID)
			getCarsByColorAndOwner(contract, color, ownerID)

		case 5:
			fmt.Printf("Enter car ID: ")
			var carID string
			fmt.Scanf("%s", &carID)

			fmt.Printf("Enter new owner ID: ")
			var newOwnerID string
			fmt.Scanf("%s", &newOwnerID)

			fmt.Printf("Does the owner accept malfunctioned car, with a price compensation? (Y/n): ")
			var acceptMalfunctionedStr string
			fmt.Scanf("%s", &acceptMalfunctionedStr)
			var acceptMalfunctionedBool bool
			if acceptMalfunctionedStr == "no" {
				acceptMalfunctionedBool = false
			} else {
				acceptMalfunctionedBool = true
			}

			transferCarAsset(contract, carID, newOwnerID, acceptMalfunctionedBool)

		case 6:
			fmt.Printf("Enter car ID: ")
			var carID string
			fmt.Scanf("%s", &carID)

			fmt.Println("Enter malfunction description:")
			var description string
			fmt.Scanf("%s", &description)

			fmt.Printf("Enter malfunction repair price: ")
			var repairPrice float32
			fmt.Scanf("%f", &repairPrice)
			addCarMalfunction(contract, carID, description, repairPrice)

		case 7:
			fmt.Printf("Enter car ID: ")
			var carID string
			fmt.Scanf("%s", carID)

			fmt.Printf("Enter new car color: ")
			var newColor string
			fmt.Scanf("%s", &newColor)
			changeCarColor(contract, carID, newColor)

		case 8:
			fmt.Printf("Enter car ID: ")
			var carID string
			fmt.Scanf("%s", carID)
			repairCar(contract, carID)

		case 9:
			fmt.Printf("Exiting...")
			break loop

		default:
			fmt.Printf("Invalid input! Please enter a number in the range [1, 9]!")
		}

		fmt.Printf("\n\n")
	}

	log.Println("============ application-golang ends ============")
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	files, err := ioutil.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := ioutil.ReadFile(path.Join(keyPath, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

/*
 This type of transaction would typically only be run once by an application the first time it was started after its
 initial deployment. A new version of the chaincode deployed later would likely not need to run an "init" function.
*/
func initLedger(contract *client.Contract) {
	fmt.Printf("Submit Transaction: InitLedger, function creates the initial set of assets on the ledger \n")

	_, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func readPersonAsset(contract *client.Contract, id string) {
	fmt.Printf("Evaluate Transaction: ReadPersonAsset, function returns person asset attributes\n")

	evaluateResult, err := contract.EvaluateTransaction("ReadPersonAsset", id)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

func readCarAsset(contract *client.Contract, id string) {
	fmt.Printf("Evaluate Transaction: ReadCarAsset, function returns car asset attributes\n")

	evaluateResult, err := contract.EvaluateTransaction("ReadCarAsset", id)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

func getCarsByColor(contract *client.Contract, color string) {
	fmt.Println("Evaluate Transaction: GetCarsByColor, function returns all the cars with the given color")

	evaluateResult, err := contract.EvaluateTransaction("GetCarsByColor", color)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

func getCarsByColorAndOwner(contract *client.Contract, color string, ownerID string) {
	fmt.Println("Evaluate Transaction: GetCarsByColor, function returns all the cars with the given color and owner")

	evaluateResult, err := contract.EvaluateTransaction("GetCarsByColorAndOwner", color, ownerID)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

func transferCarAsset(contract *client.Contract, id string, newOwner string, acceptMalfunction bool) {
	fmt.Printf("Submit Transaction: TransferCarAsset, change car owner \n")

	_, err := contract.SubmitTransaction("TransferCarAsset", id, newOwner, strconv.FormatBool(acceptMalfunction))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func addCarMalfunction(contract *client.Contract, id string, description string, repairPrice float32) {
	fmt.Printf("Submit Transaction: AddCarMalfunction, record a new car malfunction \n")

	_, err := contract.SubmitTransaction("AddCarMalfunction", id, description, fmt.Sprintf("%f", repairPrice))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func changeCarColor(contract *client.Contract, id string, newColor string) {
	fmt.Printf("Submit Transaction: ChangeCarColor, change the color of a car \n")

	_, err := contract.SubmitTransaction("ChangeCarColor", id, newColor)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func repairCar(contract *client.Contract, id string) {
	fmt.Printf("Submit Transaction: RepairCar, fix all of the car's malfunctions \n")

	_, err := contract.SubmitTransaction("RepairCar", id)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

//Format JSON data
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, " ", ""); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
