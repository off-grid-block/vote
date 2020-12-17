package blockchain

import (
	"fmt"
	"github.com/off-grid-block/fabric-sdk-go/pkg/client/channel"
    "github.com/off-grid-block/fabric-sdk-go/pkg/fabsdk"
	"time"
)


func UpdatePollStatusSDK(s *SetupSDK, pollID, status string) (string, error) {

    eventID := "updateEvent"

    clientContext := s.Fsdk.ChannelContext(s.ChannelID, fabsdk.WithUser("Voting"))
    client, err := s.CreateChannelClient(clientContext)
    if err != nil {
        return "", fmt.Errorf("failed to create new channel client: %v\n", err)
    }

    event, registered, notifier, err := s.CreateEventClient(clientContext, eventID)
    if err != nil {
        return "", fmt.Errorf("failed to create new event client: %v\n", err)
    }

    // unregister chaincode event
    defer event.Unregister(registered)

    // Create a request for poll update and send it
    response, err := client.Execute(channel.Request{ChaincodeID: "vote", Fcn: "updatePollStatus", Args: [][]byte{[]byte(pollID), []byte(status)}})
    if err != nil {
        return "", fmt.Errorf("failed to update: %v", err)
    }
    fmt.Println("Fabric transaction created")

    // Wait for the result of the submission
    select {
    case ccEvent := <-notifier:
        fmt.Printf("Received CC event: %v\n", ccEvent)
    case <-time.After(time.Second * 10):
        return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
    }

    return string(response.Payload), nil
}
