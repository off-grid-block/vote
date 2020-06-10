package vote

import (
	// "fmt"
	"github.com/off-grid-block/vote/blockchain"
	"github.com/off-grid-block/vote/web"
	"github.com/pkg/errors"
	// "os"
	// "bufio"
	// "strings"

	ipfs "github.com/ipfs/go-ipfs-api"
)


func SetupApp() (*web.Application, error) {

	var app *web.Application

	fabricSDK, err := SetupSDK()
	if err != nil {
		return app, err
	}

	app.FabricSDK = &fabricSDK
	app.IpfsShell = ipfs.NewShell("localhost:5001")

	return app, nil
}


// Sets up the SDK for the voting application
func SetupSDK() (blockchain.SDKConfig, error) {

	fSetup := blockchain.SDKConfig {
		OrdererID: 			"orderer.example.com",
		ChannelID: 			"mychannel",
		ChannelConfig:		"/Users/brianli/deon/fabric-samples/first-network/channel-artifacts/channel.tx",
		ChainCodeID:		"vote",
		ChaincodeGoPath:	"/Users/brianli/deon",
		ChaincodePath:		"vote/chaincode",
		OrgAdmin:			"Admin",
		OrgName:			"org1",
		ConfigFile:			"config.yaml",
		UserName:			"User1",
	}

	err := fSetup.Initialization()
	if err != nil {
		return fSetup, errors.WithMessage(err, "Unable to initialize the Fabric SDK")
	}

	// Close SDK
	defer fSetup.CloseSDK()

	err = fSetup.AdminSetup()
	if err != nil {
		return fSetup, errors.WithMessage(err, "Failed to set up admin")
	}

	err = fSetup.ChannelSetup()
	if err != nil {
		// fmt.Printf("Failed to set up channel: %v\n", err)
		return fSetup, errors.WithMessage(err, "Failed to set up channel")
	}

	err = fSetup.ChainCodeInstallationInstantiation()
	if err != nil {
		// fmt.Printf("Failed to install and instantiate chaincode: %v\n", err)
		return fSetup, errors.WithMessage(err, "Failed to install and instantiate chaincode")
	}

	return fSetup, nil
}


// func main() {

// 	fSetup := blockchain.SDKConfig {
// 		OrdererID: 			"orderer.example.com",
// 		ChannelID: 			"mychannel",
// 		ChannelConfig:		"/Users/brianli/deon/fabric-samples/first-network/channel-artifacts/channel.tx",
// 		ChainCodeID:		"vote",
// 		ChaincodeGoPath:	"/Users/brianli/deon",
// 		ChaincodePath:		"vote/chaincode",
// 		OrgAdmin:			"Admin",
// 		OrgName:			"org1",
// 		ConfigFile:			"config.yaml",
// 		UserName:			"User1",
// 	}

// 	err := fSetup.Initialization()
// 	if err != nil {
// 		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
// 		return
// 	}

// 	// Close SDK
// 	defer fSetup.CloseSDK()

// 	err = fSetup.AdminSetup()
// 	if err != nil {
// 		fmt.Printf("Failed to set up network admin: %v\n", err)
// 		return
// 	}

// 	// fSetup.ChannelID = userInput("Enter channel name: ")
// 	// fmt.Printf("ChannelID is named %s\n", fSetup.ChannelID)

// 	err = fSetup.ChannelSetup()
// 	if err != nil {
// 		fmt.Printf("Failed to set up channel: %v\n", err)
// 		return
// 	}

// 	err = fSetup.ChainCodeInstallationInstantiation()
// 	if err != nil {
// 		fmt.Printf("Failed to install and instantiate chaincode: %v\n", err)
// 		return
// 	}

// 	// err = fSetup.ClientSetup()
// 	// if err != nil {
// 	// 	fmt.Printf("Failed to set up client: %v\n", err)
// 	// 	return
// 	// }

// 	// create shell to connect to IPFS
// 	sh := ipfs.NewShell("localhost:5001")

// 	app := &web.Application{
// 		FabricSDK: &fSetup,
// 		IpfsShell: sh,
// 	}

// 	web.Serve(app)
// }


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