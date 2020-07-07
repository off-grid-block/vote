package blockchain

import (
	"fmt"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
)

// FabricSetup implementation
type SetupSDK struct {
	ConfigFile      string
	OrgID           string
	OrdererID       string
	ChannelID       string
	initialized     bool
	ChannelConfig   string
	ChaincodeGoPath string
	ChaincodePath   map[string]string
	OrgAdmin        string
	OrgName         string
	UserName        string
	Client          *channel.Client
	Mgmt            *resmgmt.Client
	Fsdk            *fabsdk.FabricSDK
	Event           *event.Client
	MgmtIdentity	msp.SigningIdentity
}

// Initialization setups new sdk
func (s *SetupSDK) Initialization() error {

	// Add parameters for the initialization
	if s.initialized {
		return errors.New("sdk is already initialized")
	}

	// Initialize the SDK with the configuration file
	fsdk, err := fabsdk.New(config.FromFile(s.ConfigFile))
	if err != nil {
		return errors.WithMessage(err, "failed to create SDK")
	}
	s.Fsdk = fsdk
	fmt.Println("SDK is now created")

	fmt.Println("Initialization Successful")
	s.initialized = true

	return nil

}

func (s *SetupSDK) AdminSetup() error {

	// The resource management client is responsible for managing channels (create/update channel)
	resourceManagerClientContext := s.Fsdk.Context(fabsdk.WithUser(s.OrgAdmin), fabsdk.WithOrg(s.OrgName))
//	if err != nil {
//		return errors.WithMessage(err, "failed to load Admin identity")
//	}
	resMgmtClient, err := resmgmt.New(resourceManagerClientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create channel management client from Admin identity")
	}
	s.Mgmt = resMgmtClient
	fmt.Println("Resource management client created")

	// The MSP client allow us to retrieve user information from their identity, like its signing identity which we will need to save the channel
	mspClient, err := mspclient.New(s.Fsdk.Context(), mspclient.WithOrg(s.OrgName))
	if err != nil {
		return errors.WithMessage(err, "failed to create MSP client")
	}

	s.MgmtIdentity, err = mspClient.GetSigningIdentity(s.OrgAdmin)
	if err != nil {
		return errors.WithMessage(err, "failed to get mgmt signing identity")
	}

	return nil
}

func (s *SetupSDK) ChannelSetup() error {

	req := resmgmt.SaveChannelRequest{ChannelID: s.ChannelID, ChannelConfigPath: s.ChannelConfig, SigningIdentities: []msp.SigningIdentity{s.MgmtIdentity}}
	//create channel
	txID, err := s.Mgmt.SaveChannel(req, resmgmt.WithOrdererEndpoint(s.OrdererID))
	if err != nil || txID.TransactionID == "" {
		return errors.WithMessage(err, "failed to save channel")
	}
	fmt.Println("Channel created")

	// Make mgmt user join the previously created channel
	if err = s.Mgmt.JoinChannel(s.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(s.OrdererID)); err != nil {
		return errors.WithMessage(err, "failed to make mgmt join channel")
	}
	fmt.Println("Channel joined")

	return nil
}

// Create collection config to for chaincode instantiation
func newCollectionConfig(colName, policy string, reqPeerCount, maxPeerCount int32, blockToLive uint64) (*cb.CollectionConfig, error) {
	p, err := cauthdsl.FromString(policy)
	if err != nil {
        return nil, err
    }
    cpc := &cb.CollectionPolicyConfig{
        Payload: &cb.CollectionPolicyConfig_SignaturePolicy{
            SignaturePolicy: p,
        },
    }
    return &cb.CollectionConfig{
        Payload: &cb.CollectionConfig_StaticCollectionConfig{
            StaticCollectionConfig: &cb.StaticCollectionConfig{
                Name:              colName,
                MemberOrgsPolicy:  cpc,
                RequiredPeerCount: reqPeerCount,
                MaximumPeerCount:  maxPeerCount,
                BlockToLive:       blockToLive,
            },
        },
    }, nil
}

