package web

import (
	"net/http"
	"fmt"
	"encoding/json"	
	
	"github.com/gorilla/mux"
)

type PollContent struct {
	FirstChoice 	bool
	SecondChoice 	bool
	ThirdChoice		bool
}

type FabricResponsePollPrivateDetails struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	Salt 		string 	`json:"salt"`
	PollHash 	string 	`json:"pollHash"`
}

// push poll data to IPFS and IPFS CID to ledger
func (app *Application) initPollHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	pollID := query.Get("pollid")

	var p PollContent

	// Decode HTTP request body and marshal into Vote struct.
	// If the bytes in the request body do not match the fields
	// of the Vote struct, the operation will fail.
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Push vote data to IPFS
	cid, err := app.IpfsAddPoll(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate a random salt to concatenate with the vote's IPFS CID
	salt := GenerateSalt()

	// Print out values of arguments to initVote
	fmt.Println("poll ID:   ", pollID)
	fmt.Println("salt:      ", salt)
	fmt.Println("CID:       ", cid)

	// Call InitVoteSDK() to initialize a vote on the Fabric network
	resp, err := app.FabricSDK.InitPollSDK(pollID, salt, cid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

// Retrieve vote from the Fabric network
func (app *Application) getPollPrivateDetailsHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	pollID := query.Get("pollid")
	salt := query.Get("salt")

	resp, err := app.FabricSDK.GetPollPrivateDetailsSDK(pollID, salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var fabResp FabricResponsePollPrivateDetails

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