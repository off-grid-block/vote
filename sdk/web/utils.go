package web

import (
	"time"
	"bytes"
	"strings"
	"encoding/json"
	"math/rand"
	ipfs "github.com/ipfs/go-ipfs-api"
)

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