package main

import (
	"fmt"
	"os"
	"time"

	"github.com/prachin77/pkr/pb"
	"github.com/prachin77/pkr/root"
	"github.com/prachin77/pkr/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = ":8080"
)

func main() {
	var choice string

	utils.ClearScreen()

	conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("error connecting to server : ", err)
		return
	}
	fmt.Println("connection established on port : ", port)

	background_service_client := pb.NewBackgroundServiceClient(conn)
	fmt.Println(background_service_client)

	for {
		DisplayMenu()

		fmt.Printf("Enter choice : ")
		fmt.Scan(&choice)

		utils.ClearScreen()

		switch choice {
		case "S":
			root.Setup(background_service_client)
		case "I":
			root.Init(background_service_client)
		case "C":
			root.Clone(background_service_client)
		case "P":
			root.Push(background_service_client)
		case "Q":
			fmt.Println("Terminating ...")
			time.Sleep(2 * time.Second)
			os.Exit(0)
		default:
			fmt.Println("pls select a valid choice")
		}
	}
}

func DisplayMenu() {
	fmt.Println("--------------  WELCOME  --------------")
	fmt.Println("Enter S To Setup Service Configuration")
	fmt.Println("Enter P To Push Into Workspace")
	fmt.Println("Enter I To Initialize/Install Workspace")
	fmt.Println("Enter C To Clone Workspace")
	fmt.Println("Enter Q To Quit")
}
