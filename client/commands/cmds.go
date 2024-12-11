package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/prachin77/pkr/client/utils"
)

func ShowCommands() {
	fmt.Println("Available Commands:")
	fmt.Println("NOTE : PLS ENTER ABSOLUTE PATH OF FILE/FOLDER FOR SMOOTH OPERATION\n\n")
	// fmt.Println("1. start server")
	// fmt.Println("   Starts the server on this system, turning it into a node that can interact with other systems.\n")
	fmt.Println("1. movefile <file_path> <username>")
	fmt.Println("   Moves a file/folder to the specified user's directory.\n")
	fmt.Println("2. follow <username>")
	fmt.Println("   Follow a user, ensuring the username exists and preventing duplicates.\n")
	fmt.Println("3. unfollow <username>")
	fmt.Println("   Unfollow a user, ensuring the username exists and handling any errors.\n")
	fmt.Println("4. createpool <poolname>")
	fmt.Println("   Create a new pool with a unique name and optional description/tags.\n")
	fmt.Println("5. joinpool <poolname or pool_id>")
	fmt.Println("   Join an existing pool, checking membership and capacity limits.\n")
	fmt.Println("6. leavepool <poolname or pool_id>")
	fmt.Println("   Leave an existing pool and notify the admin if required.\n")
	fmt.Println("7. getmyipadd")
	fmt.Println("   Gets your current IP address")
}

func HandleCommand(cmdChoice string) {
	parts := strings.Fields(cmdChoice)
	command := parts[0]
	args := parts[1:]

	switch command {
	case "H":
		ShowCommands()
	case "getmyipadd":
		utils.ShowClientIpAdd()
	case "movefile":
		if len(args) != 2 {
			// If the number of arguments is not correct, display usage info
			fmt.Println("Usage: movefile <file_path> <username>")
		} else {
			// If the correct number of arguments is passed, call MoveFile
			MoveFile(args[0], args[1])
		}
	default:
		fmt.Println("Unknown command. Type 'H' for help.")
	}
}

func MoveFile(filePath string, username string) {
	fmt.Println("file path:", filePath)
	fmt.Println("username:", username)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("file or folder doesn't exists in your system\n")
		fmt.Println("pls specify correct path")
	} else {
		fmt.Println("absolute file path : ", filePath)
	}
}
