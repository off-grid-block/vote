# Guide

### InitVoteSDK ###
Description:
Creates a new vote object, splitting up the object's metadata into two groups: the poll ID, voter ID, voter sex, and voter age are committed to the private database of collectionVote (which Org1 and Org2 can both access); the poll ID, voter ID, IPFS CID of the vote's contents, and a randomly generated salt are committed to the private database of collectionVotePrivateDetails (which only Org1 can access). Note: a vote cannot be created until the associated poll has been initialized with InitPollSDK.

Usage:
// fSetup is an instance of SetupSDK, defined in blockchain/setup.go  
resp, err := fSetup.InitVoteSDK("1234", "5678", "m", "21", "309fzlvvj5t4ofj3fvh")  
if err != nil {  
	fmt.Printf("error while initializing vote: %v\n", err)  
}  

### GetVoteSDK ###
Description  
Reads vote metadata (poll ID, voter ID, voter sex, voter age) stored in collectionVote.  

pollID: ID of the poll that the vote is cast for  
voterID: ID of the voter who cast the vote  

Usage:
resp, err := fSetup.GetVoteSDK("1234", "5678")
if err != nil {
	fmt.Printf("error while reading vote: %v\n", err)
}


### GetVotePrivateDetailsSDK ###
Description  
Reads a vote's IPFS CID and randomly generated salt, which are stored in collectionVotePrivateDetails.  

pollID: ID of the poll that the vote is cast for  
voterID: ID of the voter who cast the vote  

Usage:  
resp, err := fSetup.GetVotePrivateDetailsSDK("1234", "5678")  
if err != nil {  
	fmt.Printf("error while reading vote private details: %v\n", err)  
}  


### GetVotePrivateDetailsHashSDK ###
Description  
Returns the hash of the collectionVotePrivateDetails entry keyed by the passed poll ID and voter ID args. The hash is retrieved from the channel ledger.  

pollID: ID of the poll that the vote is cast for  
voterID: ID of the voter who cast the vote  

Usage:  
resp, err := fSetup.GetVotePrivateDetailsHashSDK("1234", "5678")  
if err != nil {  
	fmt.Printf("error while reading private details hash: %v\n", err)  
}  
// resp holds the private details hash of the vote associated with poll ID 1234 and voter ID 5678  


### QueryVotesByPollSDK ###
Description  
Returns the metadata of all the votes cast for poll with ID pollID.  

pollID: ID of poll of interest  

Usage:  
votes, err := fSetup.QueryVotesByPollSDK("1234")  
if err != nil {  
	fmt.Printf("error while querying all votes of poll: %v\n", pollID)  
	return  
}  
// votes is a []string, with each string a different vote entry from poll with ID pollID  


### QueryVotePrivateDetailsByPollSDK ###
Description:  
Returns the IPFS CIDs of all the votes cast for poll with ID pollID.  

pollID: ID of poll of interest  

Usage:  
cidList, err := fSetup.QueryVotePrivateDetailsByPollSDK("1234")  
if err != nil {  
	fmt.Printf("error while querying private details of votes from poll: %v\n", pollID)  
	return  
}  


### QueryVotesByVoterSDK ###
Description:  
Returns the metadata of all the votes cast by voter with ID voterID.  

voterID: ID of voter of interest  

Usage:  
votes, err := fSetup.QueryVotesByVoterSDK("5678")  
if err != nil {  
	fmt.Printf("error while query all votes from voter: %v\n", voterID)  
	return  
}  
// votes is a []string, with each string a different vote entry from poll with ID pollID  


### InitPollSDK ###
Description:  
Creates a new poll object, splitting the object's metadata into two groups: the poll ID, poll status, and number of votes cast are committed to the private database of collectionPoll (which Org1 and Org2 can both access); the poll ID, randomly generated salt, and IPFS CID of the poll's contents are committed to the private database of collectionPollPrivateDetails (which only Org1 can access).  

pollID: poll's ID  
pollHash: IPFS CID of poll's contents  

Usage:  
resp, err := fSetup.InitPollSDK("1234", "4930rfdlk5jds2oidfh")  
if err != nil {  
	fmt.Printf("error while initializing new poll: %v\n", err)  
	return  
}  


### GetPollSDK ###
Description:  
Reads poll metadata (poll ID, poll status, number of votes cast) stored in collectionPoll.  

pollID: ID of poll of interest  

Usage:  
resp, err := fSetup.GetPollSDK("1234")  
if err != nil {  
	fmt.Printf("error while reading poll: %v\n", err)  
	return  
}  


### GetPollPrivateDetailsSDK ###
Description:  
Reads a poll's IPFS CID and randomly generated salt, which are stored in collectionPollPrivateDetails.  

pollID: ID of poll of interest  

Usage:  
resp, err := fSetup.GetPollPrivateDetailsSDK("1234")  
if err != nil {  
	fmt.Printf("error while reading poll private details: %v\n", err)  
	return  
}  

### UpdatePollStatusSDK ###
Description:  
Updates the status of a poll. Possible statuses are: ongoing, paused, closed.  

pollID: ID of poll of interest  
status: the new status of the poll  

Usage:  
resp, err := fSetup.UpdatePollStatusSDK("1234", "closed")  
if err != nil {  
	fmt.Printf("error while updating poll status: %v\n", err)  
	return  
}  

// resp is nil if update poll transaction is successful  