package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	// "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"time"
    "bytes"
    "encoding/gob"
)


// add entry using SDK
func (s *SetupSDK) InitVoteSDK(PollID string, VoterID string, VoterSex string, VoterAge string, VoteHash string) (string, error) {

    // Generate a random salt to concatenate with the vote's IPFS CID
    Salt := GenerateRandomSalt()

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
    registered, notifier, err := s.event.RegisterChaincodeEvent(s.ChainCodeID, eventID)
    if err != nil {
        return "Failed to register chaincode event", err
    }

    // unregister chaincode event
    defer s.event.Unregister(registered)

    // Create a request for vote init and send it
    response, err := s.client.Execute(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "initVote", Args: [][]byte{}, TransientMap: transientDataMap})
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

//read entry on chaincode using SDK
func (s *SetupSDK) GetVoteSDK(pollID, voterID string) (string, error) {

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
func (s *SetupSDK) GetVotePrivateDetailsSDK(pollID, voterID string) (string, error) {

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
func (s *SetupSDK) GetVotePrivateDetailsHashSDK(pollID, voterID string) (string, error) {
    
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

    return string(response.Payload), nil
}

// query the private details of a vote by poll
func (s *SetupSDK) QueryVotePrivateDetailsByPollSDK(pollID string) ([]string, error) {

    var cidList []string

    response, err := s.client.Query(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "queryVotePrivateDetailsByPoll",  Args: [][]byte{[]byte(pollID)}})
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
func (s *SetupSDK) QueryVotesByPollSDK(pollID string) (string, error) {

    response, err := s.client.Query(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "queryVotesByPoll", Args: [][]byte{[]byte(pollID)}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}

// query votes of a particular poll
func (s *SetupSDK) QueryVotesByVoterSDK(voterID string) (string, error) {

    response, err := s.client.Query(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "queryVotesByVoter", Args: [][]byte{[]byte(voterID)}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}

// add entry of poll using SDK
func (s *SetupSDK) InitPollSDK(PollID string, PollHash string) (string, error) {

    // Generate a random salt to concatenate with the poll's IPFS CID
    Salt := GenerateRandomSalt()

    text := fmt.Sprintf(
        "{\"PollID\":\"%s\",\"Salt\":\"%s\",\"PollHash\":\"%s\"}",
        PollID,
        Salt,
        PollHash,
    )

    eventID := "initEvent"

    // Add data to transient map (because we are using private data, all of the data will be in the transient map)
    transientDataMap := make(map[string][]byte)
    transientDataMap["poll"] = []byte(text)

    // register chaincode event
    registered, notifier, err := s.event.RegisterChaincodeEvent(s.ChainCodeID, eventID)
    if err != nil {
        return "Failed to register chaincode event", err
    }

    // unregister chaincode event
    defer s.event.Unregister(registered)

    // Create a request for vote init and send it
    response, err := s.client.Execute(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "initPoll", Args: [][]byte{}, TransientMap: transientDataMap})
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

func (s *SetupSDK) GetPollSDK(pollID string) (string, error) {

    response, err := s.client.Query(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "getPoll",  Args: [][]byte{[]byte(pollID)}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}

// read private details of vote using SDK
func (s *SetupSDK) GetPollPrivateDetailsSDK(pollID string) (string, error) {

    // create and send request for reading an entry
    response, err := s.client.Query(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "getPollPrivateDetails",  Args: [][]byte{[]byte(pollID)}})
    if err != nil {
        return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}

func (s *SetupSDK) UpdatePollStatusSDK(pollID, status string) (string, error) {

    eventID := "updateEvent"

    // register chaincode event
    registered, notifier, err := s.event.RegisterChaincodeEvent(s.ChainCodeID, eventID)
    if err != nil {
        return "Failed to register chaincode event", err
    }
    defer s.event.Unregister(registered)

    // Create a request for poll update and send it
    response, err := s.client.Execute(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "updatePollStatus", Args: [][]byte{[]byte(pollID), []byte(status)}})
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

// //delete entry on chaincode using SDK
// func (s *SetupSDK) DeleteEntrySDK(ID string) (string, error) {

// 	//register event
// 	eventID := "deleteEvent"
// 	reg, notifier, err := s.event.RegisterChaincodeEvent(s.ChainCodeID, eventID)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer s.event.Unregister(reg)

// 	//create a request for deletion and sent it
// 	resp, err := s.client.Execute(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "deleteEntry", Args: [][]byte{[]byte(ID)} })
// 	if err != nil {
// 		return "", fmt.Errorf("failed to delete: %v",err)
// 	}

// 	// Wait for the result of the submission
//         var ccEvent *fab.CCEvent
//         select {
//         case ccEvent = <-notifier:
//                 fmt.Printf("Received CC event: %v\n", ccEvent)
//         case <-time.After(time.Second * 20):
//                 return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
//         }

// 	return string(resp.Payload), nil
// }
