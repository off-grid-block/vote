#!/bin/bash

# Initiate DEON Admin agent
curl -d '{ "agent_type": "admin" }' -H "Content-Type: application/json" -X POST http://localhost:8000/api/v1/admin/agent
sleep 2

# Initiate agent for the Voting app
VOTEAGENTID=$(curl -d '{ "agent_type": "client", "alias": "vote", "agent_url": "http://client.example.com:8031", "name": "Voting", "secret": "kerapwd", "type": "user" }' -H "Content-Type: application/json" -X POST http://localhost:8000/api/v1/admin/agent | cut -d ' ' -f4)
sleep 2

# Connect the two agents
curl -H "Content-Type: application/json" -X POST http://localhost:8000/api/v1/admin/agent/$VOTEAGENTID/connect
sleep 2

# Issue the DEON credential to the app agent
curl -H "Content-Type: application/json" -X POST http://localhost:8000/api/v1/admin/agent/$VOTEAGENTID/issue-credential -d '{ "app_name": "Vote", "app_id": "101" }'
