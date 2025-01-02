package root

import (
	"encoding/json"
	"fmt"
	"os"

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


// 1. create user config file in user's root folder(folder name = config)
// 2. setup username & password in user config file
func Setup(background_service_client pb.BackgroundServiceClient) {
	fmt.Printf("Enter username : ")
	fmt.Scan(&user_config.Username)

	fmt.Println("Hello "+user_config.Username+" setting up your system ...\n")
	fmt.Println("This may take few seconds , so be patient ...\n")

	configFolderExists := CheckUserConfigFolderExists()
	if configFolderExists{
		configFileCreated := CreateUserConfigFile()
		if configFileCreated{
			WriteInUserConfigFile(&user_config)
		}
	}
}

// 1. both sender & reciever first have to setup their shystum
// 2. config folder will be created in send & recieve directory
// 3. checks if exists or not 
// 4. if exists than continue to write data or read from it 
// if doesn't exists than create one 
func CheckUserConfigFolderExists() bool {
	if _ , err := os.Stat(ROOT_DIR); err == nil{
		fmt.Println("config folder exists")
		return true
	}else if os.IsNotExist(err){
		fmt.Println("config folder doesn't exists , creating it ...")
		err := os.Mkdir(ROOT_DIR , 0777)
		if err != nil{
			fmt.Println("failed to create config folder : ",err)
			return false
		}
		return true
	}else{
		fmt.Println("error checking config folder exists : ",err)
		return false
	}
}

// 1. creates config file inside config folder 
func CreateUserConfigFile() bool {
	file , err := os.Create(USER_CONFIG_FILE)
	if err != nil{
		fmt.Println("failed to create user config file : ",err)
		return false
	}
	defer file.Close()	
	return true
}

// 1. writes username , port , ip_address in it 
func WriteInUserConfigFile(user_config_file_data *models.UserConfig) error {
	jsonbytes , err := json.Marshal(user_config_file_data)
	if err != nil{
		fmt.Println("unable to parse name into json : ",err)
		return err
	}

	err = os.WriteFile(ROOT_DIR + "/userconfig.json",jsonbytes,0777)
	if err != nil{
		fmt.Println("unable to write into userconfig file 	: ",err)
		return err
	}
	return nil
}
