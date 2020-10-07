package blockchain

import (
	"fmt"
	cb "github.com/off-grid-block/fabric-protos-go/common"
	"github.com/off-grid-block/fabric-sdk-go/pkg/client/channel"
	"github.com/off-grid-block/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/off-grid-block/fabric-sdk-go/pkg/client/msp"
	"github.com/off-grid-block/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/off-grid-block/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/off-grid-block/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/off-grid-block/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/off-grid-block/fabric-sdk-go/pkg/common/providers/context"
	"github.com/off-grid-block/fabric-sdk-go/pkg/core/config"
	packager "github.com/off-grid-block/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/off-grid-block/fabric-sdk-go/pkg/fabsdk"
	"github.com/off-grid-block/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
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

func (s *SetupSDK) CreateChannelClient(clientContext context.ChannelProvider) (*channel.Client, error) {

    client, err := channel.New(clientContext)
    if err != nil {
        return nil, fmt.Errorf("failed to create new channel client: %v", err)
    }

    return client, nil
}


func (s *SetupSDK) CreateEventClient(clientContext context.ChannelProvider, eventID string) (*event.Client, fab.Registration, <-chan *fab.CCEvent, error) {

    event, err := event.New(clientContext)
    if err != nil {
        return nil, nil, nil, fmt.Errorf("failed to create new client context: %v", err)
    }
    fmt.Println("Event client created")

    // register chaincode event
    registered, notifier, err := event.RegisterChaincodeEvent("vote", eventID)
    if err != nil {
        return nil, nil, nil, errors.WithMessage(err, "failed to register cc event")
    }

    return event, registered, notifier, nil
}

func (s *SetupSDK) CloseSDK() {
	s.Fsdk.Close()
}
