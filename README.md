# DEON Service - Vote
This is a repository for the voting application of the off-grid network.

## Setup

### Hyperledger Fabric

1. Clone the Hyperledger `fabric-samples` repository:
```git clone https://github.com/hyperledger/fabric-samples.git```
2. Inside `fabric-samples`, checkout to version 1.4.2:
```git checkout 1.4.2```
3. Replace the `first-network/byfn.sh` script with the script found [here](https://github.com/off-grid-block/off-grid-net/blob/master/cyfn.sh).
4. Download the Hyperledger Fabric v1.4.2 docker images:
```curl -sSL https://bit.ly/2ysbOFE | bash -s -- 1.4.2```

### IPFS

Install [IPFS](https://docs.ipfs.io/install/).

### Code variables

Navigate to the directory in which you want to clone this repository.
```
mkdir src && cd src
git clone https://github.com/off-grid-block/vote.git
```

In `main.go`, modify the following fields of the `fSetup` struct literal for your environment:

- `ChannelConfig` is the absolute path to `fabric-samples/first-network/channel-artifacts/channel.tx`.
- `ChaincodeGoPath` is the directory in which you created the `src` directory.

## Launch the DEON service API

1. To start up the Hyperledger Fabric test network, navigate to `fabric-samples/first-network`:
```./byfn.sh up -s couchdb```
2. In a separate shell session, start up the IPFS daemon:
```ipfs daemon```
3. Navigate to this repository on your machine. Run `go build` inside `chaincode` and in the top-level directory. 
4. Launch the API with `./vote`

To restart the network:
```./byfn.sh restart -s couchdb```
To take down the network:
```./byfn.sh down```


## Building additional apps

To add other apps: 

1. create a new package (use the voteapp package as a guide).
2. add an entry to the `ccPaths` map defined in `main()`. The key should be the chaincode ID, and the value should be the path to your app's chaincode folder.
3. add the following in `main()` under `fSetup.ChainCodeInstallationInstantiation("vote")`:
```
err = fSetup.ChainCodeInstallationInstantiation(yourchaincodeID)
	if err != nil {
		fmt.Printf("Failed to install and instantiate chaincode: %v\n", err)
		return
	}
```