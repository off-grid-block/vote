
package blockchain

import (
	"testing"
	"github.com/off-grid-block/core-interface/pkg/sdk"
	"github.com/off-grid-block/vote/blockchain"
	// "fmt"
)

func TestChaincodeSDK(t *testing.T) {

	ccPath := map[string]string{"vote": "vote/chaincode"}

	/* SDK Setup */

	fSetup := &sdk.SDKConfig {
		OrdererID: 			"orderer.example.com",

		ChannelID: 			"mychannel",
		ChannelConfig: 		"/Users/brianli/deon/fabric-samples/first-network/channel-artifacts/channel.tx",

		// ChainCodeID: 		"vote",
		ChaincodeGoPath: 	"/Users/brianli/deon",
		ChaincodePath:   	ccPath,
		OrgAdmin:        	"Admin",
		OrgName:         	"org1",
		ConfigFile:      	"/Users/brianli/deon/src/vote/blockchain/config_test.yaml",

		UserName: 			"User1",
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

	err = fSetup.ClientSetup()
	if err != nil {
		t.Errorf("Failed to set up client: %v\n", err)
		return
	}

	/* Initialize Example Polls & Votes for Testing */

	// Example poll ID 1
	_, err = blockchain.InitPollSDK(fSetup, "1", "title", "t589t9fh98rfh23")
	if err != nil {
		t.Errorf("Failed to initialize poll: %v\n", err)
	}

	// Example vote with poll ID 1, voter ID 1
	_, err = blockchain.InitVoteSDK(fSetup, "1", "1", "m", "21", "36afgyr545tfgfd")
	if err != nil {
		t.Errorf("Failed to initialize vote: %v\n", err)
	}

	// Example vote with poll ID 1, voter ID 2
	_, err = blockchain.InitVoteSDK(fSetup, "1", "2", "f", "42", "f149j9vjf0cfqj5")
	if err != nil {
		t.Errorf("Failed to initialize vote: %v\n", err)
	}

	/* Run GetVote test */

	t.Run("GetVote", func(t *testing.T) {
		payload, err := blockchain.GetVoteSDK(fSetup, "1", "1")
		if err != nil {
			t.Errorf("Failed to get vote: %v\n", err)
		}
		t.Logf(payload)
	})

	/* Run GetVotePrivateDetails test */

	t.Run("GetVotePrivateDetails", func(t *testing.T) {
		payload, err := blockchain.GetVotePrivateDetailsSDK(fSetup, "1", "1")
		if err != nil {
			t.Errorf("Failed to get vote private details: %v\n", err)
		}
		t.Logf(payload)
	})

	/* Run GetVotePrivateDetailsHash test */

	t.Run("GetVotePrivateDetailsHash", func(t *testing.T) {
		payload, err := blockchain.GetVotePrivateDetailsHashSDK(fSetup, "1", "1")
		if err != nil {
			t.Errorf("Failed to get vote private details hash: %v\n", err)
		}
		t.Logf(payload)
	})

	/* Run QueryVotePrivateDetailsByPoll test */

	t.Run("QueryVotePrivateDetailsByPoll", func(t *testing.T) {

		queryResult, err := blockchain.QueryVotePrivateDetailsByPollSDK(fSetup, "1")
		if err != nil {
			t.Errorf("Failed to query vote private details by poll: %v\n", err)
		}

		for _, hash := range queryResult {
			t.Logf(hash)
		}
	})

	/* Run QueryVotesByPoll test */

	t.Run("QueryVotesByPoll", func(t *testing.T) {
		payload, err := blockchain.QueryVotesByPollSDK(fSetup, "1")
		if err != nil {
			t.Errorf("Failed to query vote private details by poll: %v\n", err)
		}
		t.Logf(payload)
	})

	/* Run QueryVotesByVoter test */

	t.Run("QueryVotesByVoter", func(t *testing.T) {
		payload, err := blockchain.QueryVotesByVoterSDK(fSetup, "1")
		if err != nil {
			t.Errorf("Failed to query vote private details by poll: %v\n", err)
		}
		t.Logf(payload)
	})

	/* Run GetPoll test */

	t.Run("GetPoll", func(t *testing.T) {
		payload, err := blockchain.GetPollSDK(fSetup, "1")
		if err != nil {
			t.Errorf("Failed to query poll: %v\n", err)
		}
		t.Logf(payload)
	})

	/* Run GetPollPrivateDetails test */

	t.Run("GetPollPrivateDetails", func(t *testing.T) {
		payload, err := blockchain.GetPollPrivateDetailsSDK(fSetup, "1")
		if err != nil {
			t.Errorf("Failed to query poll private details: %v\n", err)
		}
		t.Logf(payload)
	})

	/* Run UpdatePollStatus test */

	t.Run("UpdatePollStatus", func(t *testing.T) {
		_, err := blockchain.UpdatePollStatusSDK(fSetup, "1", "closed")
		if err != nil {
			t.Errorf("Failed to update poll status: %v\n", err)
		}
	})
}