package blockchain

import(
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"time"
	// "sync"
	// "log"
)

//var mlat [1000]time.Duration

// // Function to measure latency
// func (s *blockchainSDK.SetupSDK) GetLatency(start time.Time, name string, t int) {
//     elapsed := time.Since(start)
//     s.Latency[t] = elapsed

//     log.Printf("%s latency: %s", name, s.Latency[t])
// }

// add entry using SDK
func (s *SetupSDK) InitEntrySDK(PollID string, VoterID string, VoterSex string, VoterAge int, Salt string, VoteHash string) (string, error) {

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
func (s *SetupSDK) ReadEntrySDK(ID string) (string, error) {

	// create and send request for reading an entry
    response, err := s.client.Query(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "getVote",  Args: [][]byte{[]byte(ID)}})
    if err != nil {
            return "", fmt.Errorf("failed to query: %v", err)
    }

    return string(response.Payload), nil
}

//delete entry on chaincode using SDK
func (s *SetupSDK) DeleteEntrySDK(ID string) (string, error) {

	//register event
	eventID := "deleteEvent"
	reg, notifier, err := s.event.RegisterChaincodeEvent(s.ChainCodeID, eventID)
	if err != nil {
		return "", err
	}
	defer s.event.Unregister(reg)

	//create a request for deletion and sent it
	resp, err := s.client.Execute(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "deleteEntry", Args: [][]byte{[]byte(ID)} })
	if err != nil {
		return "", fmt.Errorf("failed to delete: %v",err)
	}

	// Wait for the result of the submission
        var ccEvent *fab.CCEvent
        select {
        case ccEvent = <-notifier:
                fmt.Printf("Received CC event: %v\n", ccEvent)
        case <-time.After(time.Second * 20):
                return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
        }

	return string(resp.Payload), nil
}

//search by username on chaincode using SDK
func (s *SetupSDK) SearchByOwnerSDK(Owner string) (string, error) {

	//creat and send request for reading an entry
        response, err := s.client.Query(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "searchByOwner",  Args: [][]byte{[]byte(Owner)}})
        if err != nil {
                return "", fmt.Errorf("failed to query: %v", err)
        }

        return string(response.Payload), nil
}

