package voteapp

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
    "github.com/off-grid-block/vote/blockchain"
	// "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"time"
)


// add entry using SDK
func InitVoteSDK(s *blockchain.SetupSDK, PollID string, VoterID string, VoterSex string, VoterAge string, VoteHash string) (string, error) {

    // Generate a random salt to concatenate with the vote's IPFS CID
    Salt := blockchain.GenerateRandomSalt()

    text := fmt.Sprintf(
        "{\"PollID\":\"%s\",\"VoterID\":\"%s\",\"VoterSex\":\"%s\",\"VoterAge\":%s,\"Salt\":\"%s\",\"VoteHash\":\"%s\"}",
        PollID,
        VoterID,
        VoterSex,
        VoterAge,
        Salt,
        VoteHash,
    )

    eventID := "initEvent"

    // Add data to transient map (because we are using private data, all of the data will be in the transient map)
	transientDataMap := make(map[string][]byte)
	transientDataMap["vote"] = []byte(text)

    // register chaincode event
    registered, notifier, err := s.Event.RegisterChaincodeEvent("vote", eventID)
    if err != nil {
        return "Failed to register chaincode event", err
    }

    // unregister chaincode event
    defer s.Event.Unregister(registered)

    // Create a request for vote init and send it
    response, err := s.Client.Execute(channel.Request{ChaincodeID: "vote", Fcn: "initVote", Args: [][]byte{}, TransientMap: transientDataMap})
    if err != nil {
        return "", fmt.Errorf("failed to initiate: %v", err)
    }

    // Wait for the result of the submission
    select {
    case ccEvent := <-notifier:
        fmt.Printf("Received CC event: %v\n", ccEvent)
    case <-time.After(time.Second * 10):
        return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
    }

    return string(response.Payload), nil
}
