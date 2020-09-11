package main

import (
	"fmt"
	"github.com/off-grid-block/vote/blockchain"
	"github.com/off-grid-block/vote/web"
	"os"
	// "bufio"
	// "strings"

	ipfs "github.com/ipfs/go-ipfs-api"
)

func main() {

	ccPaths := map[string]string{
		"vote": "chaincode",
	}

	fSetup := blockchain.SetupSDK {
		OrdererID: 			"orderer.example.com",
		ChannelID: 			"mychannel",
		ChannelConfig:		os.Getenv("CHANNEL_CONFIG"),
		ChaincodeGoPath:	os.Getenv("CHAINCODE_GOPATH"),
		ChaincodePath:		ccPaths,
		OrgAdmin:			"Admin",
		OrgName:			"org1",
		ConfigFile:			"/src/config.yaml",
		UserName:			"User1",
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

	err = fSetup.ChainCodeInstallationInstantiation("vote")
	if err != nil {
		fmt.Printf("Failed to install and instantiate chaincode: %v\n", err)
		return
	}

	err = fSetup.ClientSetup()
	if err != nil {
		fmt.Printf("Failed to set up client: %v\n", err)
		return
	}

	// create shell to connect to IPFS
	sh := ipfs.NewShell(os.Getenv("IPFS_ENDPOINT"))

	app := &web.Application{
		FabricSDK: &fSetup,
		IpfsShell: sh,
	}

	web.Serve(app)
}


// //function to read string
// func userInput(request string) (inputval string) {
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Println(request)
// 	inputval, err := reader.ReadString('\n')
// 	if err != nil {
// 		fmt.Printf("Error %v\n", err)
// 	}
// 	returnval := strings.TrimRight (inputval, "\n")
// 	return returnval
// }