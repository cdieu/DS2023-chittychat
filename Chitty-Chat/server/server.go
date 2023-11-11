package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"strconv"

	"DS2023-chittychat/Chitty-Chat/lamport"
	proto "DS2023-chittychat/Chitty-Chat/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var messages = []proto.SentMessage{}

type Server struct {
	proto.UnimplementedChatServer
	name    string
	port    int
	Clients []string
}

var servername = flag.String("name", "defaultserver", "server name")
var port = flag.Int("port", 0, "server port number")
var time int
var lam = lamport.LamportTime{Client: *servername}

func main() {
	flag.Parse()

	server := &Server{
		name: *servername,
		port: *port,
	}

	go startServer(server)

	for {

	}
}

func startServer(server *Server) {

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.port))

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Printf("Started server at port: %d\n", server.port)

	// Register the grpc server and serve its listener
	proto.RegisterChatServer(grpcServer, server)

	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

// Før havde vi func (server) sendandreceive (ctx, in proto.sentmessage) (proto.chat_sendandreceiveclient, error)
// Den sidste parentes betyder at vi returner dette. Så vi returnede proto.chat_sendandreceiveclient, samt en error.
// Men nu siden vi har stream to stream grpc, så returnerer vi ikke en ting kun baseret på et metode kald, men altid return whenever, så derfor returner vi nu kun en error.
func (s *Server) SendAndReceive(stream proto.Chat_SendAndReceiveServer) error {
	for {
		//This block receives the messages
		// create a receiver variable for the stream
		//This receiver variable represents a SentMessage from proto: clientName, Message, time
		receiver, err := stream.Recv()
		// If err is the exception of end of file exception, then return nil
		if err == io.EOF {
			return nil
		}

		if receiver.Time < lam.GetTimestamp() {
			receiver.Time = lam.GetTimestamp()
		}
		lam.Increment()
		time++
		// If an error is not null, then there exists an error, and we return it.
		if err != nil {
			return err
		}
		log.Printf("RECEIVING: Server received new message, making timestamp: %d\n", time)
		messages = append(messages, *receiver)

		//This block broadcasts the messages
		if receiver.Time < lam.GetTimestamp() {
			receiver.Time = lam.GetTimestamp()
		}
		lam.Increment()
		time++
		log.Printf("BROADCASTING: Client with name %s sent %s, making timestamp: %d\n", receiver.ClientName, receiver.Message, time)
		for _, msg := range messages {
			if err := stream.Send(&msg); err != nil {
				return err
			}
		}
	}
}

func (c *Server) Join(ctx context.Context, in *proto.JoinRequest) (*proto.JoinResponse, error) {
	if in.Time < lam.GetTimestamp() {
		in.Time = lam.GetTimestamp()
	}
	lam.Increment()
	time++

	messageToSend := "Participant " + in.ClientName + " joined Chitty-Chat at Lamport time " + strconv.FormatUint(uint64(in.Time), 10)
	messages = append(messages, proto.SentMessage{ClientName: "Server broadcast", Message: messageToSend})

	//log.Printf("BROADCASTING: Participant %s joined Chitty-Chat at Lamport time %d\n", in.ClientName, time)
	return &proto.JoinResponse{
		ServerName: c.name,
		Time:       lam.GetTimestamp(),
	}, nil
}

/*
func (s *Server) Leave(ctx context.Context, in *proto.LeftRequest, opts ...grpc.CallOption) (*proto.LeftResponse, error) {
	if in.Time < lam.GetTimestamp() {
		in.Time = lam.GetTimestamp()
	}
	lam.Increment()
	time++

	log.Printf("Client with name %s wants to leave, making timestamp: %d\n", in.ClientName, time)
	return &proto.LeftResponse{
		ServerName: s.name,
		Time:       lam.GetTimestamp(),
	}, nil
}*/

func connectToServer() (proto.ChatClient, error) {
	// Dial the server at the specified port.
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *port)
	} else {
		log.Printf("Connected to the server at port %d\n", *port)
	}
	return proto.NewChatClient(conn), nil
}
