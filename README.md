# DS2023


Hi TA
Unfortunately we did not have time to finish the assignment for this time, and so we obviously do not live up to the requirements, hence the missing system requirements e.g broadcasting whenever a client joins and leaves the server etc. 
We also do not expect to pass this assignment, BUT we would like to receive some feedback on our current approach and what we have done so far, as well as what we can do to improve, and what we can do to get further in the system and technical requirements.
 We are sorry to waste your time :( 

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative Chitty-Chat/proto/proto.proto


Run the server by specifying a port: 
go run Chitty-Chat/server/server.go -port someportnumber


Run the a client by specifying a client name and the same port as server:
go run Chitty-Chat/client/client.go -user someclientname -port someportnumber
