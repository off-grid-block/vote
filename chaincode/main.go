package main

import (
	"fmt"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"errors"
		
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


type poll struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	Title 		string 	`json:"title"`
	Status		string 	`json:"status"`
	NumVotes	int 	`json:"numVotes"`
}


type pollPrivateDetails struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	Salt 		string 	`json:"salt"`
	PollHash 	string 	`json:"pollHash"`
}


type VoteChaincode struct {
}


func main() {
	err := shim.Start(new(VoteChaincode))
	if err != nil {
		fmt.Printf("Error starting Vote chaincode: %s", err)
	}
}

// ============================
// Init - initializes chaincode
// ============================
func (vc *VoteChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// =================================================
// Invoke - starting point for chaincode invocations
// =================================================
func (vc *VoteChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	indycreatorbytes, _ := stub.GetCreator()
	type IndyCreator struct {
		Did string
	}
	indycreator := &IndyCreator{}
	if err := json.Unmarshal(indycreatorbytes, &indycreator); err != nil {
		panic(err)
	}
	fmt.Println(indycreator)
	status, err := VerifyVoterProof(indycreator.Did)
	if status == false || err != nil {
		return shim.Error("Proof verification failed : " + err.Error())
	}
	fmt.Println("proof verification success")

	fn, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + fn)

	switch fn {
	case "initVote":
		return vc.initVote(stub, args)
	case "getVote":
		return vc.getVote(stub, args)
	case "getVotePrivateDetails":
		return vc.getVotePrivateDetails(stub, args)
	case "getVotePrivateDetailsHash":
		return vc.getVotePrivateDetailsHash(stub, args)
	case "amendVote":
		return vc.updateVote(stub, args)
	case "queryVotePrivateDetailsByPoll":
		return vc.queryVotePrivateDetailsByPoll(stub, args)
	case "queryVotesByPoll":								// parametrized rich query w/ poll ID
		return vc.queryVotesByPoll(stub, args)
	case "queryVotesByVoter":								// parametrized rich query w/ voter ID
		return vc.queryVotesByVoter(stub, args)			
	case "queryVotes":										// ad hoc rich query
		return vc.queryVotes(stub, args)
	case "queryAllPolls":
		return vc.queryAllPolls(stub, args)
	case "initPoll":
		return vc.initPoll(stub, args)
	case "getPoll":
		return vc.getPoll(stub, args)
	case "getPollPrivateDetails":
		return vc.getPollPrivateDetails(stub, args)
	case "updatePollStatus":
		return vc.updatePollStatus(stub, args)
	}

	fmt.Println("invoke did not find fn: " + fn)
	return shim.Error("Received unknown function invocation")
}

// ===========================================================================================
// Taken from fabric-samples/marbles_chaincode.go.
// constructQueryResponseFromIterator constructs a JSON array containing query results from
// a given result iterator
// ===========================================================================================
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

// =========================================================================================
// Taken from fabric-samples/marbles_chaincode.go.
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}


//VerifyVoterProof function receives did as argument and verifies that the identity has required attributes.
func VerifyVoterProof(Did string) (bool, error) {

	//Check validity of arguments
	DidSize := len(Did)
	if DidSize == 0 {
		return false, errors.New("Empty DID received")
	}
	if (DidSize) != 22 {
		return false, errors.New("Size of DID is not 22")
	}

	//Initialize
	type Attributes struct {
		AppName string `json:"app_name"`
		AppID   string `json:"app_id"`
	}
	type IndyResponse struct {
		Status     string     `json:"status"`
		Attributes Attributes `json:"attributes"`
	}
	ProofAttributes := "app_name,app_id"
	RequiredAppName := "voter"
	RequiredAppID := "101"

	//Prepare Payload for Indy
	IndyURL := "http://localhost:7997/verify_proof"
	Payload := []byte("{\"proof_attr\" : \"" + ProofAttributes + "\",\"their_did\" : \"" + Did + "\"}")
	Request, _ := http.NewRequest("POST", IndyURL, bytes.NewBuffer(Payload))
	Request.Header.Add("content-type", "text/plain")
	Response, err := http.DefaultClient.Do(Request)
	if err != nil || Response == nil || Response.StatusCode != 200 {
		return false, errors.New("!!!!!!!!!!! Error connecting to Indy Server ")
	}
	defer Response.Body.Close()

	//Validate Response from Indy
	Body, _ := ioutil.ReadAll(Response.Body)
	ResponseJSON := IndyResponse{}
	err = json.Unmarshal(Body, &ResponseJSON)
	if err != nil {
		return false, errors.New("Error unmarshaling Indy response")
	}
	if ResponseJSON.Status != "true" {
		return false, errors.New("Proof verification failed: attributes missing")
	}
	if !(ResponseJSON.Attributes.AppName == RequiredAppName && ResponseJSON.Attributes.AppID == RequiredAppID) {
		return false, errors.New("Attribute values didn't match")
	}
	return true, nil
}
