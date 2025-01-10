package root

import (
	"fmt"
	"time"

	"github.com/prachin77/pkr/models"
	pb "github.com/prachin77/pkr/pb"
	"github.com/prachin77/pkr/root/files"
	"github.com/prachin77/pkr/security"
)

var (
	user_config = models.UserConfig{}
)

// 1. create user config folder in user's root folder (folder name = config)
// 2. create userconfig.json file inside config folder
// 3. setup username, password, private & public key, port in user config file
func Setup(background_service_client pb.BackgroundServiceClient) {
	if files.UserConfigFileExists() {
		err := files.ReadFromUserConfigFile(&user_config)
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

	user_config.Port = models.Port

	fmt.Println("Hello " + user_config.Username + " setting up your system ...\n")
	fmt.Println("This may take a few seconds, so be patient ...")
	time.Sleep(1 * time.Second)

	configFolderExists := files.CheckUserConfigFolderExists()
	if !configFolderExists {
		fmt.Println("Creating configuration folder ...")
		if !files.CreateConfigFolder() {
			fmt.Println("Failed to create configuration folder. Exiting setup.")
			return
		}
	}

	fmt.Println("Creating user configuration file !")
	if files.CreateUserConfigFile() {
		files.WriteInUserConfigFile(&user_config)
	}

	private_key, public_key, err := security.GenerateRSAKeys()
	if private_key == nil || public_key == nil || err != nil {
		fmt.Println("error generating public & private keys:", err)
		return
	}

	if err = security.StorePrivateKeys(private_key , models.PRIVATE_KEY_FILE); err != nil {
		fmt.Println("error storing private keys in file:", err)
		return
	}

	if err = security.StorePublicKeys(public_key, models.PUBLIC_KEY_FILE); err != nil {
		fmt.Println("error storing public keys in file:", err)
		return
	}

	fmt.Println("~ created user : ", user_config.Username)
}
