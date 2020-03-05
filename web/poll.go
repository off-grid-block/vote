package web

import (
	"fmt"
	"net/http"
	"encoding/json"	
	"reflect"
	
	"github.com/gorilla/mux"
)

// PRELIMINARY DEF: struct to hold vote data.
type VoteContent struct {
	YesOrNo bool
}

type PollContent struct {
	FirstChoice 	bool
	SecondChoice 	bool
	ThirdChoice		bool
}

type VoteFabricResponse struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	VoterID		string 	`json:"voterID"`
	VoterSex 	string 	`json:"voterSex"`
	VoterAge	int 	`json:"voterAge"`
	PrivateHash string 	`json:"privateHash"`
}

type VoteFabricResponsePrivateDetails struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	VoterID		string 	`json:"voterID"`
	Salt 		string 	`json:"salt"`
	VoteHash 	string 	`json:"voteHash"`
}

type PollFabricResponsePrivateDetails struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	Salt 		string 	`json:"salt"`
	PollHash 	string 	`json:"pollHash"`
}

// Retrieve vote from the Fabric network
func (app *Application) getVoteHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pollID := vars["pollid"]
	voterID := vars["voterid"]

	resp, err := app.FabricSDK.GetVoteSDK(pollID, voterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If a salt is provided, marshal Fabric resp into VoteFabricResponse so we can
	// add the associated private data hash to the http response
	var fabResp VoteFabricResponse
	err = json.Unmarshal([]byte(resp), &fabResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// retrieve private data hash from Fabric ledger
	votePrivateDetailsHash, err := app.FabricSDK.GetVotePrivateDetailsHashSDK(pollID, voterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(reflect.TypeOf(votePrivateDetailsHash))

	// *** TODO ***: decode unicode votePrivateDetailsHash to string

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

	vars := mux.Vars(r)
	pollID := vars["pollid"]
	voterID := vars["voterid"]

	resp, err := app.FabricSDK.GetVotePrivateDetailsSDK(pollID, voterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var fabResp VoteFabricResponsePrivateDetails

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

// Retrieve vote from the Fabric network
func (app *Application) getPollPrivateDetailsHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pollID := vars["pollid"]

	resp, err := app.FabricSDK.GetPollPrivateDetailsSDK(pollID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var fabResp PollFabricResponsePrivateDetails

	err = json.Unmarshal([]byte(resp), &fabResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	pollContentResp, err := app.IpfsGetData(fabResp.PollHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(pollContentResp))
}

func (app *Application) queryVotesByPollHandler(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	pollID := vars["pollid"]

	resp, err := app.FabricSDK.QueryVotesByPollSDK(pollID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}