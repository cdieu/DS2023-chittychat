# DS2023

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative Chitty-Chat/proto/proto.proto


Run the server by specifying a port: 
go run Chitty-Chat/server/server.go -port someportnumber


Run the a client by specifying a client name and the same port as server:
go run Chitty-Chat/client/client.go -user someclientname -port someportnumber
