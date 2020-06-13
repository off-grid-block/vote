package web

import (
	"bytes"
	"net/http"
	"encoding/gob"
	"github.com/off-grid-block/vote/blockchain"
	"github.com/gorilla/mux"
)


func (app *Application) QueryVotesByPollHandler(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	pollID := vars["pollid"]

	resp, err := blockchain.QueryVotesByPollSDK(app.FabricSDK, pollID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}


func (app *Application) QueryVotePrivateDetailsByPollHandler(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	pollID := vars["pollid"]

	cidList, err := blockchain.QueryVotePrivateDetailsByPollSDK(app.FabricSDK, pollID)
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


func (app *Application) QueryVotesByVoterHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	voterID := vars["voterid"]
	resp, err := blockchain.QueryVotesByVoterSDK(app.FabricSDK, voterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))

}