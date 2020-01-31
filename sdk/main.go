package main

import (
	"fmt"
	// "github.com/chainHero/heroes-service/blockchain"
	"github.com/off-grid-block/vote/sdk/blockchain"
	"os"
)

func main() {
	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		// Network parameters 
		OrdererID: "orderer.example.com",

		// Channel parameters
		ChannelID:     "mychannel",
		ChannelConfig: "/Users/brianli/playground/fabric/fabric-samples/first-network/channel-artifacts/channel.tx",

		// Chaincode parameters
		ChainCodeID:     "vote",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/off-grid-block/off-grid-cc/vote",
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",

		// User parameters
		UserName: "User1",
	}

	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}
	// Close SDK
	defer fSetup.CloseSDK()	
}