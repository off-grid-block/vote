package web

import (
	"fmt"
	"net/http"
	"encoding/json"
	
	"github.com/gorilla/mux"
	"github.com/off-grid-block/vote/blockchain"
)


func (app *Application) getPollHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pollID := vars["pollid"]

	// Retrieve public details from Fabric
	public, err := blockchain.GetPollSDK(app.FabricSDK, pollID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Unmarshal public details
	var fabPublicResp pollResponseSDK

	err = json.Unmarshal([]byte(public), &fabPublicResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var content interface{}

	private, err := blockchain.GetPollPrivateDetailsSDK(app.FabricSDK, pollID)
	// if there is an error, that means the peer does not have access to
	// the private details. So only proceed with the retrieval of private
	// data from IPFS if GetPollPrivateDetailsSDK succeeds.
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("You do not have permission")
	} else {
		var fabPrivateResp pollPrivateDetailsResponseSDK

		// Retrieve IPFS CID from private details
		err = json.Unmarshal([]byte(private), &fabPrivateResp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Retrieve poll content from IPFS
		content, err = app.IpfsGet(fabPrivateResp.PollHash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Build HTTP response
	var pollResp pollDetailsHttpResponse

	pollResp.PollID = fabPublicResp.PollID
	pollResp.Title = fabPublicResp.Title
	pollResp.Status = fabPublicResp.Status
	pollResp.NumVotes = fabPublicResp.NumVotes
	pollResp.Content = content

	resp, err := json.Marshal(pollResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Write([]byte(resp))
}

// // Retrieve vote from the Fabric network
// func (app *Application) getPollPrivateDetailsHandler(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	pollID := vars["pollid"]

// 	resp, err := blockchain.GetPollPrivateDetailsSDK(app.FabricSDK, pollID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	var fabResp pollPrivateDetailsResponseSDK

// 	err = json.Unmarshal([]byte(resp), &fabResp)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		fmt.Println(err)
// 		return
// 	}

// 	pollContentResp, err := app.IpfsGet(fabResp.PollHash)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Write([]byte(pollContentResp))
// }

func (app *Application) queryAllPollsHandler(w http.ResponseWriter, r *http.Request) {

	resp, err := blockchain.QueryAllPollsSDK(app.FabricSDK)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}
