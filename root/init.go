package root

import (
	"fmt"
	"os"

	"github.com/prachin77/pkr/models"
	"github.com/prachin77/pkr/pb"
)

var (
	sending_workspace = models.SendWorkSpaceFolder{}
)

// 1. register a folder / workspace to send / export
// 2. create log file outside config folder 
// 3. every user who's hosting a folder / workspace will have a log file 
func Init(background_service_client pb.BackgroundServiceClient) {
	fmt.Print("Enter Workspace Password : ")
	fmt.Scan(&sending_workspace.Workspace_Password)

	// get current working directory path 
	pwd , err := os.Getwd()
	if err != nil{
		fmt.Println("failed to retrive current working directory path : ",err)
		return
	}else{
		fmt.Println("current working directory path : ",pwd)
	}

	// get directory name 
	
}

// 1. checks if workspace is already Initialized or not 
// 2. if present then do nothing
// 3. if not then write into config file
func CheckWorkSpaceInUserConfigFile() bool {

	return false
}
