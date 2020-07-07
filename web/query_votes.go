package web

import (
	"bytes"
	"net/http"
	"encoding/gob"
	
	"github.com/gorilla/mux"
	"github.com/off-grid-block/vote/voteapp"
)


func (app *Application) queryVotesByPollHandler(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	pollID := vars["pollid"]

	resp, err := voteapp.QueryVotesByPollSDK(app.FabricSDK, pollID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}


func (app *Application) queryVotePrivateDetailsByPollHandler(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	pollID := vars["pollid"]

	cidList, err := voteapp.QueryVotePrivateDetailsByPollSDK(app.FabricSDK, pollID)
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


func (app *Application) queryVotesByVoterHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	voterID := vars["voterid"]

	resp, err := voteapp.QueryVotesByVoterSDK(app.FabricSDK, voterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))

}