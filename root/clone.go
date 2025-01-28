package root

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/prachin77/pkr/pb"
	"github.com/prachin77/pkr/root/files"
	"github.com/prachin77/pkr/security"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	workspace_ip string
	choice       string
)

func Clone(background_service_client pb.BackgroundServiceClient) {
	// 1. get host PC public key
	// 2. get client PC public key
	// 3. encrypt password
	// 4. init -> dont know what to do ?
	// 5. get files from host PC zip it & encrypt it
	// 6. unzip it & decrypt it
	// 7. store it in client folder

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	// for now enter manually , later automatically it'll be retrieved
	fmt.Print("Enter Workspace IP : ")
	fmt.Scan(&workspace_ip)

	fmt.Print("Enter Worspace Name : ")
	fmt.Scan(&sending_workspace.Workspace_Name)

	fmt.Print("Enter Workspace Password : ")
	fmt.Scan(&sending_workspace.Workspace_Password)

	res, err := background_service_client.GetHostPcPublicKey(ctx, &emptypb.Empty{})
	if err != nil {
		fmt.Println("error retrieving host PC public key : ", err)
		return
	}

	my_public_key_filepath := files.GetClientPublicKeyFilepath()

	my_public_key, err := os.ReadFile(my_public_key_filepath)
	if err != nil {
		fmt.Println("error retrieving client public key file path : ", err)
	}

	// encrypt workspace password using host PC public key
	// bcz when host will decrypt password it will have its public key
	encrypted_password, err := security.EncryptData(&sending_workspace, string(res.PublicKey))
	if err != nil {
		fmt.Println("Unable to encrypt workspace password:", err)
		return
	}

	init_res, err := background_service_client.InitWorkspaceConnWithPort(ctx, &pb.InitRequest{
		WorkspaceName:     sending_workspace.Workspace_Name,
		WorkspacePassword: encrypted_password,
		Port:              "8080",
		WorkspaceIp:       workspace_ip,
		PublicKey:         []byte(my_public_key),
	})
	if err != nil {
		fmt.Println("error establishing connection to workspace to be cloned : ", err)
		return
	}

	stream, err := background_service_client.GetFiles(ctx, &pb.CloneRequest{
		WorkspacePath: init_res.WorkspacePath,
		WorkspaceName: init_res.WorkspaceName,
		Port:          "8080",
	})
	if err != nil {
		fmt.Println("Error receiving files from server:", err)
		return
	}

	var data_bytes []byte
	var key_bytes []byte
	var nonce_bytes []byte

	for {
		data, err := stream.Recv()
		if err != nil {
			fmt.Println("error recieving data chunks")
		}

		if data.Filetype == 0 {
			data_bytes = append(data_bytes, data.FileContent...)
		} else if data.Filetype == 1 {
			key_bytes = append(key_bytes, data.FileContent...)
		} else if data.Filetype == 2 {
			nonce_bytes = append(nonce_bytes, data.FileContent...)
			break
		}
	}

	// for {
	// 	// Receive a chunk continously , eventually recieving whole zip file but in bytes
	// 	fileChunk, err := stream.Recv()
	// 	if err == io.EOF {
	// 		fmt.Println("File transfer completed.")
	// 		break
	// 	}
	// 	if err != nil {
	// 		fmt.Println("Error receiving chunk:", err)
	// 		return
	// 	}

	// 	fmt.Printf("Received chunk of type: %d (%d bytes)\n", fileChunk.Filetype, len(fileChunk.FileContent))
	// 	fmt.Printf("Chunk Content (as string): %s\n", string(fileChunk.FileContent))
	// }

	// decrypted_key , err := security.DecryptData(string(key_bytes))
	// if err != nil{
	// 	fmt.Println("error decrypting AES key : ",err)
	// }

	fmt.Println("host public key : ", string(res.PublicKey))
	fmt.Println("client public key : ", string(my_public_key))
	fmt.Println("client public key file path : ", my_public_key_filepath)
	fmt.Println("host public key file path : ", res.PublicKeyFilepath)
	fmt.Println("encrypted workspace password : ", encrypted_password)
	fmt.Println("\nresponse of connection to workspace to be cloned ... ")
	fmt.Println("workspace path : ", init_res.WorkspacePath)
	fmt.Println("workspace hosted date : ", init_res.WorkspaceHostedDate)
	fmt.Println("\nClone response ...")
	fmt.Println("zipped file data bytes : ", string(data_bytes))
	fmt.Println("zipped AES key bytes : ", string(key_bytes))
	fmt.Println("zipped file nonce bytes : ", string(nonce_bytes))
}
