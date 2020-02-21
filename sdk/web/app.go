package web

import (
	"github.com/off-grid-block/vote/sdk/blockchain"
	"net/http"
	"log"
	"fmt"
	"time"
	"reflect"
	"encoding/json"
	"github.com/gorilla/mux"
)

type Application struct {
	FabricSDK *blockchain.SetupSDK
}

type Vote struct {
	YesOrNo bool
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

	fmt.Println(reflect.TypeOf(pollID))
	fmt.Println(reflect.TypeOf(voterID))
	fmt.Println(reflect.TypeOf(sex))
	fmt.Println(reflect.TypeOf(age))
	fmt.Println(salt)
	fmt.Println(cid)
	// return

	resp, err := app.FabricSDK.InitVoteSDK(pollID, voterID, sex, age, salt, cid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

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

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Homepage\n"))
}

func Serve(app *Application) {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	// test api homepage
	api.HandleFunc("/", HomeHandler)

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