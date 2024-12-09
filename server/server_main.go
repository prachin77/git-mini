package main

import (
	"fmt"
	"net"

	pb "github.com/prachin77/pkr/proto"
	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

type authServer struct {
	pb.AuthServiceServer
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error listening to tcp server at port : ", lis.Addr())
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &authServer{})

	fmt.Println("welcome to server ✌️✌️")
	
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Println("Failed to server at port : ", lis.Addr())
	}
	fmt.Println("server started at port : ", lis.Addr())
}
