package blockchain

import (
    "fmt"
    "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
    // "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
    "reflect"
)


//read entry on chaincode using SDK
func (s *SDKConfig) GetVoteSDK(pollID, voterID string) (string, error) {

    // concatenate poll ID and voter ID to get vote key
    pollIdBytes := []byte(pollID)
    voterIdBytes := []byte(voterID)

	// create and send request for reading an entry
    response, err := s.client.Query(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "getVote",  Args: [][]byte{pollIdBytes, voterIdBytes}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}


// read private details of vote using SDK
func (s *SDKConfig) GetVotePrivateDetailsSDK(pollID, voterID string) (string, error) {

    pollIdBytes := []byte(pollID)
    voterIdBytes := []byte(voterID)

    // create and send request for reading an entry
    response, err := s.client.Query(
        channel.Request{
            ChaincodeID: s.ChainCodeID, 
            Fcn: "getVotePrivateDetails",  
            Args: [][]byte{pollIdBytes, voterIdBytes}})
            
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}

// get the private data hash of a vote
func (s *SDKConfig) GetVotePrivateDetailsHashSDK(pollID, voterID string) (string, error) {
    
    pollIdBytes := []byte(pollID)
    voterIdBytes := []byte(voterID)

    response, err := s.client.Query(
        channel.Request{
            ChaincodeID: s.ChainCodeID, 
            Fcn: "getVotePrivateDetailsHash", 
            Args: [][]byte{pollIdBytes, voterIdBytes}})

    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    fmt.Println(reflect.TypeOf(response.Payload))
    fmt.Println(response.Payload)
    fmt.Println(string(response.Payload))
    
    return string(response.Payload), nil
}