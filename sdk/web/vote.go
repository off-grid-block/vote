package web

import (
	"net/http"
	"fmt"
	"encoding/json"
	// "unicode/utf8"
)

// PRELIMINARY DEF: struct to hold vote data.
type VoteContent struct {
	YesOrNo bool
}

type FabricResponseVote struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	VoterID		string 	`json:"voterID"`
	VoterSex 	string 	`json:"voterSex"`
	VoterAge	int 	`json:"voterAge"`
	PrivateHash string 	`json:"privateHash"`
}

type FabricResponseVotePrivateDetails struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	VoterID		string 	`json:"voterID"`
	Salt 		string 	`json:"salt"`
	VoteHash 	string 	`json:"voteHash"`
}

// Initialize & push votes on the Fabric network
func (app *Application) initVoteHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	pollID := query.Get("pollid")
	voterID := query.Get("voterid")
	sex := query.Get("sex")
	age := query.Get("age")

	var v VoteContent

	// Decode HTTP request body and marshal into Vote struct.
	// If the bytes in the request body do not match the fields
	// of the Vote struct, the operation will fail.
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Push vote data to IPFS
	cid, err := app.IpfsAddVote(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate a random salt to concatenate with the vote's IPFS CID
	salt := GenerateSalt()

	// Print out values of arguments to initVote
	fmt.Println("poll ID:   ", pollID)
	fmt.Println("voter ID:  ", voterID)
	fmt.Println("voter sex: ", sex)
	fmt.Println("voter age: ", age)
	fmt.Println("salt:      ", salt)
	fmt.Println("CID:       ", cid)

	// Call InitVoteSDK() to initialize a vote on the Fabric network
	resp, err := app.FabricSDK.InitVoteSDK(pollID, voterID, sex, age, salt, cid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

// Retrieve vote from the Fabric network
func (app *Application) getVoteHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	pollID := query.Get("pollid")
	voterID := query.Get("voterid")
	salt := query.Get("salt")

	resp, err := app.FabricSDK.GetVoteSDK(pollID, voterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If a salt is not provided, send vote metadata response w/o private data hash
	if salt == "" {
		w.Write([]byte(resp))
		return
	}

	// If a salt is provided, marshal Fabric resp into FabricResponseVote so we can
	// add the associated private data hash to the http response
	var fabResp FabricResponseVote
	err = json.Unmarshal([]byte(resp), &fabResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// retrieve private data hash from Fabric ledger
	votePrivateDetailsHash, err := app.FabricSDK.GetVotePrivateDetailsHashSDK(pollID, voterID, salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add private data hash to response
	fabResp.PrivateHash = votePrivateDetailsHash

	httpResp, err := json.Marshal(fabResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(httpResp))
}

// Retrieve private details of a vote from the Fabric network
func (app *Application) getVotePrivateDetailsHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	pollID := query.Get("pollid")
	voterID := query.Get("voterid")
	salt := query.Get("salt")

	resp, err := app.FabricSDK.GetVotePrivateDetailsSDK(pollID, voterID, salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var fabResp FabricResponseVotePrivateDetails

	err = json.Unmarshal([]byte(resp), &fabResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	voteContentResp, err := app.IpfsGetData(fabResp.VoteHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(voteContentResp))
}