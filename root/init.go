package root

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/prachin77/pkr/models"
	"github.com/prachin77/pkr/pb"
	"github.com/prachin77/pkr/root/files"
)

var (
	sending_workspace    = models.SendWorkSpaceFolder{}
	log_file_msg         string
)

// 1. register/host a folder/workspace to send into userconfig file
// 2. create log file inside config folder
// 3. every user who's hosting a folder / workspace will have a log file
func Init(background_service_client pb.BackgroundServiceClient) {
	fmt.Print("Enter Workspace Password: ")
	fmt.Scan(&sending_workspace.Workspace_Password)

	// Get current working directory path
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to retrieve current working directory path:", err)
		return
	}
	fmt.Println("Current working directory path:", pwd)

	sending_workspace.Workspace_Path = pwd
	sending_workspace.Workspace_Name = filepath.Base(pwd)

	// Check if the workspace already exists
	workspace_initialized, workspaces := files.CheckWorkSpaceInUserConfigFile(&sending_workspace)
	if !workspace_initialized {
		fmt.Println("Workspace already exists!")
		return
	}

	fmt.Println("Initializing workspace, wait a moment ...")
	files.InitalizeWorkspace(&sending_workspace)

	log_file_msg = "Workspace " + sending_workspace.Workspace_Name + " initialized on " + sending_workspace.Workspace_Hosted_Date + "\n"

	log_file_created := files.CreateLogFile()
	if log_file_created {
		if err := files.WriteInLogFile(log_file_msg); err == nil {
			fmt.Println("Successfully written in log file!")
		}
	}

	fmt.Println("Initialized workspaces:", workspaces)
}
