package main

import (
	"fmt"
	"strings"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)


type vote struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	VoterID		string 	`json:"voterID"`
	VoterSex 	string 	`json:"voterSex"`
	VoterAge	int 	`json:"voterAge"`
}

type votePrivateDetails struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	VoterID		string 	`json:"voterID"`
	Salt 		string 	`json:"salt"`
	VoteHash 	string 	`json:"voteHash"`
}


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

	existingVoteAsBytes, err := stub.GetPrivateData("collectionVote", voteInput.PollID + voteInput.VoterID)
	if err != nil {
		return shim.Error("Failed to get vote: " + err.Error())
	} else if existingVoteAsBytes != nil {
		fmt.Println("This vote already exists: " + voteInput.PollID + voteInput.VoterID)
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

	err = stub.PutPrivateData("collectionVote", voteInput.PollID + voteInput.VoterID, voteJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
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

	err = stub.PutPrivateData(
		"collectionVotePrivateDetails", 
		voteInput.PollID + voteInput.VoterID + voteInput.Salt, 
		votePrivateDetailsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//register event
	err = stub.SetEvent("initEvent", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end init vote (success)")
	return shim.Success([]byte(voteInput.Salt))
}

// =====================================================
// getVote - retrieve vote metadata from chaincode state
// =====================================================

func (vc *VoteChaincode) getVote(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting vote key to query")
	}

	voteKey := args[0]
	// ==== retrieve the vote ====
	voteAsBytes, err := stub.GetPrivateData("collectionVote", voteKey)
	if err != nil {
		return shim.Error("{\"Error\":\"Failed to get state for " + voteKey + "\"}")
	} else if voteAsBytes == nil {
		return shim.Error("{\"Error\":\"Vote does not exist: " + voteKey + "\"}")
	}

	return shim.Success(voteAsBytes)
}

// ===============================================================
// getVotePrivateDetails - retrieve vote hash from chaincode state
// ===============================================================

func (vc *VoteChaincode) getVotePrivateDetails(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting vote key to query")
	}

	voteKey := args[0]
	voteAsBytes, err := stub.GetPrivateData("collectionVotePrivateDetails", voteKey)
	if err != nil {
		return shim.Error("{\"Error\":\"Failed to get private details for " + voteKey + "\"}")
	} else if voteAsBytes == nil {
		return shim.Error("{\"Error\":\"Vote private details do not exist: " + voteKey + "\"}")
	}

	return shim.Success(voteAsBytes)
}

func (vc *VoteChaincode) getVotePrivateDetailsHash(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting vote key to query")
	}

	voteKey := args[0]
	voteHashAsBytes, err := stub.GetPrivateDataHash("collectionVotePrivateDetails", voteKey)
	if err != nil {
		return shim.Error("{\"Error\":\"Failed to get private data hash for " + voteKey + "\"}")
	} else if voteHashAsBytes == nil {
		return shim.Error("{\"Error\":\"Vote private data does not exist: " + voteKey + "\"}")
	}

	return shim.Success(voteHashAsBytes)
}

// =================================================
// amendVote - replace vote hash with new vote hash
// =================================================

func (vc *VoteChaincode) amendVote(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Pass data using transient map")		
	}

	fmt.Println("- begin amend vote")

	type amendVoteTransientInput struct {
		VoteKey string `json:"voteKey"`
		NewHash string `json:"newHash"`
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	amendVoteJsonBytes, ok := transMap["amend_vote"]
	if !ok {
		return shim.Error("amend_vote must be key in transient map")
	}

	if len(amendVoteJsonBytes) == 0 {
		return shim.Error("amend_vote value in transient map must be non-empty")
	}

	var amendVoteInput amendVoteTransientInput
	err = json.Unmarshal(amendVoteJsonBytes, &amendVoteInput)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(amendVoteJsonBytes))
	}

	if len(amendVoteInput.VoteKey) == 0 {
		return shim.Error("vote key field must be a non-empty string")
	}

	if len(amendVoteInput.NewHash) == 0 {
		return shim.Error("New hash field must be a non-empty string")
	}

	voteAsBytes, err := stub.GetPrivateData("collectionVotePrivateDetails", amendVoteInput.VoteKey)
	if err != nil {
		return shim.Error("Failed to get private vote data:" + err.Error())
	} else if voteAsBytes == nil {
		return shim.Error("Vote does not exist: " + amendVoteInput.VoteKey)
	}

	amendedVote := votePrivateDetails{}
	err = json.Unmarshal(voteAsBytes, &amendedVote)
	if err != nil {
		return shim.Error(err.Error())
	}
	amendedVote.VoteHash = amendVoteInput.NewHash

	voteJSONasBytes, _ := json.Marshal(amendedVote)
	err = stub.PutPrivateData(
		"collectionVotePrivateDetails", 
		amendedVote.PollID + amendedVote.VoterID + amendedVote.Salt,
		voteJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end amend vote (success)")
	return shim.Success(nil)
}


// ===== Parametrized rich queries =========================================================

// =========================================================================================
// queryVotesByPoll takes the poll ID as a parameter, builds a query string using
// the passed poll ID, executes the query, and returns the result set.
// =========================================================================================
func (vc *VoteChaincode) queryVotesByPoll(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: Poll ID")
	}

	pollID := strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"vote\",\"pollID\":\"%s\"}}", pollID)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

// =========================================================================================
// queryVotesByVoter takes the voter ID as a parameter, builds a query string using
// the passed voter ID, executes the query, and returns the result set.
// =========================================================================================	
func (vc *VoteChaincode) queryVotesByVoter(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: Voter ID")
	}

	voterID := strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"vote\",\"voterID\":\"%s\"}}", voterID)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

// ===== Ad hoc rich queries ===============================================================

// =========================================================================================
// Taken from fabric-samples/marbles_chaincode.go.
// queryVotes uses a query string to perform a query for votes.
// Query string matching state database syntax is passed in and executed as is.
// Supports ad hoc queries that can be defined at runtime by the client.
// =========================================================================================
func (vc *VoteChaincode) queryVotes(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := args[0]
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}
