package blockchain

import (
    "fmt"
    "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
    "github.com/off-grid-block/core-interface/pkg/sdk"
    // "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
    "reflect"
)


//read entry on chaincode using SDK
func GetVoteSDK(s *sdk.SDKConfig, pollID, voterID string) (string, error) {

    // concatenate poll ID and voter ID to get vote key
    pollIdBytes := []byte(pollID)
    voterIdBytes := []byte(voterID)

	// create and send request for reading an entry
    response, err := s.Client.Query(channel.Request{ChaincodeID: ccID, Fcn: "getVote",  Args: [][]byte{pollIdBytes, voterIdBytes}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}


// read private details of vote using SDK
func GetVotePrivateDetailsSDK(s *sdk.SDKConfig, pollID, voterID string) (string, error) {

    pollIdBytes := []byte(pollID)
    voterIdBytes := []byte(voterID)

    // create and send request for reading an entry
    response, err := s.Client.Query(
        channel.Request{
            ChaincodeID: ccID, 
            Fcn: "getVotePrivateDetails",  
            Args: [][]byte{pollIdBytes, voterIdBytes}})
            
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}

// get the private data hash of a vote
func GetVotePrivateDetailsHashSDK(s *sdk.SDKConfig, pollID, voterID string) (string, error) {
    
    pollIdBytes := []byte(pollID)
    voterIdBytes := []byte(voterID)

    response, err := s.Client.Query(
        channel.Request{
            ChaincodeID: ccID, 
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