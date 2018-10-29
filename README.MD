# IMG Manipulation

A full-stack application written in Go and React, where users can upload videos and the server splits it up to individual frames and converts the pictures to greyscale.

## Installation

1. Git clone this repo

### Server side
2. navigate to server directory and install necessary dependencies using
`go get -d ./...`

3. build server file by typing `go build server.go`

4. start server executable

5. by default, the server runs on localhost:8888

There is also a websocket communication between the server and client, which runs on localhost:5000

### Client Side

6. navigate to client directory

7. run `npm install`

8. run `npm start`

9. navigate to localhost:3000

## Usage

Files can be uploaded by either via the client 

Or files can also be curl-ed to the webserver

`
curl -X PUT http://localhost:8888/upload --upload-file example.mp4
`

Generated files are placed in the output folder of the server directory.
