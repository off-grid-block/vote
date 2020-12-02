
# DEON Service - Vote
This is the repository for the example voting application for the DEON platform. DEON Service Vote allows users to dynamically create polls and submit votes for user-defined polls.

## Start up the demo

### Prerequisites

Docker Desktop (2.2.0.0)

### Fabric Network
The DEON services are dependent on our modified version of Hyperledger Fabric. Use the DEON example Fabric network scripts located in the ```off-grid-block/off-grid-net``` repository. 

To clone the repository:
```git clone https://github.com/off-grid-block/off-grid-net.git```

To bring up the network, run:
```./cyfn.sh up -s couchdb```

These scripts will automatically pull the correct Docker images needed to bring up the example Fabric network. 

### VON Network (Indy)
The DEON services rely on VON Network, an implementation of a development level Indy Node network, developed by BCGov. For more information on the project and for additional instructions, see their [github repository](https://github.com/bcgov/von-network).

1. clone the repository: ```git clone https://github.com/bcgov/von-network.git```
2. Generate the Docker images: ```./manage build```
3. Start up the network: ```./manage start```

### DEON Services
The DEON services are dependent on a number of different components developed by the Yale Institute for Network Science. To bring up all necessary nodes:
1. Download  ```docker-compose-demo.yml``` from this repository.
2. Run ```export DOCKERHOST=`docker run --rm --net=host eclipse/che-ip````
3. Run ```docker-compose -f docker-compose-demo.yml up```

The docker-compose file will bring up:
 - the API exposing the endpoints services, hosted at http://localhost:8000/
-  the DEON vote service, using the code in this repository
 - the DEON core-service (github.com/off-grid-block/core-service)
 - a reverse proxy server to redirect requests to the correct component
 - the Admin Aries Cloud Agent
 - the Client Aries Cloud Agent
 - UIs to send instructions to & interact with both agents

## Test the demo

To test the demo, the first step is establishing a connection between the client and CI/MSP Aries Cloud Agents and creating a verifiable credential.
1. Create a controller for the client agent by sending a POST request to `localhost:8000/api/v1/admin/agent` with the following body:
    ```
    {
        "alias": "client",
        "agent_url": "http://client.example.com:8031",
    	"name": "Voting",
    	"secret": "kerapwd",
    	"type": "user"
    }
    ```
    This endpoint will create a signing DID & verkey pair for the application and store that information inside the client agent and VON Network ledger.
2. Using the ID returned in the previous POST request, send another POST request to `localhost:8000/api/v1/admin/agent/{client_agent_id}/connect` to establish a connection between the client and admin agents.
3. Register a public DID for the admin agent on the VON Network ledger by sending a POST request to `localhost:8000/api/v1/admin/agent/1/register-ledger`. (No need to repeat this process for the client agent; a public DID was automatically registered when you created the controller in step 1)
4. Issue a credeential by sending a POST to `http://localhost:8000/api/v1/admin/agent/{client_agent_id}/issue-credential` with the following body:
    ```
    {
        "app_name": "voter",
        "app_id": "101"
    }
    ```

### Create your first poll
Now we can initialize a poll! Send a post request to http://localhost:8000/api/v1/vote-app/poll with the following body: `{"PollID": "1", "Title": "My first poll", "Content": {"First choice": "DEON is good", "Second choice": "DEON is great", "Third choice": "DEON is amazing"}}`

### Summary

You can bring down the demo with the fo llowing commands:
1. ```./manage down``` inside ```von-network``` directory
2. ```./cyfn.sh down``` inside ```off-grid-net``` directory
3. ```docker-compose -f docker-compose-demo.yml down```
4. ```docker-compose -f docker-compose-demo.yml rm -f```
5. ```docker volume prune```


You've now created a poll and pushed it to the Fabric network. For more information on what else you can do with the vote service, check out the API documentation at https://app.swaggerhub.com/apis/haniavis/deon-core/0.3.0.