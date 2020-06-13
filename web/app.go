package web

import (
	"github.com/off-grid-block/core-interface/pkg/sdk"
	// "net/http"
	// "log"
	// "fmt"
	// "time"
	// "github.com/gorilla/mux"
	ipfs "github.com/ipfs/go-ipfs-api"
)

// Struct containing Fabric SDK setup data. Objects of type
// Application have access to the SDK's chaincode interface.
type Application struct {
	FabricSDK *sdk.SDKConfig
	IpfsShell *ipfs.Shell
} 

// PRELIMINARY DEF: struct to hold vote data.
type initVoteRequestBodyAPI struct {
	PollID 			string 		`json:"pollID"`
	VoterID 		string 		`json:"voterID"`
	Sex 			string 		`json:"sex"`
	Age 			string 		`json:"age"`
	Content 		interface{} `json:"content"`
}

type initPollRequestBodyAPI struct {
	PollID 			string 		`json:"pollID"`
	Title 			string 		`json:"title"`
	Content 		interface{} `json:"content"`
}

type updatePollStatusRequestBodyAPI struct {
	Status 			string 		`json:"status"`
}

type voteDetailsHttpResponse struct {
	PollID 			string 		`json:"pollID"`
	VoterID 		string 		`json:"voterID"`
	VoterSex 		string 		`json:"voterSex"`
	VoterAge 		int 		`json:"voterAge"`
	Content 		interface{} `json:"content"`
	PrivateHash 	string 		`json:"privateHash"`
}

type pollDetailsHttpResponse struct {
	PollID 			string
	Title 			string
	Status 			string
	NumVotes 		int
	Content 		interface{}
}

type voteResponseSDK struct {
	ObjectType 		string 		`json:"docType"`
	PollID			string 		`json:"pollID"`
	VoterID			string 		`json:"voterID"`
	VoterSex 		string 		`json:"voterSex"`
	VoterAge		int 		`json:"voterAge"`
	PrivateHash 	string 		`json:"privateHash"`
}

type votePrivateDetailsResponseSDK struct {
	ObjectType 		string 		`json:"docType"`
	PollID			string 		`json:"pollID"`
	VoterID			string 		`json:"voterID"`
	Salt 			string 		`json:"salt"`
	VoteHash 		string 		`json:"voteHash"`
}

type pollResponseSDK struct {
	ObjectType 		string
	PollID 			string
	Title 			string
	Status  		string
	NumVotes 		int
}

type pollPrivateDetailsResponseSDK struct {
	ObjectType 		string 		`json:"docType"`
	PollID			string 		`json:"pollID"`
	Salt 			string 		`json:"salt"`
	PollHash 		string 		`json:"pollHash"`
}
