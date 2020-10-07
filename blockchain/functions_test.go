
package blockchain

import (
	"testing"
	// "fmt"
)

func TestChaincodeSDK(t *testing.T) {

	/* SDK Setup */

	fSetup := SetupSDK {
		OrdererID: 			"orderer.example.com",

		ChannelID: 			"mychannel",
		ChannelConfig: 		"/Users/brianli/deon/fabric-samples/first-network/channel-artifacts/channel.tx",

		ChainCodeID: 		"vote",
		ChaincodeGoPath: 	"/Users/brianli/deon",
		ChaincodePath:   	"vote/cc",
		OrgAdmin:        	"Admin",
		OrgName:         	"org1",
		ConfigFile:      	"/Users/brianli/deon/src/vote/config.yaml",

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
	_, err = fSetup.InitPollSDK("1", "t589t9fh98rfh23")
	if err != nil {
		t.Errorf("Failed to initialize poll: %v\n", err)
	}

	// Example vote with poll ID 1, voter ID 1
	_, err = fSetup.InitVoteSDK("1", "1", "m", "21", "36afgyr545tfgfd")
	if err != nil {
		t.Errorf("Failed to initialize vote: %v\n", err)
	}

	// Example vote with poll ID 1, voter ID 2
	_, err = fSetup.InitVoteSDK("1", "2", "f", "42", "f149j9vjf0cfqj5")
	if err != nil {
		t.Errorf("Failed to initialize vote: %v\n", err)
	}

	/* Run GetVote test */

	t.Run("GetVote", func(t *testing.T) {
		payload, err := fSetup.GetVoteSDK("1", "1")
		if err != nil {
			t.Errorf("Failed to get vote: %v\n", err)
		}
		t.Logf(payload)
	})

	/* Run GetVotePrivateDetails test */

	t.Run("GetVotePrivateDetails", func(t *testing.T) {
		payload, err := fSetup.GetVotePrivateDetailsSDK("1", "1")
		if err != nil {
			t.Errorf("Failed to get vote private details: %v\n", err)
		}
		t.Logf(payload)
	})

	/* Run GetVotePrivateDetailsHash test */

	t.Run("GetVotePrivateDetailsHash", func(t *testing.T) {
		payload, err := fSetup.GetVotePrivateDetailsHashSDK("1", "1")
		if err != nil {
			t.Errorf("Failed to get vote private details hash: %v\n", err)
		}
		t.Logf(payload)
	})

	/* Run QueryVotePrivateDetailsByPoll test */

	t.Run("QueryVotePrivateDetailsByPoll", func(t *testing.T) {

		queryResult, err := fSetup.QueryVotePrivateDetailsByPollSDK("1")
		if err != nil {
			t.Errorf("Failed to query vote private details by poll: %v\n", err)
		}

		for _, hash := range queryResult {
			t.Logf(hash)
		}
	})

	/* Run QueryVotesByPoll test */

	t.Run("QueryVotesByPoll", func(t *testing.T) {
		payload, err := fSetup.QueryVotesByPollSDK("1")
		if err != nil {
			t.Errorf("Failed to query vote private details by poll: %v\n", err)
		}
		t.Logf(payload)
	})

	/* Run QueryVotesByVoter test */

	t.Run("QueryVotesByVoter", func(t *testing.T) {
		payload, err := fSetup.QueryVotesByVoterSDK("1")
		if err != nil {
			t.Errorf("Failed to query vote private details by poll: %v\n", err)
		}
		t.Logf(payload)
	})

	/* Run GetPoll test */

	t.Run("GetPoll", func(t *testing.T) {
		payload, err := fSetup.GetPollSDK("1")
		if err != nil {
			t.Errorf("Failed to query poll: %v\n", err)
		}
		t.Logf(payload)
	})

	/* Run GetPollPrivateDetails test */

	t.Run("GetPollPrivateDetails", func(t *testing.T) {
		payload, err := fSetup.GetPollPrivateDetailsSDK("1")
		if err != nil {
			t.Errorf("Failed to query poll private details: %v\n", err)
		}
		t.Logf(payload)
	})

	/* Run UpdatePollStatus test */

	t.Run("UpdatePollStatus", func(t *testing.T) {
		_, err := fSetup.UpdatePollStatusSDK("1", "closed")
		if err != nil {
			t.Errorf("Failed to update poll status: %v\n", err)
		}
	})
}
