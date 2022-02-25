# pdasp-2021

This is an example application implemented using the Hyperledger Fabric blockchain framework.

The "./project" directory mainly contains the code from the official Hyperledger Fabric samples (https://hyperledger-fabric.readthedocs.io/en/release-2.2/install.html). Some of the code has been modified to fit the purposes of this project.

# Using the application

## Starting the network
To start the network, enter the project/test-network directory and run "./network.sh createChannel -ca". This will run the network with 3 organizations, each containing 3 peers. The option "-ca" ensures that all the cryptomaterial will be generated using the Certificate Authority. If this option is not specified, the cryptogen tool will be used for this purpose which will not start the network properly.

## Deploying the chaincode
To deploy the chaincode after the network is started, run "./network.sh deployCC -ccn basic -ccp ../cars-and-persons-chaincodes/ -ccl go", while in the project/test-network directory. This will package, install and commit the chaincode on all 9 peers in the network (3 peers in 3 organizations). After deploying the chaincode, it will be approved by all 3 organizations. The chaincode name is "basic". The source code for the smart contract that makes up the chaincode is written in Golang and is located in the project/cars-and-persons-chaincodes directory.

### Endorsement policy
The endorsement policy was changed from MAJORITY to a Signature endorsement policy type, which specifies that a transaction must be endorsed by at least one peer per organization, that is all 3 organizations participate in trasaction endorsement. This is a more strict policy, since the MAJORITY policy accepts a transaction when the majority of the organizations endorses it, which, in this case, means that an endorsement from two out of three organizations would be enough.

## Running the client application
To run the client application, enter the project/cars-and-persons-appliction, and run "go run source.go". After this, follow the instructions from the console in order to interact with the network. The application communicates with one of the peers from the organization that you choose after starting the application. In order to change which peer this is, open the project/cars-and-persons-application/app_config.json and change the desired fields.