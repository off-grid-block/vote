package voteapp

import (
	"fmt"
	"github.com/off-grid-block/fabric-sdk-go/pkg/client/channel"
    "github.com/off-grid-block/vote/blockchain"
	// "github.com/off-grid-block/fabric-sdk-go/pkg/common/providers/fab"
	"time"
)


func UpdatePollStatusSDK(s *blockchain.SetupSDK, pollID, status string) (string, error) {

    eventID := "updateEvent"

    // register chaincode event
    registered, notifier, err := s.Event.RegisterChaincodeEvent("vote", eventID)
    if err != nil {
        return "Failed to register chaincode event", err
    }
    defer s.Event.Unregister(registered)

    // Create a request for poll update and send it
    response, err := s.Client.Execute(channel.Request{ChaincodeID: "vote", Fcn: "updatePollStatus", Args: [][]byte{[]byte(pollID), []byte(status)}})
    if err != nil {
        return "", fmt.Errorf("failed to update: %v", err)
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