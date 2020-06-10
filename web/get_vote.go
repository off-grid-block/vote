package web

import (
	"net/http"
	"encoding/json"

	"fmt"
	
	"github.com/gorilla/mux"
)


// Retrieve vote from the Fabric network
func (app *Application) GetVoteHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pollID := vars["pollid"]
	voterID := vars["voterid"]

	resp, err := app.FabricSDK.GetVoteSDK(pollID, voterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Unmarshal into our defined struct
	var fabResp voteResponseSDK
	err = json.Unmarshal([]byte(resp), &fabResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve contents of vote
	var voteContentResp interface{}

	resp, err = app.FabricSDK.GetVotePrivateDetailsSDK(pollID, voterID)
	if err != nil { // if no access, simply ignore request
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

	// Retrieve private data hash
	hash, err := app.FabricSDK.GetVotePrivateDetailsHashSDK(pollID, voterID)
	if err != nil { // if user doesn't have access, ignore request
		fmt.Println("You do not have permission to see the private data hash")
		hash = ""
	}

	var vote voteDetailsHttpResponse

	vote.PollID = fabResp.PollID
	vote.VoterID = fabResp.VoterID
	vote.VoterSex = fabResp.VoterSex
	vote.VoterAge = fabResp.VoterAge
	vote.Content = voteContentResp
	vote.PrivateHash = hash

	httpResp, err := json.Marshal(vote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(httpResp))
}

// // Retrieve private details of a vote from the Fabric network
// func (app *Application) getVotePrivateDetailsHandler(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	pollID := vars["pollid"]
// 	voterID := vars["voterid"]

// 	resp, err := app.FabricSDK.GetVotePrivateDetailsSDK(pollID, voterID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	var fabResp votePrivateDetailsResponseSDK

// 	err = json.Unmarshal([]byte(resp), &fabResp)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	voteContentResp, err := app.IpfsGet(fabResp.VoteHash)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Write([]byte(voteContentResp))
// }