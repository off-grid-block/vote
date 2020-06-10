package web

import (
	"net/http"
	"encoding/json"
)


// push poll data to IPFS and IPFS CID to ledger
func (app *Application) InitPollHandler(w http.ResponseWriter, r *http.Request) {

	var p initPollRequestBodyAPI

	// Decode HTTP request body and marshal into Poll struct.
	// If the bytes in the request body do not match the fields
	// of the Poll struct, the operation will fail.
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Push poll data to IPFS
	cid, err := app.IpfsAdd(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Call InitVoteSDK() to initialize a vote on the Fabric network
	resp, err := app.FabricSDK.InitPollSDK(p.PollID, p.Title, cid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}