package voteapp

import (
    "fmt"
    "github.com/off-grid-block/fabric-sdk-go/pkg/client/channel"
    "github.com/off-grid-block/fabric-sdk-go/pkg/fabsdk"
    // "github.com/off-grid-block/fabric-sdk-go/pkg/client/event"
    // "github.com/off-grid-block/fabric-sdk-go/pkg/common/providers/fab"
    "github.com/off-grid-block/vote/blockchain"
    "time"
)


// add entry of poll using SDK
func InitPollSDK(s *blockchain.SetupSDK, PollID, Title, PollHash string) (string, error) {

    // Generate a random salt to concatenate with the poll's IPFS CID
    Salt := blockchain.GenerateRandomSalt()

    text := fmt.Sprintf(
        "{\"PollID\":\"%s\",\"Title\":\"%s\",\"Salt\":\"%s\",\"PollHash\":\"%s\"}",
        PollID,
        Title,
        Salt,
        PollHash,
    )

    // Add data to transient map (because we are using private data, all of the data will be in the transient map)
    transientDataMap := make(map[string][]byte)
    transientDataMap["poll"] = []byte(text)

    clientContext := s.Fsdk.ChannelContext(s.ChannelID, fabsdk.WithUser("Voting"))

    client, err := s.CreateChannelClient(clientContext)
    if err != nil {
        return "", fmt.Errorf("failed to create new channel client: %v\n", err)
    }

    event, registered, notifier, err := s.CreateEventClient(clientContext)
    if err != nil {
        return "", fmt.Errorf("failed to create new event client: %v\n", err)
    }

    // unregister chaincode event
    defer event.Unregister(registered)

    // Create a request for vote init and send it
    response, err := client.Execute(channel.Request{ChaincodeID: "vote", Fcn: "initPoll", Args: [][]byte{}, TransientMap: transientDataMap})
    if err != nil {
        return "", fmt.Errorf("failed to initiate: %v", err)
    }

    // Wait for the result of the submission
    select {
    case ccEvent := <-notifier:
        fmt.Printf("Received CC event: %v\n", ccEvent)
    case <-time.After(time.Second * 10):
        return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", "initEvent")
    }

    return string(response.Payload), nil
}