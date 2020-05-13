package main

import (
	"fmt"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)


// ===========================================================
// initPoll - create a new poll and store into chaincode state
// ===========================================================

func (vc *VoteChaincode) initPoll(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	type pollTransientInput struct {
		Title 			string 	`json:"title"`
		PollID			string 	`json:"pollID"`
		Salt 			string 	`json:"salt"`
		PollHash 		string 	`json:"pollHash"`
	}

	fmt.Println("- start init vote")

	if len(args) != 0 {
		return shim.Error("Private data should be passed in transient map.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	pollJsonBytes, success := transMap["poll"]
	if !success {
		return shim.Error("poll must be a key in the transient map")
	}

	if len(pollJsonBytes) == 0 {
		return shim.Error("poll value in transient map cannot be empty JSON string")
	}

	var transInput pollTransientInput
	err = json.Unmarshal(pollJsonBytes, &transInput)
	if err != nil {
		return shim.Error("failed to decode JSON of: " + string(pollJsonBytes))
	}

	// Transient input validation

	if len(transInput.PollID) == 0 {
		return shim.Error("poll ID field must be a non-empty string")
	}

	if len(transInput.Title) == 0 {
		return shim.Error("title field must be a non-empty string")
	}

	if len(transInput.Salt) == 0 {
		return shim.Error("salt field must be a non-empty string")
	} 

	if len(transInput.PollHash) == 0 {
		return shim.Error("poll hash field must be a non-empty string")
	}

	// Put public poll information on chain
	
	existingPollAsBytes, err := stub.GetState(transInput.PollID)
	if err != nil {
		return shim.Error("Failed to get poll: " + err.Error())
	} else if existingPollAsBytes != nil {
		return shim.Error("This poll already exists: " + transInput.PollID)
	}

	poll := &poll{
		ObjectType: "poll",
		PollID: transInput.PollID,
		Title: transInput.Title,
		Status: "ongoing",
		NumVotes: 0,
	}

	pollJSONasBytes, err := json.Marshal(poll)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(transInput.PollID, pollJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// create a composite key for poll private details collections using the poll ID and salt
	pollPrivateDetailsCompositeKey, err := stub.CreateCompositeKey("poll", []string{transInput.PollID, transInput.Salt})
	if err != nil {
		return shim.Error("Failed to create composite key for poll private details: " + err.Error())
	}

	pollPrivateDetails := &pollPrivateDetails {
		ObjectType: "pollPrivateDetails",
		PollID: transInput.PollID,
		Salt: transInput.Salt,
		PollHash: transInput.PollHash,
	}

	pollPrivateDetailsBytes, err := json.Marshal(pollPrivateDetails)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutPrivateData(
		"collectionPollPrivateDetails", 
		pollPrivateDetailsCompositeKey, 
		pollPrivateDetailsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//register event
	err = stub.SetEvent("initEvent", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end init poll (success)")
	return shim.Success(nil)

}
