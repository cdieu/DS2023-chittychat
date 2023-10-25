package main

import (
	"flag"
)

type Client struct {
	name       string
	portNumber int
}

var (
	clientPort = flag.Int("cPort", 0, "client port number")
	serverPort = flag.Int("sPort", 0, "server port number (should match the port used for the server)")
)

func main() {
	// Parse the flags to get the port for the client
	flag.Parse()

	// Create a client
	client := &Client{
		name:       "christine",
		portNumber: *clientPort,
	}

	// Wait for the client (user) to ask for the time
	go waitForMessage(client)

	for {

	}
}
