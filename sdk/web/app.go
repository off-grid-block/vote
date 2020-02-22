package web

import (
	"github.com/off-grid-block/vote/sdk/blockchain"
	"net/http"
	"log"
	"fmt"
	"time"
	"encoding/json"
	"github.com/gorilla/mux"
)

// Struct containing Fabric SDK setup data. Objects of type
// Application have access to the SDK's chaincode interface.
type Application struct {
	FabricSDK *blockchain.SetupSDK
} 

// PRELIMINARY DEF: struct to hold vote data.
type Vote struct {
	YesOrNo bool
}

// Initialize & push votes on the Fabric network
func (app *Application) initVoteHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	pollID := query.Get("pollid")
	voterID := query.Get("voterid")
	sex := query.Get("sex")
	age := query.Get("age")

	var v Vote

	// Decode HTTP request body and marshal into Vote struct.
	// If the bytes in the request body do not match the fields
	// of the Vote struct, the operation will fail.
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Push vote data to IPFS
	cid, err := IpfsAddVote(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate a random salt to concatenate with the vote's IPFS CID
	salt := GenerateSalt()

	// Print out values of arguments to initVote
	fmt.Println("poll ID:   ", pollID)
	fmt.Println("voter ID:  ", voterID)
	fmt.Println("voter sex: ", sex)
	fmt.Println("voter age: ", age)
	fmt.Println("salt:      ", salt)
	fmt.Println("CID:       ", cid)

	// Call InitVoteSDK() to initialize a vote on the Fabric network
	resp, err := app.FabricSDK.InitVoteSDK(pollID, voterID, sex, age, salt, cid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

// Retrieve vote from the Fabric network
func (app *Application) getVoteHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	pollID := query.Get("pollid")
	voterID := query.Get("voterid")

	resp, err := app.FabricSDK.GetVoteSDK(pollID, voterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

// Retrieve private details of a vote from the Fabric network
func (app *Application) getVotePrivateDetailsHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	pollID := query.Get("pollid")
	voterID := query.Get("voterid")
	salt := query.Get("salt")

	resp, err := app.FabricSDK.GetVotePrivateDetailsSDK(pollID, voterID, salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

// Homepage
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Homepage\n"))
}

// Initiate the web server
func Serve(app *Application) {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	// test api homepage
	api.HandleFunc("/", HomeHandler)

	api.HandleFunc("/vote/private", app.getVotePrivateDetailsHandler).
		Methods("GET").
		Queries(
			"pollid", "{pollid}",
			"voterid", "{voterid}")

	// handle for postVote
	api.HandleFunc("/vote", app.initVoteHandler).
		Methods("POST").
		Queries(
			"pollid", "{pollid}", 
			"voterid", "{voterid}", 
			"sex", "{sex}", 
			"age", "{age}")

	// handler for getVote
	api.HandleFunc("/vote", app.getVoteHandler).
		Methods("GET").
		Queries(
			"pollid", "{pollid}", 
			"voterid", "{voterid}")

	srv := &http.Server{
		Handler: 	r,
		Addr:		"127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Listening on http://127.0.0.1:8000/...")
	log.Fatal(srv.ListenAndServe())
}