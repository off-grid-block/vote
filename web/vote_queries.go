package web

import (
	"bytes"
	"net/http"
	"encoding/gob"
	
	"github.com/gorilla/mux"
)


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


func (app *Application) queryVotesByVoterHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	voterID := vars["voterid"]

	resp, err := app.FabricSDK.QueryVotesByVoterSDK(voterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))

}