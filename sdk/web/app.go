package web

import (
	"github.com/off-grid-block/vote/sdk/blockchain"
	"net/http"
	// "io/ioutil"
	"log"
	"fmt"
	"time"
	// "encoding/json"
	"github.com/gorilla/mux"
	ipfs "github.com/ipfs/go-ipfs-api"
)

// Struct containing Fabric SDK setup data. Objects of type
// Application have access to the SDK's chaincode interface.
type Application struct {
	FabricSDK *blockchain.SetupSDK
	IpfsShell *ipfs.Shell
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

	// api.HandleFunc("/vote/private/public-hash", app.getVotePrivateDetailsHandler).
	// 	Methods("GET").
	// 	Queries(
	// 		"pollid", "{pollid}",
	// 		"voterid", "{voterid}",
	// 		"salt", "{salt}")

	api.HandleFunc("/vote/private", app.getVotePrivateDetailsHandler).
		Methods("GET").
		Queries(
			"pollid", "{pollid}",
			"voterid", "{voterid}",
			"salt", "{salt}")

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
			"voterid", "{voterid}",
			"salt", "{salt}")

	api.HandleFunc("/poll/{pollid:[0-9]+}/all", app.queryVotesByPollHandler).Methods("GET")

	// handler for getPollPrivateDetails
	api.HandleFunc("/poll/private", app.getPollPrivateDetailsHandler).
		Methods("GET").
		Queries("pollid", "{pollid}",
				"salt", "{salt}")

	// handler for initPoll
	api.HandleFunc("/poll", app.initPollHandler).
		Methods("POST").
		Queries("pollid", "{pollid}")

	srv := &http.Server{
		Handler: 	r,
		Addr:		"127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Listening on http://127.0.0.1:8000/...")
	log.Fatal(srv.ListenAndServe())
}