package root

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/prachin77/pkr/models"
	pb "github.com/prachin77/pkr/pb"
)

const (
	ROOT_DIR         = "config"
	USER_CONFIG_FILE = ROOT_DIR + "\\userconfig.json"
	LOG_FILE         = ROOT_DIR + "\\logs.txt"
)

var (
	user_config = models.UserConfig{}
)

// 1. create user config folder in user's root folder (folder name = config)
// 2. create userconfig.json file inside config folder
// 3. setup username, password, private & public key, port in user config file
func Setup(background_service_client pb.BackgroundServiceClient) {
	if UserConfigFileExists() {
		err := ReadFromUserConfigFile(&user_config)
		if err == nil {
			fmt.Println("User is already set up! Welcome back : " + user_config.Username)
			return
		} else {
			fmt.Println("Error reading user configuration file:", err)
			return
		}
	}

	fmt.Printf("Enter username: ")
	fmt.Scan(&user_config.Username)

	fmt.Println("Hello " + user_config.Username + " setting up your system ...\n")
	fmt.Println("This may take a few seconds, so be patient ...")
	time.Sleep(1 * time.Second)

	configFolderExists := CheckUserConfigFolderExists()
	if !configFolderExists {
		fmt.Println("Creating configuration folder ...")
		if !CreateConfigFolder() {
			fmt.Println("Failed to create configuration folder. Exiting setup.")
			return
		}
	}

	fmt.Println("Creating user configuration file ...")
	if CreateUserConfigFile() {
		WriteInUserConfigFile(&user_config)
	}
	fmt.Println("~ created user: ", user_config.Username)
}

// 1. both sender & receiver first have to setup their system
// 2. config folder will be created in send & receive directory
// 3. checks if exists or not
// 4. if exists than continue to write data or read from it
// 5. if doesn't exist than creates one
func CheckUserConfigFolderExists() bool {
	if _, err := os.Stat(ROOT_DIR); err == nil {
		fmt.Println("Config folder already exists")
		return true
	} else if os.IsNotExist(err) {
		fmt.Println("Config folder doesn't exist, creating it ...")
		err := os.Mkdir(ROOT_DIR, 0777)
		if err != nil {
			fmt.Println("Failed to create config folder:", err)
			return false
		}
		fmt.Println("Config folder successfully created!")
		return true
	} else {
		fmt.Println("Error checking config folder existence:", err)
		return false
	}
}

func UserConfigFileExists() bool {
	if _, err := os.Stat(USER_CONFIG_FILE); err == nil {
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
	err := os.Mkdir(ROOT_DIR, 0777)
	if err != nil {
		fmt.Println("Failed to create config folder:", err)
		return false
	}
	return true
}

func CreateUserConfigFile() bool {
	file, err := os.Create(USER_CONFIG_FILE)
	if err != nil {
		fmt.Println("Failed to create user config file:", err)
		return false
	}
	defer file.Close()
	return true
}

// 1. writes username, port, ip_address in it
func WriteInUserConfigFile(user_config_file_data *models.UserConfig) error {
	jsonBytes, err := json.MarshalIndent(user_config_file_data, "", "    ")
	if err != nil {
		fmt.Println("Unable to parse user data into JSON:", err)
		return err
	}

	err = os.WriteFile(USER_CONFIG_FILE, jsonBytes, 0777)
	if err != nil {
		fmt.Println("Unable to write into user config file:", err)
		return err
	}
	return nil
}

// 1. reads username, port, ip_address from userconfig.json file
func ReadFromUserConfigFile(user_config_file_data *models.UserConfig) error {
	jsonBytes, err := os.ReadFile(USER_CONFIG_FILE)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonBytes, user_config_file_data)
	if err != nil {
		return err
	}
	return nil
}
