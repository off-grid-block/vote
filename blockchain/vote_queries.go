package blockchain

import (
    "fmt"
    "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
    // "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
    "bytes"
    "encoding/gob"
)


// query the private details of a vote by poll
func (s *SetupSDK) QueryVotePrivateDetailsByPollSDK(pollID string) ([]string, error) {

    var cidList []string

    response, err := s.client.Query(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "queryVotePrivateDetailsByPoll",  Args: [][]byte{[]byte(pollID)}})
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
func (s *SetupSDK) QueryVotesByPollSDK(pollID string) (string, error) {

    response, err := s.client.Query(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "queryVotesByPoll", Args: [][]byte{[]byte(pollID)}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}

// query votes of a particular poll
func (s *SetupSDK) QueryVotesByVoterSDK(voterID string) (string, error) {

    response, err := s.client.Query(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "queryVotesByVoter", Args: [][]byte{[]byte(voterID)}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}