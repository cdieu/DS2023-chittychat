package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"net"
	"os"
	proto "simpleGuide/grpc"
	"strconv"
	"time"

	"github.com/cdieu/DS2023-chittychat.git/Chitty-Chat/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	proto.UnimplementedChatServer
	name 		string
	port		int
}

var port = flag.Int("port", 0, "server port number")

func main() {
	flag.Parse()

	server := &Server{
		name: "dima",
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

func (s *Server) SendAndReceive(ctx context.Context, in *proto.SentMessage) (*proto.Chat_SendAndReceiveClient, error) {
	log.Printf("Client with name %d sent message\n", in.clientName)
	return &proto.ReceivedMessage{
		Time: time.Now().String(),
		ServerName: s.name,
	}, nil
}
