package web

import (
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"
	"encoding/gob"
	
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

func (app *Application) getPollHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pollID := vars["pollid"]

	resp, err := app.FabricSDK.GetPollSDK(pollID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
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

	pollContentResp, err := app.IpfsGet(fabResp.PollHash)
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


func (app *Application) queryVotePrivateDetailsByPollHandler(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	pollID := vars["pollid"]

	cidList, err := app.FabricSDK.QueryVotePrivateDetailsByPollSDK(pollID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var votesByPoll []string

	for _, cid := range cidList {
		vote, err := app.IpfsGet(cid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		votesByPoll = append(votesByPoll, vote)
	}

	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err = enc.Encode(votesByPoll)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(buf.Bytes())
}

func (app *Application) updatePollStatusHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pollID := vars["pollid"]

	var body UpdatePollStatusRequestBody

	// Decode HTTP request body and marshal into Vote struct.
	// If the bytes in the request body do not match the fields
	// of the Vote struct, the operation will fail.
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := app.FabricSDK.UpdatePollStatusSDK(pollID, body.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}