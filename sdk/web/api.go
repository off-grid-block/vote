package web

import (
	"github.com/off-grid-block/vote/blockchain"
	"net/http"
	"log"
	"github.com/gorilla/mux"
)

type Application struct {
	FabricSDK *blockchain.SetupSDK
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test successful.\n"))
}

func Serve(app *Application) {
	r := mux.NewRouter()
	r.HandleFunc("/", TestHandler)
	log.Fatal(http.ListenAndServe(":8000", r))
}