package main

import (
	"fmt"
	"github.com/off-grid-block/vote/sdk/blockchain"
	"os"
	"bufio"
	"strings"
)

//function to read string
func userInput(request string) (inputval string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(request)
	inputval, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}
	returnval := strings.TrimRight (inputval, "\n")
	return returnval
}

func main() {

	fSetup := blockchain.SetupSDK {
		OrdererID: "orderer.example.com",

		ChannelID: 			"mychannel",
		ChannelConfig: 		"/Users/brianli/deon/fabric-samples/first-network/channel-artifacts/channel.tx",

		ChainCodeID: 		"vote",
		ChaincodeGoPath: 	"/Users/brianli/deon",
		ChaincodePath:   	"vote/cc",
		OrgAdmin:        	"Admin",
		OrgName:         	"org1",
		ConfigFile:      	"config.yaml",

		UserName: "User1",
	}

	err := fSetup.Initialization()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}

	// Close SDK
	defer fSetup.CloseSDK()

	err = fSetup.AdminSetup()
	if err != nil {
		fmt.Printf("Failed to set up network admin: %v\n", err)
		return
	}

	// fSetup.ChannelID = userInput("Enter channel name: ")
	// fmt.Printf("ChannelID is named %s\n", fSetup.ChannelID)

	err = fSetup.ChannelSetup()
	if err != nil {
		fmt.Printf("Failed to set up channel: %v\n", err)
		return
	}

	err = fSetup.ChainCodeInstallationInstantiation()
	if err != nil {
		fmt.Printf("Failed to install and instantiate chaincode: %v\n", err)
		return
	}

	err = fSetup.ClientSetup()
	if err != nil {
		fmt.Printf("Failed to set up client: %v\n", err)
		return
	}

	payload, err := blockchain.InitEntrySDK("1", "1", "male", 30, "12345", "54321")
	if err != nil {
		fmt.Printf("Failed to initialize entry: %v\n", err)
		return
	}

	// text := "{\"PollID\":\"1\",\"VoterID\":\"1\",\"VoterSex\":\"male\",\"VoterAge\":30,\"Salt\":\"12345\",\"VoteHash\":\"54321\"}"

}