package root

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/prachin77/pkr/pb"
	"github.com/prachin77/pkr/security"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	workspace_ip string
)

func Clone(background_service_client pb.BackgroundServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	fmt.Print("Enter Workspace IP [Ip:Port] : ")
	fmt.Scan(&workspace_ip)

	fmt.Print("Enter Worspace Name : ")
	fmt.Scan(&sending_workspace.Workspace_Name)

	fmt.Print("Enter Workspace Password : ")
	fmt.Scan(&sending_workspace.Workspace_Password)

	res, err := background_service_client.GetHostPcPublicKey(ctx, &emptypb.Empty{})
	if err != nil{
		fmt.Println("error retrieving host PC public key : ",err)
		return
	}

	my_public_key_filepath := GetClientPublicKeyFilepath()
	
	my_public_key , err := os.ReadFile(my_public_key_filepath)
	if err != nil{
		fmt.Println("failed to read client public key !")
		return
	}

	encrypted_password , err := security.EncryptWorkspacePassword(&sending_workspace , string(res.PublicKey))
	if err != nil {
		fmt.Println("Unable to encrypt workspace password:", err)
		return
	}
	
	fmt.Println("host public key : ",res.PublicKey)
	fmt.Println("client public key : ",my_public_key)
	fmt.Println("client public key file path : ",my_public_key_filepath)
	fmt.Println("host public key file path : ",res.PublicKeyFilepath)
	fmt.Println("encrypted workspace password : ",encrypted_password)
}

func GetClientPublicKeyFilepath() (string) {
	my_public_key_filepath , err := filepath.Abs("./config/publickey.pem")
	if err != nil{
		fmt.Println("error retrieving client public key file path !")
		return "" 
	}else{
		return my_public_key_filepath 
	}
}
