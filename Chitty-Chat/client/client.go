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
	serverPort = flag.Int("sPort", 5454, "server port number (should match the port used for the server)")
	timestamp  = lamport.LamportTime{Client: *user}
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
			if in.Time < timestamp.GetTimestamp() {
				in.Time = timestamp.GetTimestamp()
			}
			timestamp.Increment()
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
