package web

import (
	"github.com/off-grid-block/vote/sdk/blockchain"
	"net/http"
	"log"
	"fmt"
	"time"
	"github.com/gorilla/mux"
	ipfs "github.com/ipfs/go-ipfs-api"
)

// Struct containing Fabric SDK setup data. Objects of type
// Application have access to the SDK's chaincode interface.
type Application struct {
	FabricSDK *blockchain.SetupSDK
	IpfsShell *ipfs.Shell
} 

// PRELIMINARY DEF: struct to hold vote data.
type InitVoteRequestBody struct {
	YesOrNo bool
}

type InitPollRequestBody struct {
	FirstChoice 	bool
	SecondChoice 	bool
	ThirdChoice		bool
}

type UpdatePollStatusRequestBody struct {
	Status 		string 	`json:"status"`
}

type VoteFabricResponse struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	VoterID		string 	`json:"voterID"`
	VoterSex 	string 	`json:"voterSex"`
	VoterAge	int 	`json:"voterAge"`
	PrivateHash string 	`json:"privateHash"`
}

type VoteFabricResponsePrivateDetails struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	VoterID		string 	`json:"voterID"`
	Salt 		string 	`json:"salt"`
	VoteHash 	string 	`json:"voteHash"`
}

type PollFabricResponsePrivateDetails struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	Salt 		string 	`json:"salt"`
	PollHash 	string 	`json:"pollHash"`
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

	// handler for getVotePrivateDetailsHash
	poll.HandleFunc("/{pollid:[0-9]+}/{voterid:[0-9]+}/hash", app.getVotePrivateDetailsHashHandler).Methods("GET")

	// handler for getVote
	poll.HandleFunc("/{pollid:[0-9]+}/{voterid:[0-9]+}", app.getVoteHandler).Methods("GET")

	// handler for queryVotePrivateDetailsByPoll
	poll.HandleFunc("/{pollid:[0-9]+}/all/private", app.queryVotePrivateDetailsByPollHandler).Methods("GET")

	// handler for queryVotesByPoll
	poll.HandleFunc("/{pollid:[0-9]+}/all", app.queryVotesByPollHandler).Methods("GET")

	// handler for getPoll
	poll.HandleFunc("/{pollid:[0-9]+}/summary", app.getPollHandler).Methods("GET")

	// handler for getPollPrivateDetails
	poll.HandleFunc("/{pollid:[0-9]+}/private", app.getPollPrivateDetailsHandler).Methods("GET")

	// handler for updatePollStatus
	poll.HandleFunc("/{pollid:[0-9]+}/status", app.updatePollStatusHandler).Methods("POST")

	/*********************************/
	/*	subrouter for "voter" prefix */
	/*********************************/
	voter := api.PathPrefix("/voter").Subrouter()

	// handler for getVotesByVoter
	voter.HandleFunc("/{voterid:[0-9]+}/all", app.queryVotesByVoterHandler).Methods("GET")

	/*********************/
	/*  create requests  */
	/*********************/

	// handler for initVote
	api.HandleFunc("/create/vote", app.initVoteHandler).
		Methods("POST").
		Queries(
			"pollid", "{pollid}",
			"voterid", "{voterid}",
			"sex", "{sex}", 
			"age", "{age}")

	// handler for initPoll
	api.HandleFunc("/create/poll", app.initPollHandler).
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