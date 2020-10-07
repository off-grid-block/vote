package core

import (
	caMsp "github.com/off-grid-block/fabric-sdk-go/pkg/client/msp"
	"github.com/pkg/errors"
	//"github.com/off-grid-block/fabric-sdk-go/pkg/client/channel"
	"fmt"
	"github.com/off-grid-block/vote/blockchain"
)

// InvokeHello
func RegUser(s *blockchain.SetupSDK, data caMsp.RegistrationRequest) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	// new User information

	caClient, err := caMsp.New(s.Fsdk.Context())
	//fmt.Println("caclient", caClient)
	enrollSecret, err := caClient.Register(&data)
	fmt.Println("enrollSecret", enrollSecret)
	if err != nil {
		return "", errors.WithMessage(err, "Unable to register user with CA")
	}

	return string(enrollSecret), nil
}
