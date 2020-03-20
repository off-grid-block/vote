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
	PollID 			string
	VoterID 		string
	Sex 			string
	Age 			string
	Content 		struct {
		YesOrNo 		bool
	}
}

type InitPollRequestBody struct {
	PollID 			string
	Content 		struct {
		FirstChoice 	bool
		SecondChoice 	bool
		ThirdChoice		bool	
	}
}

type UpdatePollStatusRequestBody struct {
	Status 			string 	`json:"status"`
}

type VoteResponseSDK struct {
	ObjectType 		string 	`json:"docType"`
	PollID			string 	`json:"pollID"`
	VoterID			string 	`json:"voterID"`
	VoterSex 		string 	`json:"voterSex"`
	VoterAge		int 	`json:"voterAge"`
	PrivateHash 	string 	`json:"privateHash"`
}

type VotePrivateDetailsResponseSDK struct {
	ObjectType 		string 	`json:"docType"`
	PollID			string 	`json:"pollID"`
	VoterID			string 	`json:"voterID"`
	Salt 			string 	`json:"salt"`
	VoteHash 		string 	`json:"voteHash"`
}

type PollPrivateDetailsResponseSDK struct {
	ObjectType 		string 	`json:"docType"`
	PollID			string 	`json:"pollID"`
	Salt 			string 	`json:"salt"`
	PollHash 		string 	`json:"pollHash"`
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

	// handler for initPoll
	poll.HandleFunc("/create", app.initPollHandler).Methods("POST")

	// handler for updatePollStatus
	poll.HandleFunc("/{pollid:[0-9]+}/status", app.updatePollStatusHandler).Methods("PUT")

	// handler for getPoll
	poll.HandleFunc("/{pollid:[0-9]+}", app.getPollHandler).Methods("GET")

	// handler for getPollPrivateDetails
	poll.HandleFunc("/{pollid:[0-9]+}/private", app.getPollPrivateDetailsHandler).Methods("GET")

	/*********************************/
	/*	subrouter for "vote" prefix  */
	/*********************************/
	vote := api.PathPrefix("/vote").Subrouter()

	// handler for initVote
	vote.HandleFunc("/create", app.initVoteHandler).Methods("POST")

	// handler for getVotePrivateDetails
	vote.HandleFunc("/{pollid:[0-9]+}/{voterid:[0-9]+}/private", app.getVotePrivateDetailsHandler).Methods("GET")

	// handler for getVotePrivateDetailsHash
	vote.HandleFunc("/{pollid:[0-9]+}/{voterid:[0-9]+}/hash", app.getVotePrivateDetailsHashHandler).Methods("GET")

	// handler for getVote
	vote.HandleFunc("/{pollid:[0-9]+}/{voterid:[0-9]+}", app.getVoteHandler).Methods("GET")

	// handler for queryVotePrivateDetailsByPoll
	vote.HandleFunc("", app.queryVotePrivateDetailsByPollHandler).
		Methods("GET").
		Queries("type", "private", "pollid", "{pollid}")

	// handler for queryVotesByPoll
	vote.HandleFunc("", app.queryVotesByPollHandler).
		Methods("GET").
		Queries("type", "public", "pollid", "{pollid}")

	// handler for getVotesByVoter
	vote.HandleFunc("", app.queryVotesByVoterHandler).
		Methods("GET").
		Queries("voter", "{voterid}")

	srv := &http.Server{
		Handler: 	r,
		Addr:		"127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Listening on http://127.0.0.1:8000/...")
	log.Fatal(srv.ListenAndServe())
}