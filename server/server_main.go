// THIS IS BASE FILE OR FILE THAT ACTS LIKE SERVER WHICH SENDS DATA To CLIENT IN P2P ARCHITECHTURE
// SENDER -> SERVER
// 	   -> SENDS DATA

// CLIENT -> RECIEVER
// 	   -> RECIEVES DATA

package main

import (
	"fmt"
	"net"

	"github.com/prachin77/pkr/pb"
	roothandler "github.com/prachin77/pkr/server/root_handler"
	"github.com/prachin77/pkr/utils"
	"google.golang.org/grpc"
)

const (
	ip = "0.0.0.0" // public IP -> so that anyone (client) can access from anywhere
	port = ":8080"
)

// type background_service_server struct {
// 	pb.UnimplementedBackgroundServiceServer
// }

func main() {
	lis, err := net.Listen("tcp", ip + port)
	if err != nil {
		fmt.Println("error lisening to tcp server at port : ", lis.Addr())
	}
	grpcServer := grpc.NewServer(
		// for colored & structured logging into terminal
		grpc.UnaryInterceptor(utils.StructuredLoggerInterceptor()),
	)
	// pb.RegisterBackgroundServiceServer(grpcServer , &background_service_server{})

	background_service_server := &roothandler.BackgroundServiceServer{}
	pb.RegisterBackgroundServiceServer(grpcServer , background_service_server)

	fmt.Println("welcome to server ✌️✌️")
	fmt.Println("server running on port : ", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		fmt.Println("failed to start server at port : ", lis.Addr())
	}
}