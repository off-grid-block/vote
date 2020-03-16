# Guide

### initVoteHandler ###
Description
Initialize & push votes on the Fabric network. Under the hood, the API adds the vote data to IPFS and pushes the IPFS Content Identifier (IPFS CID), along with some metadata, to the network.

Usage: See web/app.go


### initPollHandler ###
Description
Initialize & push poll data on the Fabric network. Under the hood, the API adds the poll data to IPFS and pushes the IPFS Content Identifier (IPFS CID), along with some metadata, to the network.

Usage: See web/app.go


### GetVoteHandler ###
Description
Retrieve public vote metadata from the Fabric network.

Usage: See web/app.go


### GetVotePrivateDetailsHandler ###
Description
Retrieve vote IPFS CID from the Fabric network private details database (collectionVotePrivateDetails) and use CID to retrieve vote data from IPFS.

Usage: See web/app.go


### GetVotePrivateDetailsHashHandler ###
Description
Retrieve has of vote private details from the Fabric network.

Usage: See web/app.go


### GetPollHandler ###
Description
Retrieve public poll metadata from the Fabric network.

Usage: See web/app.go


### GetPollPrivateDetailsHandler ###
Description
Retrieve poll IPFS CID from the Fabric network private details database (collectionPollPrivateDetails) and use CID to retrieve poll data from IPFS.

Usage: See web/app.go


### queryVotesByPollHandler ###
Description
Retrieve public metadata of all votes of a particular poll.

Usage: See web/app.go


### queryVotePrivateDetailsByPollHandler ###
Description
Retrieve private data of all votes of a particular poll.

Usage: See web/app.go


### updatePollStatusHandler ###
Description
Change the status of a poll. Choices are ongoing, paused, closed.

Usage: See web/app.go