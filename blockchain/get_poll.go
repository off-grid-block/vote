package blockchain

import (
    "fmt"
    "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
    // "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)


func (s *SDKConfig) GetPollSDK(pollID string) (string, error) {

    response, err := s.client.Query(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "getPoll",  Args: [][]byte{[]byte(pollID)}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}

// read private details of vote using SDK
func (s *SDKConfig) GetPollPrivateDetailsSDK(pollID string) (string, error) {

    // create and send request for reading an entry
    response, err := s.client.Query(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "getPollPrivateDetails",  Args: [][]byte{[]byte(pollID)}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}

// retrieve all poll objects in state database
func (s *SDKConfig) QueryAllPollsSDK() (string, error) {

    response, err := s.client.Query(
        channel.Request{
            ChaincodeID: s.ChainCodeID, 
            Fcn: "queryAllPolls", 
            Args: [][]byte{}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}