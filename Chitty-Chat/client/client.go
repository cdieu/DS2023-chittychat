package main

import (
	"DS2023-chittychat/Chitty-Chat/lamport"
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
	user       = flag.String("user", "defaultUser", "client name")
	serverPort = flag.Int("port", 0, "server port number (should match the port used for the server)")
	timestamp  = lamport.LamportTime{Client: *user}
)

func main() {

	// Parse the flags to get the port for the client
	flag.Parse()
	//Create a connection variable, setting up connection to server
	conn, err := grpc.Dial(":"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))

	//Catch and handle error if it happens
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	//Create a client for the service Chat

	c := proto.NewChatClient(conn)

	//Join method
	c.Join(context.Background(), &proto.JoinRequest{
		ClientName: *user,
		Time:       timestamp.GetTimestamp(),
	})

	//For client to know when to type
	log.Println("Enter your input")

	//Create stream variable to begin sending and receiving. c is type ChatClient, a client of our grpc Chat
	stream, _ := c.SendAndReceive(context.Background())

	//Create a channel to continue the function until the reading is done.
	waitchannel := make(chan struct{})

	//Start a goroutine that runs the function
	go func() {
		for {
			//Create a variable "in" that contains the received message from our grpc SendAndReceive call
			in, err := stream.Recv()
			//Catches an End of File exception by closing the wait channel, signalling the function should stop
			if err == io.EOF {
				// read done.
				close(waitchannel)
				return
			}

			//Check for the highest current lamport timestamp
			//if in.Time < timestamp.GetTimestamp() {
				//in.Time = timestamp.GetTimestamp()
			//}
			//Increase timestamp
			timestamp.Increment()

			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}

			//Formatting to see which client sent which message in the client terminal.
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
		log.Println("> ", userInput)

		timestamp.Increment()

		if err := stream.Send(&proto.SentMessage{
			ClientName: *user,
			Message:    userInput,
			Time:       timestamp.GetTimestamp(),
		}); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}

	}
	stream.CloseSend()
	<-waitchannel
}
