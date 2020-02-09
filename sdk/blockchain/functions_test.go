
package blockchain

import (
	"testing"
)

func TestSDKChaincodeFunctions(t *testing.T) {

	// Set up test SDK parameters
	fSetup := SetupSDK {
		OrdererID: "orderer.example.com",

		ChannelID: 			"mychannel",
		ChannelConfig: 		"/Users/brianli/deon/fabric-samples/first-network/channel-artifacts/channel.tx",

		ChainCodeID: 		"vote",
		ChaincodeGoPath: 	"/Users/brianli/deon",
		ChaincodePath:   	"vote/cc",
		OrgAdmin:        	"Admin",
		OrgName:         	"org1",
		ConfigFile:      	"/Users/brianli/deon/src/vote/sdk/config.yaml",

		UserName: "User1",
	}

	err := fSetup.Initialization()
	if err != nil {
		t.Errorf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}

	// Close SDK
	defer fSetup.CloseSDK()

	err = fSetup.AdminSetup()
	if err != nil {
		t.Errorf("Failed to set up network admin: %v\n", err)
		return
	}

	// fSetup.ChannelID = userInput("Enter channel name: ")
	// t.Errorf("ChannelID is named %s\n", fSetup.ChannelID)

	err = fSetup.ChannelSetup()
	if err != nil {
		t.Errorf("Failed to set up channel: %v\n", err)
		return
	}

	err = fSetup.ChainCodeInstallationInstantiation()
	if err != nil {
		t.Errorf("Failed to install and instantiate chaincode: %v\n", err)
		return
	}

	// Set up client for chaincode invocations
	err = fSetup.ClientSetup()
	if err != nil {
		t.Errorf("Failed to set up client: %v\n", err)
		return
	}

	// Initialize vote entry
	payload, err := fSetup.InitVoteSDK("1", "1", "male", "30", "12345", "54321")
	if err != nil {
		t.Errorf("Failed to initialize entry: %v\n", err)
		return
	}

	t.Log("Successfully initialized test vote entry")

	// Retrieve vote entry
	payload, err = fSetup.GetVoteSDK("1", "1")
	if err != nil {
		t.Errorf("Failed to read entry: %v\n", err)
		return
	}

	t.Logf("Payload of test getVote: %s\n", payload)
}