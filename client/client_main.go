package main

import (
	"fmt"
	"os"
	"time"

	"github.com/prachin77/pkr/client/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const port = ":8080"

func main() {
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	for {
		displayMenu()
		var choice string
		fmt.Scanln(&choice)
		utils.ClearScreen()

		switch choice {
		case "H":
			fmt.Println("help ...")
		case "Q":
			fmt.Println("Terminating ...")
			time.Sleep(time.Second)
			os.Exit(0)
		default:
			fmt.Println("Please select a valid option ...")
		}
	}
}

func displayMenu() {
	fmt.Println("\nWelcome to Picker")
	fmt.Println("NOTE :  SAME PC CAN'T HAVE 2 HOSTS/CLIENTS\n")
	fmt.Println("Need to enter virtual enviorment to interact with application .")
	fmt.Println("1. Enter H for HELP")
	fmt.Println("2. Enter Q to QUIT")
	fmt.Print("Enter your choice : ")
}
