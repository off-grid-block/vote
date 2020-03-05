package web

import (
	"bytes"
	// "context"
	"encoding/json"
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

func (app *Application) IpfsAddPoll(p PollContent) (string, error) {
	var cid string

	pollBytes, err := json.Marshal(p)
	if err != nil {
		return cid, err
	}

	// create io reader of bytes
	reader := bytes.NewReader(pollBytes)

	// add byte data to IPFS
	cid, err = app.IpfsShell.Add(reader)
	if err != nil {
		return cid, err
	}

	return cid, nil
}

func (app *Application) IpfsGetData(cid string) (string, error) {
	tmpFilePath := "/tmp/file"
	app.IpfsShell.Get(cid, tmpFilePath)

	data, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		return "", err
	}
	os.Remove(tmpFilePath)
	return string(data), nil
}