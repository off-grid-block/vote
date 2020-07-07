package web

import (
	"net/http"
	"encoding/json"
	
	"github.com/gorilla/mux"
	"github.com/off-grid-block/vote/voteapp"
)


func (app *Application) updatePollStatusHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pollID := vars["pollid"]

	var body updatePollStatusRequestBodyAPI

	// Decode HTTP request body and marshal into Vote struct.
	// If the bytes in the request body do not match the fields
	// of the Vote struct, the operation will fail.
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := voteapp.UpdatePollStatusSDK(app.FabricSDK, pollID, body.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}