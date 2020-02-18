package web

import (
	"github.com/off-grid-block/vote/sdk/blockchain"
	"net/http"
	"log"
	"fmt"
	"time"
	"bytes"
	"strings"
	"math/rand"
	"encoding/json"
	"github.com/gorilla/mux"
	ipfs "github.com/ipfs/go-ipfs-api"
)

type Application struct {
	FabricSDK *blockchain.SetupSDK
}

type Vote struct {
	YesOrNo string
}


func IpfsAddVote(v Vote) (string, error) {
	var cid string

	voteBytes, err := json.Marshal(v)
	if err != nil {
		return cid, err
	}

	// create io reader of bytes
	reader := bytes.NewReader(voteBytes)

	// create shell to connect to IPFS
	sh := ipfs.NewShell("localhost:5001")

	// add byte data to IPFS
	cid, err = sh.Add(reader)
	if err != nil {
		return cid, err
	}

	return cid, nil
}

// https://yourbasic.org/golang/generate-random-string/
func GenerateSalt() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune(
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" + 
		"abcdefghijklmnopqrstuvwxyz" + 
		"0123456789")
	length := 8

	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	return b.String()
}


func (app *Application) initVoteHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	pollID := query.Get("pollid")
	voterID := query.Get("voterid")
	sex := query.Get("sex")
	age := query.Get("age")

	var v Vote

	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cid, err := IpfsAddVote(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	salt := GenerateSalt()

	resp, err := app.FabricSDK.InitVoteSDK(pollID, voterID, sex, age, salt, cid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

func (app *Application) getVoteHandler(w http.ResponseWriter, r *http.Request) {

}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test successful.\n"))
}

func Serve(app *Application) {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	// test api homepage
	api.HandleFunc("/", TestHandler)

	// handle for postVote
	api.HandleFunc("/vote/", app.initVoteHandler).
		Methods("POST").
		Queries(
			"pollid", "{pollid}", 
			"voterid", "{voterid}", 
			"sex", "{sex}", 
			"age", "{age}")

	// handler for getVote
	api.HandleFunc("/vote/", app.getVoteHandler).
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