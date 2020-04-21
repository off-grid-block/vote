package web

import (
	"fmt"
	"net/http"
	"encoding/json"
	
	"github.com/gorilla/mux"
)


func (app *Application) getPollHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pollID := vars["pollid"]

	resp, err := app.FabricSDK.GetPollSDK(pollID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

// Retrieve vote from the Fabric network
func (app *Application) getPollPrivateDetailsHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pollID := vars["pollid"]

	resp, err := app.FabricSDK.GetPollPrivateDetailsSDK(pollID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var fabResp pollPrivateDetailsResponseSDK

	err = json.Unmarshal([]byte(resp), &fabResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	pollContentResp, err := app.IpfsGet(fabResp.PollHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(pollContentResp))
}