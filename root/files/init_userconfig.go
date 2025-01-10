package files

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/prachin77/pkr/models"
)

var (
	user_config = models.UserConfig{}
)

// 1. checks if workspace is already Initialized or not -> present in send workspace slice in userconfig file
func CheckWorkSpaceInUserConfigFile(sending_workspace *models.SendWorkSpaceFolder) (bool, []models.SendWorkSpaceFolder) {
	data, err := os.ReadFile(models.USER_CONFIG_FILE)
	if err != nil {
		fmt.Println("error reading from user config file:", err)
		return false, nil
	}

	if err := json.Unmarshal(data, &user_config); err != nil {
		fmt.Println("error unmarshalling data:", err)
		return false, nil
	}

	// Check if SendWorkSpaces slice is empty
	if len(user_config.SendWorkSpaces) == 0 {
		fmt.Println("No workspaces initialized.")
		return false, nil
	}

	for _, workspace := range user_config.SendWorkSpaces {
		fmt.Printf("Existing workspace : %+v\n", workspace)
	}

	return true, user_config.SendWorkSpaces
}

func InitalizeWorkspace(sending_workspace *models.SendWorkSpaceFolder) {
	sending_workspace.Workspace_Hosted_Date = time.Now().Format("2006-01-02")

	data, err := os.ReadFile(models.USER_CONFIG_FILE)
	if err != nil {
		fmt.Println("error reading user config file : ", err)
		return
	}
	if err := json.Unmarshal(data, &user_config); err != nil {
		fmt.Println("error unmarshalling data:", err)
		return
	}

	// Append the new workspace to the SendWorkSpaces slice
	user_config.SendWorkSpaces = append(user_config.SendWorkSpaces, *sending_workspace)

	updatedData, err := json.MarshalIndent(user_config, "", "    ")
	if err != nil {
		fmt.Println("error marshalling updated user config data : ", err)
		return
	}

	err = os.WriteFile(models.USER_CONFIG_FILE, updatedData, os.ModePerm)
	if err != nil {
		fmt.Println("error writing updated data into user config file : ", err)
		return
	}
	fmt.Println("Workspace initialized and added to the configuration successfully!")
}

func CreateLogFile() bool {
	file, err := os.Create(models.LOG_FILE)
	if err != nil {
		fmt.Println("failed to create log file : ")
		return false
	}
	defer file.Close()
	fmt.Println("log file successfully created ")
	return true
}

func WriteInLogFile(log_file_msg string) error {
	file, err := os.OpenFile(models.LOG_FILE, os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println("error opening log file , issue in path : ", err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(log_file_msg)
	if err != nil {
		fmt.Println("error writing in log file : ", err)
		return err
	}
	return nil
}
