package main

import (
	"DS2023-chittychat/Chitty-Chat/proto"
	"bufio"
	"context"
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//var messages = []proto.SentMessage{}

var (
	//clientPort = flag.Int("cPort", 0, "client port number")
	serverPort = flag.Int("sPort", 5454, "server port number (should match the port used for the server)")
)

func main() {
	// Parse the flags to get the port for the client
	flag.Parse()
	conn, err := grpc.Dial(":"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	//Create a client
	c := proto.NewChatClient(conn)

	log.Println("Enter your input")

	stream, _ := c.SendAndReceive(context.Background())
	waitchannel := make(chan struct{})

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitchannel)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			log.Printf("[%s] %s\n", in.ClientName, in.Message)
		}
	}()

	for {
		// Get user input
		consoleReader := bufio.NewReader(os.Stdin)
		userInput, _ := consoleReader.ReadString('\n')
		userInput = strings.Split(userInput, "\n")[0]

		// Close connection
		if userInput == "quit" {
			break
		}

		// Print the user input
		//log.Println("> ", userInput)

		if err := stream.Send(&proto.SentMessage{
			ClientName: "christine",
			Message:    userInput,
		}); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}

	}
	stream.CloseSend()
	<-waitchannel
}
