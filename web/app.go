package web

import (
	"github.com/off-grid-block/vote/blockchain"
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
type initVoteRequestBodyAPI struct {
	PollID 			string 		`json:"pollID"`
	VoterID 		string 		`json:"voterID"`
	Sex 			string 		`json:"sex"`
	Age 			string 		`json:"age"`
	Content 		interface{} `json:"content"`
}

type initPollRequestBodyAPI struct {
	PollID 			string 		`json:"pollID"`
	Title 			string 		`json:"title"`
	Content 		interface{} `json:"content"`
}

type updatePollStatusRequestBodyAPI struct {
	Status 			string 		`json:"status"`
}

type voteDetailsHttpResponse struct {
	PollID 			string 		`json:"pollID"`
	VoterID 		string 		`json:"voterID"`
	VoterSex 		string 		`json:"voterSex"`
	VoterAge 		int 		`json:"voterAge"`
	Content 		interface{} `json:"content"`
	PrivateHash 	string 		`json:"privateHash"`
}

type pollDetailsHttpResponse struct {
	PollID 			string
	Title 			string
	Status 			string
	NumVotes 		int
	Content 		interface{}
}

type voteResponseSDK struct {
	ObjectType 		string 		`json:"docType"`
	PollID			string 		`json:"pollID"`
	VoterID			string 		`json:"voterID"`
	VoterSex 		string 		`json:"voterSex"`
	VoterAge		int 		`json:"voterAge"`
	PrivateHash 	string 		`json:"privateHash"`
}

type votePrivateDetailsResponseSDK struct {
	ObjectType 		string 		`json:"docType"`
	PollID			string 		`json:"pollID"`
	VoterID			string 		`json:"voterID"`
	Salt 			string 		`json:"salt"`
	VoteHash 		string 		`json:"voteHash"`
}

type pollResponseSDK struct {
	ObjectType 		string
	PollID 			string
	Title 			string
	Status  		string
	NumVotes 		int
}

type pollPrivateDetailsResponseSDK struct {
	ObjectType 		string 		`json:"docType"`
	PollID			string 		`json:"pollID"`
	Salt 			string 		`json:"salt"`
	PollHash 		string 		`json:"pollHash"`
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

	/********************************/
	/* identity management endpoint */
	/********************************/
	api.HandleFunc("/application", app.userHandler).Methods("POST")

	/*********************************/
	/*	subrouter for "poll" prefix  */
	/*********************************/
	poll := api.PathPrefix("/poll").Subrouter()

	// handler for initPoll
	poll.HandleFunc("", app.initPollHandler).Methods("POST")

	// handler for queryAllPolls
	poll.HandleFunc("", app.queryAllPollsHandler).Methods("GET")

	// handler for updatePollStatus
	poll.HandleFunc("/{pollid}/status", app.updatePollStatusHandler).Methods("PUT")

	// handler for getPoll
	poll.HandleFunc("/{pollid}", app.getPollHandler).Methods("GET")

	// // handler for getPollPrivateDetails
	// poll.HandleFunc("/{pollid}/private", app.getPollPrivateDetailsHandler).Methods("GET")

	/*********************************/
	/*	subrouter for "vote" prefix  */
	/*********************************/
	vote := api.PathPrefix("/vote").Subrouter()

	// handler for initVote
	vote.HandleFunc("", app.initVoteHandler).Methods("POST")

	// // handler for getVotePrivateDetails
	// vote.HandleFunc("/{pollid}/{voterid}/private", app.getVotePrivateDetailsHandler).Methods("GET")

	// // handler for getVotePrivateDetailsHash
	// vote.HandleFunc("/{pollid}/{voterid}/hash", app.getVotePrivateDetailsHashHandler).Methods("GET")

	// handler for getVote
	vote.HandleFunc("/{pollid}/{voterid}", app.getVoteHandler).Methods("GET")

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
		Queries("voterid", "{voterid}")

	srv := &http.Server{
		Handler: 	r,
		Addr:		"127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Listening on http://127.0.0.1:8000/...")
	log.Fatal(srv.ListenAndServe())
}