// Installs and instantiates chaincode
func (s *SetupSDK) ChainCodeInstallationInstantiation(ccID string) error {

	// Create the chaincode package that will be sent to the peers
	ccPackage, err := packager.NewCCPackage(s.ChaincodePath[ccID], s.ChaincodeGoPath)
	if err != nil {
		return errors.WithMessage(err, "failed to create chaincode package")
	}
	fmt.Println("Chaincode package created")

	// Install the chaincode to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: ccID, Path: s.ChaincodePath[ccID], Version: "0", Package: ccPackage}
	_, err = s.Mgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return errors.WithMessage(err, "failed to install chaincode")
	}
	fmt.Println("Chaincode installed")

	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"})

	// Create collection config for collectionVotePrivateDetails
	var collCfgPrivVoteRequiredPeerCount, collCfgPrivVoteMaximumPeerCount int32
	var collCfgPrivVoteBlockToLive uint64 

	collCfgPrivVoteName              := "collectionVotePrivateDetails"
	collCfgPrivVoteBlockToLive       = 3
	collCfgPrivVoteRequiredPeerCount = 0
	collCfgPrivVoteMaximumPeerCount  = 3
	collCfgPrivVotePolicy            := "OR('Org1MSP.member')"

	collCfgPrivVote, err := newCollectionConfig(
		collCfgPrivVoteName,
		collCfgPrivVotePolicy,
		collCfgPrivVoteRequiredPeerCount,
		collCfgPrivVoteMaximumPeerCount,
		collCfgPrivVoteBlockToLive)
	if err != nil {
	    return errors.WithMessage(err, "failed to create collection config")
	}

	// Create collection config for collectionPollPrivateDetails
	var collCfgPrivPollRequiredPeerCount, collCfgPrivPollMaximumPeerCount int32
	var collCfgPrivPollBlockToLive uint64 

	collCfgPrivPollName              := "collectionPollPrivateDetails"
	collCfgPrivPollBlockToLive       = 3
	collCfgPrivPollRequiredPeerCount = 0
	collCfgPrivPollMaximumPeerCount  = 3
	collCfgPrivPollPolicy            := "OR('Org1MSP.member')"

	collCfgPrivPoll, err := newCollectionConfig(
		collCfgPrivPollName,
		collCfgPrivPollPolicy,
		collCfgPrivPollRequiredPeerCount,
		collCfgPrivPollMaximumPeerCount,
		collCfgPrivPollBlockToLive)
	if err != nil {
	    return errors.WithMessage(err, "failed to create collection config")
	}

	cfg := []*cb.CollectionConfig{collCfgPrivVote, collCfgPrivPoll}

	// instantiate chaincode with cc policy and collection configs
	resp, err := s.Mgmt.InstantiateCC(
		// Channel ID
		s.ChannelID, 
		// InstantiateCCRequest struct
		resmgmt.InstantiateCCRequest{
			Name: ccID, 
			Path: s.ChaincodeGoPath, 
			Version: "0", 
			Args: [][]byte{[]byte("init")}, 
			Policy: ccPolicy, 
			CollConfig: cfg,
		},
		// options
		resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil || resp.TransactionID == "" {
		return errors.WithMessage(err, "failed to instantiate the chaincode")
	}
	fmt.Println("Chaincode instantiated")
	return nil
}

//setup client and setupt access to channel events
func (s*SetupSDK)  ClientSetup() error {
	// Channel client is used to Query or Execute transactions
	var err error
	clientChannelContext := s.Fsdk.ChannelContext(s.ChannelID, fabsdk.WithUser(s.UserName))
	s.Client, err = channel.New(clientChannelContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new channel client")
	}
	fmt.Println("Channel client created")

	// Creation of the client which will enables access to our channel events
	s.Event, err = event.New(clientChannelContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new event client")
	}
	fmt.Println("Event client created")

	return nil
}

func (s *SetupSDK) CloseSDK() {
	s.Fsdk.Close()
}