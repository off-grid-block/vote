package web

import (
	"net/http"
	"encoding/json"
	"github.com/off-grid-block/vote/voteapp"
)

// Initialize & push votes on the Fabric network
func (app *Application) initVoteHandler(w http.ResponseWriter, r *http.Request) {

	var v initVoteRequestBodyAPI

	// Decode HTTP request body and marshal into Vote struct.
	// If the bytes in the request body do not match the fields
	// of the Vote struct, the operation will fail.
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Push vote data to IPFS
	cid, err := app.IpfsAdd(v.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Call InitVoteSDK() to initialize a vote on the Fabric network
	resp, err := voteapp.InitVoteSDK(app.FabricSDK, v.PollID, v.VoterID, v.Sex, v.Age, cid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}