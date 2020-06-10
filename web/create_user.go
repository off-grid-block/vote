package web

import (
	"encoding/json"
	"net/http"
	"strings"
	"fmt"
	caMsp "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
)

func (app *Application) UserHandler(w http.ResponseWriter, r *http.Request) {
	// data := &struct {
	// 	Name 		string
	// 	Secret      string
	// 	Type      	string
	// }{
	// 	Name: 		"",
	// 	Secret:     "",
	// 	Type:      	"",
	// }
	affl := strings.ToLower("org1") + ".department1"

	data := caMsp.RegistrationRequest{
		Name:           "email",
		Secret:         "password",
		Type:           "peer",
		MaxEnrollments: -1,
		Affiliation:    affl,
		Attributes: []caMsp.Attribute{
			{
				Name:  "role",
				Value: "user",
				ECert: true,
			},
		},
		CAName: "ca.org1.hf.sample.io",
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	_, err = app.FabricSDK.RegUser(data)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to invoke hello in the blockchain", 500)
	}
//	fmt.Println(txid)
	// data.TransactionId = txid
	// data.Success = true
	// data.Response = true
}
