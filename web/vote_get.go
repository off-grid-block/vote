package web

import (
	"net/http"
	"encoding/json"
	
	"github.com/gorilla/mux"
)


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
	var fabResp VoteResponseSDK
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

	var fabResp VotePrivateDetailsResponseSDK

	err = json.Unmarshal([]byte(resp), &fabResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	voteContentResp, err := app.IpfsGet(fabResp.VoteHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(voteContentResp))
}

func (app *Application) getVotePrivateDetailsHashHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pollID := vars["pollid"]
	voterID := vars["voterid"]

	resp, err := app.FabricSDK.GetVotePrivateDetailsHashSDK(pollID, voterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}