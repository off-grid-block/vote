package main

import (
	"fmt"
	"encoding/json"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)


// ============================================================
// initVote - create a new vote and store into chaincode state
// ============================================================
func (vc *VoteChaincode) initVote(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	type voteTransientInput struct {
		PollID		string 	`json:"pollID"`
		VoterID		string 	`json:"voterID"`
		VoterSex 	string 	`json:"voterSex"`
		VoterAge	int 	`json:"voterAge"`
		Salt 		string 	`json:"salt"`
		VoteHash 	string 	`json:"voteHash"`
	}

	fmt.Println("- start init vote")

	if len(args) != 0 {
		return shim.Error("Private data should be passed in transient map.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	voteJsonBytes, success := transMap["vote"]
	if !success {
		return shim.Error("vote must be a key in the transient map")
	}

	if len(voteJsonBytes) == 0 {
		return shim.Error("vote value in transient map cannot be empty JSON string")
	}

	var voteInput voteTransientInput
	err = json.Unmarshal(voteJsonBytes, &voteInput)
	if err != nil {
		return shim.Error("failed to decode JSON of: " + string(voteJsonBytes))
	}

	// input sanitation

	if len(voteInput.PollID) == 0 {
		return shim.Error("poll ID field must be a non-empty string")
	} 

	if len(voteInput.VoterID) == 0 {
		return shim.Error("voter ID field must be a non-empty string")
	} 

	if voteInput.VoterAge <= 0 {
		return shim.Error("age field must be > 0")
	}

	if len(voteInput.VoterSex) == 0 {
		return shim.Error("sex field must be a non-empty string")
	} 

	if len(voteInput.Salt) == 0 {
		return shim.Error("salt must be > 0")
	}

	if len(voteInput.VoteHash) == 0 {
		return shim.Error("vote hash field must be a non-empty string")
	}

	var p poll

	existingPollAsBytes, err := stub.GetState(voteInput.PollID)
	if err != nil {
		return shim.Error("Failed to get associated poll: " + err.Error())
	} else if existingPollAsBytes == nil {
		return shim.Error("Poll does not exist: " + voteInput.PollID)
	}

	err = json.Unmarshal(existingPollAsBytes, &p)
	if err != nil {
		return shim.Error(err.Error())
	}

	// deserialize close date string into Time object
	closeDate, err := time.Parse(dateFormat, p.CloseDate)
	if err != nil {
		return shim.Error(err.Error())
	}

	// check if the poll has already closed
	if closeDate.Before(time.Now()) || p.Status == "closed" {
		return shim.Error("Poll has already closed.")
	}

	// Increment num votes of poll
	p.NumVotes++
	pollJSONasBytes, err := json.Marshal(p)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(voteInput.PollID, pollJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// create a composite key for vote collection using the poll ID and voter ID
	attrVoteCompositeKey := []string{voteInput.PollID, voteInput.VoterID}
	voteCompositeKey, err := stub.CreateCompositeKey("vote", attrVoteCompositeKey)
	if err != nil {
		return shim.Error("Failed to create composite key for vote: " + err.Error())
	}

	// check if value for voteCompositeKey already exists
	existingVoteAsBytes, err := stub.GetState(voteCompositeKey)
	if err != nil {
		return shim.Error("Failed to get vote: " + err.Error())
	} else if existingVoteAsBytes != nil {
		return shim.Error("This vote already exists: " + voteInput.PollID + voteInput.VoterID)
	}

	vote := &vote{
		ObjectType: "vote",
		PollID: voteInput.PollID,
		VoterID: voteInput.VoterID,
		VoterAge: voteInput.VoterAge,
		VoterSex: voteInput.VoterSex,
	}
	voteJSONasBytes, err := json.Marshal(vote)
	if err != nil {
		return shim.Error(err.Error())
	}

	// put state for voteCompositeKey
	err = stub.PutState(voteCompositeKey, voteJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// create a composite key for vote private details collections using the poll ID, voter ID, and salt
	attrVotePrivateDetailsCompositeKey := []string{voteInput.PollID, voteInput.VoterID, voteInput.Salt}
	votePrivateDetailsCompositeKey, err := stub.CreateCompositeKey("vote", attrVotePrivateDetailsCompositeKey)
	if err != nil {
		return shim.Error("Failed to create composite key for vote private details: " + err.Error())
	}

	votePrivateDetails := &votePrivateDetails {
		ObjectType: "votePrivateDetails",
		PollID: voteInput.PollID,
		VoterID: voteInput.VoterID,
		Salt: voteInput.Salt,
		VoteHash: voteInput.VoteHash,
	}
	votePrivateDetailsBytes, err := json.Marshal(votePrivateDetails)
	if err != nil {
		return shim.Error(err.Error())
	}

	// put state for votePrivateDetailsCompositeKey
	err = stub.PutPrivateData("collectionVotePrivateDetails", votePrivateDetailsCompositeKey, votePrivateDetailsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//register event
	err = stub.SetEvent("initEvent", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end init vote (success)")
	return shim.Success(nil)
}
