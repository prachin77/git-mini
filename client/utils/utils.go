package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"github.com/prachin77/pkr/client/cmd"
)

// ClearScreen clears the terminal screen based on the operating system.
func ClearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		fmt.Println("Unsupported platform.")
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// StartVirtualTerminal starts the virtual app terminal session.
func StartVirtualTerminal() {
	var cmdChoice string
	fmt.Println("Virtual App Terminal Started.... \nType 'exit' to return to the main menu . \nType 'clear' to clear screen . \nType 'H' for Help")

	for {
		fmt.Print(">> ")
		fmt.Scanln(&cmdChoice)

		if cmdChoice == "exit" {
			fmt.Println("Exiting virtual terminal...")
			return
		}

		if cmdChoice == "clear" {
			ClearScreen()
			continue
		}

		handleCommand(cmdChoice)
	}
}

// handleCommand processes commands in the virtual terminal.	
func handleCommand(cmdChoice string) {
	switch cmdChoice {
	case "H":
		cmd.ShowCommands()
	case "startserver":
		fmt.Println("Starting server...")
	case "movefile":
		fmt.Println("Move file command executed...")
	case "follow":
		fmt.Println("Following user...")
	case "unfollow":
		fmt.Println("Unfollowing user...")
	case "createpool":
		fmt.Println("Creating pool...")
	case "joinpool":
		fmt.Println("Joining pool...")
	case "leavepool":
		fmt.Println("Leaving pool...")
	default:
		fmt.Println("Unknown command. Type 'H' for help.")
	}
}
