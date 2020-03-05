package web

import (
	"fmt"
	"net/http"
	"encoding/json"

	// "github.com/gorilla/mux"
)

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

	// Print out values of arguments to initVote
	fmt.Println("poll ID:   ", pollID)
	fmt.Println("voter ID:  ", voterID)
	fmt.Println("voter sex: ", sex)
	fmt.Println("voter age: ", age)
	fmt.Println("CID:       ", cid)

	// Call InitVoteSDK() to initialize a vote on the Fabric network
	resp, err := app.FabricSDK.InitVoteSDK(pollID, voterID, sex, age, cid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
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

	// Print out values of arguments to initVote
	fmt.Println("poll ID:   ", pollID)
	fmt.Println("CID:       ", cid)

	// Call InitVoteSDK() to initialize a vote on the Fabric network
	resp, err := app.FabricSDK.InitPollSDK(pollID, cid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}