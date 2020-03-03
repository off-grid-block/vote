package main

import (
	"fmt"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type poll struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
}

type pollPrivateDetails struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	Salt 		string 	`json:"salt"`
	PollHash 	string 	`json:"pollHash"`
}

// ===========================================================
// initPoll - create a new poll and store into chaincode state
// ===========================================================

func (vc *VoteChaincode) initPoll(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	type pollTransientInput struct {
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

	var pollInput pollTransientInput
	err = json.Unmarshal(pollJsonBytes, &pollInput)
	if err != nil {
		return shim.Error("failed to decode JSON of: " + string(pollJsonBytes))
	}

	if len(pollInput.PollID) == 0 {
		return shim.Error("poll ID field must be a non-empty string")
	}

	if len(pollInput.Salt) == 0 {
		return shim.Error("salt field must be a non-empty string")
	} 

	if len(pollInput.PollHash) == 0 {
		return shim.Error("poll hash field must be a non-empty string")
	}

	existingPollAsBytes, err := stub.GetPrivateData("collectionPoll", pollInput.PollID)
	if err != nil {
		return shim.Error("Failed to get vote: " + err.Error())
	} else if existingPollAsBytes != nil {
		fmt.Println("This poll already exists: " + pollInput.PollID)
		return shim.Error("This poll already exists: " + pollInput.PollID)
	}

	poll := &poll{
		ObjectType: "poll",
		PollID: pollInput.PollID,
	}
	pollJSONasBytes, err := json.Marshal(poll)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutPrivateData("collectionPoll", pollInput.PollID, pollJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	pollPrivateDetails := &pollPrivateDetails {
		ObjectType: "pollPrivateDetails",
		PollID: pollInput.PollID,
		Salt: pollInput.Salt,
		PollHash: pollInput.PollHash,
	}
	pollPrivateDetailsBytes, err := json.Marshal(pollPrivateDetails)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutPrivateData(
		"collectionPollPrivateDetails", 
		pollInput.PollID + pollInput.Salt, 
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
	return shim.Success([]byte(pollInput.Salt))

}

// ===============================================================
// getPollPrivateDetails - retrieve poll hash from chaincode state
// ===============================================================

func (vc *VoteChaincode) getPollPrivateDetails(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting cc key to query")
	}

	ccKey := args[0]
	pollAsBytes, err := stub.GetPrivateData("collectionPollPrivateDetails", ccKey)
	if err != nil {
		return shim.Error("{\"Error\":\"Failed to get private details for " + ccKey + "\"}")
	} else if pollAsBytes == nil {
		return shim.Error("{\"Error\":\"Poll private details do not exist: " + ccKey + "\"}")
	}

	return shim.Success(pollAsBytes)
}