package cmd

import (
	"fmt"
)

func ShowCommands() {
	fmt.Println("Available Commands:")
	fmt.Println("1. start server")
	fmt.Println("   Starts the server on this system, turning it into a node that can interact with other systems.\n")
	fmt.Println("2. movefile <file_path> <username>")
	fmt.Println("   Moves a file/folder to the specified user's directory.\n")
	fmt.Println("3. follow <username>")
	fmt.Println("   Follow a user, ensuring the username exists and preventing duplicates.\n")
	fmt.Println("4. unfollow <username>")
	fmt.Println("   Unfollow a user, ensuring the username exists and handling any errors.\n")
	fmt.Println("5. createpool <poolname>")
	fmt.Println("   Create a new pool with a unique name and optional description/tags.\n")
	fmt.Println("6. joinpool <poolname or pool_id>")
	fmt.Println("   Join an existing pool, checking membership and capacity limits.\n")
	fmt.Println("7. leavepool <poolname or pool_id>")
	fmt.Println("   Leave an existing pool and notify the admin if required.\n")
}