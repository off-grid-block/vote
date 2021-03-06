package web

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/off-grid-block/vote/blockchain"
	"log"
	"net/http"
	"os"
	"time"
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

type checkAgentResponse struct {
	Initialized bool `json:"initialized"`
}

func checkAgentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {

		req, err := http.NewRequest("GET", os.Getenv("CORE_URL") + "/admin/agent", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		q := req.URL.Query()
		q.Add("alias", "vote")
		req.URL.RawQuery = q.Encode()

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var response checkAgentResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !response.Initialized {
			http.Error(w, "agent has not yet been initialized", http.StatusBadRequest)
		} else {
			fmt.Println("Agent is online. Passing on request to agent...")
			next.ServeHTTP(w, r)
		}
	})
}


// Homepage
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Homepage\n"))
}

// Initiate the web server
func Serve(app *Application) {

	// create new router and register middlewares
	r := mux.NewRouter()
	r.Use(checkAgentMiddleware)

	api := r.PathPrefix("/api/v1").Subrouter()

	// test api homepage
	api.HandleFunc("/", HomeHandler)

	/*********************************/
	/*	subrouter for "poll" prefix  */
	/*********************************/
	poll := api.PathPrefix("/vote-app/poll").Subrouter()

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
	vote := api.PathPrefix("/vote-app/vote").Subrouter()

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
		Addr:		os.Getenv("API_URL"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Listening on %v...\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
