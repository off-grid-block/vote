package voteapp

import (
    "fmt"
    "github.com/off-grid-block/fabric-sdk-go/pkg/client/channel"
    // "github.com/off-grid-block/fabric-sdk-go/pkg/common/providers/fab"
    "github.com/off-grid-block/vote/blockchain"
    "bytes"
    "encoding/gob"
)


// query the private details of a vote by poll
func QueryVotePrivateDetailsByPollSDK(s *blockchain.SetupSDK, pollID string) ([]string, error) {

    var cidList []string

    response, err := s.Client.Query(channel.Request{ChaincodeID: "vote", Fcn: "queryVotePrivateDetailsByPoll",  Args: [][]byte{[]byte(pollID)}})
    if err != nil {
        return cidList, fmt.Errorf("failed to query: %v", err)
    }

    buf := bytes.NewBuffer(response.Payload)
    dec := gob.NewDecoder(buf)
    err = dec.Decode(&cidList)
    if err != nil {
        return cidList, fmt.Errorf("failed to query vote private details by poll: %v\n", err)
    }

    return cidList, nil
}

// query votes of a particular poll
func QueryVotesByPollSDK(s *blockchain.SetupSDK, pollID string) (string, error) {

    response, err := s.Client.Query(channel.Request{ChaincodeID: "vote", Fcn: "queryVotesByPoll", Args: [][]byte{[]byte(pollID)}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}

// query votes of a particular poll
func QueryVotesByVoterSDK(s *blockchain.SetupSDK, voterID string) (string, error) {

    response, err := s.Client.Query(channel.Request{ChaincodeID: "vote", Fcn: "queryVotesByVoter", Args: [][]byte{[]byte(voterID)}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}