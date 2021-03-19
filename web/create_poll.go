package web

import (
	"encoding/json"
	"github.com/off-grid-block/vote/blockchain"
	"net/http"
	"time"
)


// push poll data to IPFS and IPFS CID to ledger
func (app *Application) initPollHandler(w http.ResponseWriter, r *http.Request) {

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

	// Check if the time format is correct
	_, err = time.Parse("2006-Jan-02", p.CloseDate)
	if err != nil {
		http.Error(w, "failed to deserialize close date: " + err.Error(), http.StatusInternalServerError)
	}

	// Call InitVoteSDK() to initialize a vote on the Fabric network
	resp, err := blockchain.InitPollSDK(app.FabricSDK, p.PollID, p.Title, cid, p.CloseDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}
