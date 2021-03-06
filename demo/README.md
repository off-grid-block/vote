
# DEON demo

This demo showcases the full stack of the DEON architecture using the Voting application as an example.

## Start up the demo

### Prerequisites

Docker Desktop (2.2.0.0)
Create a folder (e.g. deon/) to store all the code.

```mkdir deon && cd deon```

### Fabric Network

The DEON services are dependent on our modified version of Hyperledger Fabric. Use the DEON example Fabric network scripts located in the ```off-grid-block/off-grid-net``` repository. 

To clone the repository:
```git clone https://github.com/off-grid-block/off-grid-net.git && cd off-grid-net/```

To bring up the Fabric network, run:
```./cyfn.sh up -s couchdb -c mychannel```

These scripts will automatically pull the correct DEON Docker images needed to bring up the Fabric network. 

### Indy Ledger (VON Network)

The DEON services rely on an Indy ledger. For demonstration purposes we use the VON Network, an implementation of a development level Indy Node network, developed by BCGov. For more information on the project and for additional instructions, see their [github repository](https://github.com/bcgov/von-network).

1. clone the repository inside the parent deon/ folder: ```git clone https://github.com/bcgov/von-network.git && cd von-network/```
2. Generate the Docker images: ```./manage build```
3. Start up the network: ```./manage start```

### DEON Services

The DEON services are dependent on a number of different components developed by the Yale Institute for Network Science. To bring up all necessary nodes:
1. Clone this repository or just download the ```docker-compose-demo.yml``` file inside the parent deon/ folder.
2. Run ```export DOCKERHOST=`docker run --rm --net=host eclipse/che-ip` ```
3. Run ```docker-compose -f docker-compose-demo.yml up```

The docker-compose file will bring up:
 - the DEON API exposing the endpoints services, hosted at http://localhost:8000/
-  the DEON vote service, using the code in this repository
 - the DEON core-service (github.com/off-grid-block/core-service)
 - a reverse proxy server to redirect requests to the correct component
 - the DEON Admin Aries agent (API at http://localhost:8021)
 - an Aries agent as an example of a client's/application's agent (API at http://localhost:8031)

## Run the demo

To test the demo, the first step is establishing a connection between the Client and the Admin Aries agents and creating a verifiable credential. You can use the init script to automate the procedure

```sh init_ctls.sh```

or follow the instructions below to post requests to the DEON API.

1. Initiate a controller for the Admin agent by sending a POST request to `localhost:8000/api/v1/admin/agent` with the following body:
    ```
    {
        "agent_type": "admin"
    }
    ```
2. Initiate a controller for the Client agent by sending a POST request to `localhost:8000/api/v1/admin/agent` with the following body:
    ```
    {
        "agent_type": "client",
        "alias": "vote",
        "agent_url": "http://client.example.com:8031",
        "name": "Voting",
        "secret": "kerapwd",
        "type": "user"
    }
    ```
    This request will create a signing DID & verkey pair for the application and store that information inside the Client agent's wallet and the VON Network ledger.
3. Using the ID returned in the previous POST request, send another POST request to `localhost:8000/api/v1/admin/agent/{client_agent_id}/connect` to establish a connection between the Client and Admin agents.
4. Issue a DEON credential to the Client agent by sending a POST to `http://localhost:8000/api/v1/admin/agent/{client_agent_id}/issue-credential` with the following body:
    ```
    {
        "app_name": "voting",
        "app_id": "101"
    }
    ```

### Create your first poll

Now you can initialize a poll! Send a post request to http://localhost:8000/api/v1/vote-app/poll with the following body: `{"PollID": "1", "Title": "My first poll", "Content": {"First choice": "DEON is good", "Second choice": "DEON is great", "Third choice": "DEON is amazing"}}`

This request is a Fabric transaction signed by Fabric SDK through DID of the Client agent, verified by the Fabric through the Admin agent

### Summary

You can bring down the demo with the following commands:
1. ```docker-compose -f docker-compose-demo.yml down```
2. ```docker-compose -f docker-compose-demo.yml rm -f```
3. ```docker volume prune```
4. ```./manage down``` inside ```von-network``` directory
5. ```./cyfn.sh down``` inside ```off-grid-net``` directory

You've now created a poll and pushed it to the Fabric network. For more information on what else you can do with the DEON API, check out the documentation at https://app.swaggerhub.com/apis/haniavis/deon-core/0.6.0.
