package web

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
)

func (app *Application) IpfsAdd(data interface{}) (string, error) {
	var cid string

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return cid, err
	}

	// create io reader of bytes
	reader := bytes.NewReader(dataBytes)

	// add byte data to IPFS
	cid, err = app.IpfsShell.Add(reader)
	if err != nil {
		return cid, err
	}

	return cid, nil
}


func (app *Application) IpfsGet(cid string) (string, error) {
	tmpFilePath := "/tmp/file"
	app.IpfsShell.Get(cid, tmpFilePath)

	data, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		return "", err
	}
	os.Remove(tmpFilePath)
	return string(data), nil
}
