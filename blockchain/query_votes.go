package blockchain

import (
    "fmt"
    "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
    "github.com/off-grid-block/core-interface/pkg/sdk"
    // "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
    "bytes"
    "encoding/gob"
)


// query the private details of a vote by poll
func QueryVotePrivateDetailsByPollSDK(s *sdk.SDKConfig, pollID string) ([]string, error) {

    var cidList []string

    response, err := s.Client.Query(channel.Request{ChaincodeID: ccID, Fcn: "queryVotePrivateDetailsByPoll",  Args: [][]byte{[]byte(pollID)}})
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
func QueryVotesByPollSDK(s *sdk.SDKConfig, pollID string) (string, error) {

    response, err := s.Client.Query(channel.Request{ChaincodeID: ccID, Fcn: "queryVotesByPoll", Args: [][]byte{[]byte(pollID)}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}

// query votes of a particular poll
func QueryVotesByVoterSDK(s *sdk.SDKConfig, voterID string) (string, error) {

    response, err := s.Client.Query(channel.Request{ChaincodeID: ccID, Fcn: "queryVotesByVoter", Args: [][]byte{[]byte(voterID)}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}