package files

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/prachin77/pkr/models"
)

var (
	workspace_path        string
	workspace_hosted_date string
	workspace_hosted_port string
	username              string
)

func GetHostWorkspaceInfo(decrypted_workspace_password string, workspace_name string) (string, string, string, string) {
	file, err := os.Open(models.USER_CONFIG_FILE)
	if err != nil {
		fmt.Println("failed to open host user config file : ", err)
		return "", "", "" , ""
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&user_config)
	if err != nil {
		fmt.Println("error decoding json data from host user config file : ", err)
		return "", "" , "" , ""
	}

	for _, workspace := range user_config.SendWorkSpaces {
		if workspace.Workspace_Name == workspace_name && workspace.Workspace_Password == decrypted_workspace_password {
			username = user_config.Username
			workspace_hosted_port = user_config.Port
			workspace_path = workspace.Workspace_Path
			workspace_hosted_date = workspace.Workspace_Hosted_Date
			return workspace_path, workspace_hosted_date, workspace_hosted_port , username
		}
	}

	fmt.Println("no matching workspace found or invalid credentials")
	return "", "", "" , ""
}
