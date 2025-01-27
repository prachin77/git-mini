package files

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/prachin77/pkr/models"
)

// 1. both sender & receiver first have to setup their system
// 2. config folder will be created in send & receive directory
// 3. checks if exists or not
// 4. if exists than continue to write data or read from it
// 5. if doesn't exist than creates one
func CheckUserConfigFolderExists() bool {
	if _, err := os.Stat(models.ROOT_DIR); err == nil {
		fmt.Println("Config folder already exists")
		return true
	} else if os.IsNotExist(err) {
		fmt.Println("Config folder doesn't exist, creating it !")
		err := os.Mkdir(models.ROOT_DIR, 0777)
		if err != nil {
			fmt.Println("Failed to create config folder:", err)
			return false
		}
		fmt.Println("Config folder created successfully ...")
		return true
	} else {
		fmt.Println("Error checking config folder existence:", err)
		return false
	}
}

func UserConfigFileExists() bool {
	if _, err := os.Stat(models.USER_CONFIG_FILE); err == nil {
		fmt.Println("User configuration file already exists")
		return true
	} else if os.IsNotExist(err) {
		fmt.Println("User configuration file does not exist")
		return false
	} else {
		fmt.Println("Error checking user configuration file:", err)
		return false
	}
}

// 1. creates config file inside config folder
func CreateConfigFolder() bool {
	err := os.Mkdir(models.ROOT_DIR, 0777)
	if err != nil {
		fmt.Println("Failed to create config folder:", err)
		return false
	}
	return true
}

func CreateUserConfigFile() bool {
	file, err := os.Create(models.USER_CONFIG_FILE)
	if err != nil {
		fmt.Println("Failed to create user config file:", err)
		return false
	}
	defer file.Close()
	fmt.Println("user config file created successfully ...")
	return true
}

// 1. writes username, port, ip_address in it
func WriteInUserConfigFile(user_config_file_data *models.UserConfig) error {
	jsonBytes, err := json.MarshalIndent(user_config_file_data, "", "    ")
	if err != nil {
		fmt.Println("Unable to parse user data into JSON:", err)
		return err
	}

	err = os.WriteFile(models.USER_CONFIG_FILE, jsonBytes, 0777)
	if err != nil {
		fmt.Println("Unable to write into user config file:", err)
		return err
	}
	return nil
}

// 1. reads username, port, ip_address from userconfig.json file
func ReadFromUserConfigFile(user_config_file_data *models.UserConfig) (error) {
	jsonBytes, err := os.ReadFile(models.USER_CONFIG_FILE)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonBytes, user_config_file_data)
	if err != nil {
		return err
	}
	return nil
}
