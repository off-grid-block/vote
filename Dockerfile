### Stage 1: Build source code ###
FROM golang:1.13-alpine AS build

ENV API_URL 0.0.0.0:8002
ENV IPFS_ENDPOINT ipfs.node.example.com:5001
ENV CHANNEL_CONFIG /config/channel-artifacts/channel.tx
ENV CHAINCODE_GOPATH /
ENV CLIENT_AGENT_URL http://10.0.0.69:8013

# Set working directory
WORKDIR /src/
# Copy project to src directory
COPY . /src/

# Get & install packages before building
RUN go get
# Compile source code (dependencies from go.mod also installed)
RUN go build -o /bin/start
CMD ["/bin/start"]

# ### Stage 2: Move executable ###
# FROM scratch

# # Copy over compiled executable from previous stage
# COPY --from=build /bin/vote /vote
# # Specify command to start container
