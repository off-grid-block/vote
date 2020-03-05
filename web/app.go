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

	/*********************************/
	/*	subrouter for "poll" prefix  */
	/*********************************/
	poll := api.PathPrefix("/poll").Subrouter()

	// handler for getVotePrivateDetails
	poll.HandleFunc("/{pollid:[0-9]+}/{voterid:[0-9]+}/private", app.getVotePrivateDetailsHandler).Methods("GET")

	// handler for getVote
	poll.HandleFunc("/{pollid:[0-9]+}/{voterid:[0-9]+}", app.getVoteHandler).Methods("GET")

	// handler for queryVotesByPoll
	poll.HandleFunc("/{pollid:[0-9]+}/all", app.queryVotesByPollHandler).Methods("GET")

	// handler for getPollPrivateDetails
	poll.HandleFunc("/{pollid:[0-9]+}/private", app.getPollPrivateDetailsHandler).Methods("GET")

	/*********************************/
	/*	subrouter for "voter" prefix */
	/*********************************/
	voter := api.PathPrefix("/voter").Subrouter()

	// handler for getVotesByVoter
	voter.HandleFunc("/{voterid:[0-9]+}/all", app.queryVotesByVoterHandler).Methods("GET")

	/*********************/
	/*  create requests  */
	/*********************/

	// handler for initPoll
	api.HandleFunc("/create/poll", app.initPollHandler).
		Methods("POST").
		Queries(
			"pollid", "{pollid}")

	// handler for initVote
	api.HandleFunc("/create/vote", app.initVoteHandler).
		Methods("POST").
		Queries(
			"pollid", "{pollid}",
			"voterid", "{voterid}",
			"sex", "{sex}", 
			"age", "{age}")

	srv := &http.Server{
		Handler: 	r,
		Addr:		"127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Listening on http://127.0.0.1:8000/...")
	log.Fatal(srv.ListenAndServe())
}