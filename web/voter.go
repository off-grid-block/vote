package web

import (
	"net/http"
	// "fmt"
	// "encoding/json"
	// "unicode/utf8"
	"github.com/gorilla/mux"
	// "reflect"
)

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