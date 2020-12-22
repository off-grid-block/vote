package blockchain

import (
    "fmt"
    "github.com/off-grid-block/fabric-sdk-go/pkg/client/channel"
    "github.com/off-grid-block/fabric-sdk-go/pkg/fabsdk"
)


//read entry on chaincode using SDK
func GetVoteSDK(s *SetupSDK, pollID, voterID string) (string, error) {

    // concatenate poll ID and voter ID to get vote key
    pollIdBytes := []byte(pollID)
    voterIdBytes := []byte(voterID)

    clientContext := s.Fsdk.ChannelContext(s.ChannelID, fabsdk.WithUser("Voting"))
    client, err := s.CreateChannelClient(clientContext)
    if err != nil {
        return "", fmt.Errorf("failed to create new channel client: %v\n", err)
    }

	// create and send request for reading an entry
    response, err := client.Query(channel.Request{ChaincodeID: "vote", Fcn: "getVote",  Args: [][]byte{pollIdBytes, voterIdBytes}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }
    fmt.Println("Fabric transaction created - get vote")

    return string(response.Payload), nil
}


// read private details of vote using SDK
func GetVotePrivateDetailsSDK(s *SetupSDK, pollID, voterID string) (string, error) {

    pollIdBytes := []byte(pollID)
    voterIdBytes := []byte(voterID)

    clientContext := s.Fsdk.ChannelContext(s.ChannelID, fabsdk.WithUser("Voting"))
    client, err := s.CreateChannelClient(clientContext)
    if err != nil {
        return "", fmt.Errorf("failed to create new channel client: %v\n", err)
    }

    // create and send request for reading an entry
    response, err := client.Query(
        channel.Request{
            ChaincodeID: "vote", 
            Fcn: "getVotePrivateDetails",  
            Args: [][]byte{pollIdBytes, voterIdBytes},
        },
    )
            
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }
    fmt.Println("Fabric transaction created - get vote private details")

    return string(response.Payload), nil
}

// get the private data hash of a vote
func GetVotePrivateDetailsHashSDK(s *SetupSDK, pollID, voterID string) (string, error) {
    
    pollIdBytes := []byte(pollID)
    voterIdBytes := []byte(voterID)

    clientContext := s.Fsdk.ChannelContext(s.ChannelID, fabsdk.WithUser("Voting"))
    client, err := s.CreateChannelClient(clientContext)
    if err != nil {
        return "", fmt.Errorf("failed to create new channel client: %v\n", err)
    }

    response, err := client.Query(
        channel.Request{
            ChaincodeID: "vote", 
            Fcn: "getVotePrivateDetailsHash", 
            Args: [][]byte{pollIdBytes, voterIdBytes},
        },
    )

    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }
    fmt.Println("Fabric transaction created - get vote private details hash")
    
    return string(response.Payload), nil
}
