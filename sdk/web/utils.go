package web

import (
	"time"
	"bytes"
	// "context"
	"strings"
	"encoding/json"
	"math/rand"
	"io/ioutil"
	"os"
	// ipfs "github.com/ipfs/go-ipfs-api"
)

func (app *Application) IpfsAddVote(v VoteContent) (string, error) {
	var cid string

	voteBytes, err := json.Marshal(v)
	if err != nil {
		return cid, err
	}

	// create io reader of bytes
	reader := bytes.NewReader(voteBytes)

	// add byte data to IPFS
	cid, err = app.IpfsShell.Add(reader)
	if err != nil {
		return cid, err
	}

	return cid, nil
}

// func (app *Application) IpfsGetVote(cid string) (*ipfs.Response, error) {
// 	var resp *ipfs.Response

// 	resp, err := app.IpfsShell.Request("get", cid).Option("create", true).Send(context.Background())
// 	if err != nil {
// 		return resp, err
// 	}

// 	if resp.Error != nil {
// 		return resp, resp.Error
// 	}

// 	return resp, nil
// }

func (app *Application) IpfsGetVote(cid string) (string, error) {

	tmpFilePath := "/tmp/vote" + GenerateSalt()
	app.IpfsShell.Get(cid, tmpFilePath)

	data, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		return "", err
	}
	os.Remove(tmpFilePath)
	return string(data), nil

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