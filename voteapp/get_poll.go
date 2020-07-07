package voteapp

import (
    "fmt"
    "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
    "github.com/off-grid-block/vote/blockchain"
    // "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)


func GetPollSDK(s *blockchain.SetupSDK, pollID string) (string, error) {

    response, err := s.Client.Query(channel.Request{ChaincodeID: "vote", Fcn: "getPoll",  Args: [][]byte{[]byte(pollID)}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}

// read private details of vote using SDK
func GetPollPrivateDetailsSDK(s *blockchain.SetupSDK, pollID string) (string, error) {

    // create and send request for reading an entry
    response, err := s.Client.Query(channel.Request{ChaincodeID: "vote", Fcn: "getPollPrivateDetails",  Args: [][]byte{[]byte(pollID)}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}

// retrieve all poll objects in state database
func QueryAllPollsSDK(s *blockchain.SetupSDK) (string, error) {

    response, err := s.Client.Query(
        channel.Request{
            ChaincodeID: "vote", 
            Fcn: "queryAllPolls", 
            Args: [][]byte{}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}