package web

import (
	"net/http"
	"encoding/json"

	"fmt"
	
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
	var fabResp voteResponseSDK
	err = json.Unmarshal([]byte(resp), &fabResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var voteContentResp interface{}

	resp, err = app.FabricSDK.GetVotePrivateDetailsSDK(pollID, voterID)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("You do not have permission to see these vote details")
	} else {
		var fabPrivateResp votePrivateDetailsResponseSDK

		err = json.Unmarshal([]byte(resp), &fabPrivateResp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		voteContentResp, err = app.IpfsGet(fabPrivateResp.VoteHash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var vote voteDetailsHttpResponse

	vote.PollID = fabResp.PollID
	vote.VoterID = fabResp.VoterID
	vote.VoterSex = fabResp.VoterSex
	vote.VoterAge = fabResp.VoterAge
	vote.Content = voteContentResp

	httpResp, err := json.Marshal(vote)
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

	var fabResp votePrivateDetailsResponseSDK

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