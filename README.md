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

### Launch using Docker

Start up the Fabric network:
1. ```cd fabric-samples/first-network``` (inside your fabric-samples repository)
2. ```./byfn.sh up -s couchdb```

Start up the DEON service API:
1. ```cd``` into your cloned version of this repository.
2. ```docker-compose up```
3. ```access the API at localhost:8000/api/v1/```

To stop the network and DEON service:
1. ```./byfn.sh down``` inside ```fabric-samples/first-network```
2. ```docker-compose down```

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

## Notes

Reference for running docker containers without docker-compose

Start IPFS container:
```
docker run -d --name ipfs-node \
  -v /tmp/ipfs-docker-staging:/export -v /tmp/ipfs-docker-data:/data/ipfs \
   -p 8080:8080 -p 4001:4001 -p 5001:5001 \
  ipfs/go-ipfs:latest
```

Build DEON service:
```
docker build -t vote_test:latest .
```

Run DEON service container:
```
docker run -it --rm \
  --name vote_test \
  --network="host" \
  --mount type=bind,source=/Users/brianli/deon/fabric-samples/first-network/channel-artifacts,target=/config/channel-artifacts \
  --mount type=bind,source=/Users/brianli/deon/fabric-samples/first-network/crypto-config,target=/config/crypto-config \
  vote_test
```

OR

```
docker run -it --rm \
  --add-host="localhost:10.0.0.69" \
  --name vote_test \
  -p 8000:8000 \
  --network="net_byfn" \
  --mount type=bind,source=/Users/brianli/deon/fabric-samples/first-network/channel-artifacts,target=/config/channel-artifacts,readonly \
  --mount type=bind,source=/Users/brianli/deon/fabric-samples/first-network/crypto-config,target=/config/crypto-config,readonly \
  vote_test
```