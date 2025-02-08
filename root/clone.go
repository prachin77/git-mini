package root

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/prachin77/pkr/models"
	"github.com/prachin77/pkr/pb"
	"github.com/prachin77/pkr/root/files"
	"github.com/prachin77/pkr/security"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	workspace_ip string
)

func Clone(background_service_client pb.BackgroundServiceClient) {
	// # check ip of host & ip given by client are same or not

	// 1. get host PC public key
	// 2. get client PC public key
	// 3. encrypt password
	// 4. init -> dont know what to do ?
	// 5. get files from host PC , zip it & encrypt it
	// 6. write in host log file who sent for cloning it & client config file when cloned it
		// in client log file written about when data cloned 
		// in host log file its left -> do it later 
	// 7. unzip it & decrypt it
	// 8. store it in client folder

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	// for now enter manually , later automatically it'll be retrieved
	fmt.Print("Enter Workspace IP [if VPN is ON then enter VPN's IP] : ")
	fmt.Scan(&workspace_ip)

	IpRes, err := background_service_client.CheckIpAddress(ctx, &pb.IpRequest{
		IpAddress: workspace_ip,
	})
	if err != nil {
		fmt.Println("Error: Client & Host IP address mismatch or verification failed:", err)
		return
	}

	fmt.Println("CLient Ip address matched with Host Ip address ... : ", IpRes.Response)

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
		if err == io.EOF {
			fmt.Println("End of stream reached.")
			break
		}
		if err != nil {
			fmt.Println("Error receiving data chunk:", err)
			return
		}

		switch data.Filetype {
		case 0:
			data_bytes = append(data_bytes, data.FileContent...)
			fmt.Printf("Received data chunk of size %d bytes\n", len(data.FileContent))
		case 1:
			key_bytes = append(key_bytes, data.FileContent...)
			fmt.Println("Received key chunk")
		case 2:
			nonce_bytes = append(nonce_bytes, data.FileContent...)
			fmt.Println("Received nonce chunk")
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

	decrypted_key, err := security.DecryptData(string(key_bytes))
	if err != nil {
		fmt.Println("error decrypting AES key : ", err)
	}

	decrypted_nonce, err := security.DecryptData(string(nonce_bytes))
	if err != nil {
		fmt.Println("error decrypting nonce (AES) key : ", err)
	}

	// get decrypted zip file data from host in bytes
	data, err := security.AESDecryptZipFile(data_bytes, decrypted_key, decrypted_nonce)
	if err != nil {
		fmt.Println("error decrypting zip file with AES key : ", err)
	}

	ZipFilePath := init_res.WorkspaceName + "_zip.zip"
	if err = files.SaveDataToZip(data, ZipFilePath); err != nil {
		fmt.Println("error saving data from host to client zip file : ", err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("error retrieving current working directory path : ", err)
	}
	fmt.Println("current working directory path : ", pwd)
	if err := files.UnZipData(ZipFilePath, pwd); err != nil {
		fmt.Println("error unzipping file : ", err)
	}

	// remove zip file
	file, err := os.Open(ZipFilePath)
	if err != nil {
		fmt.Println("zip file not found : ", err)
	}
	defer os.Remove(ZipFilePath)
	file.Close()

	// save to user config file
	err = files.WriteRecivedWorkspaceInConfigFile(init_res.WorkspaceName, init_res.WorkspacePath, workspace_ip)
	if err != nil {
		fmt.Println("error saving received workspace info to config file : ", err)
	}

	// check if log file exists or not & write in it also
	log_file_msg := "workspace " + init_res.WorkspaceName + " cloned at " + time.Now().Format("2006-01-02 15:04:05") + "\n"
	if _, err := os.Stat(models.LOG_FILE); err != nil {
		fmt.Println("Log file not found, creating log file ...")
		LogFileCreated := files.CreateLogFile()
		if !LogFileCreated {
			fmt.Println("Error creating log file!")
			return
		}
		err = files.WriteInLogFile(log_file_msg)
		if err != nil {
			fmt.Println("Error writing in log file!")
		}
	} else {
		err := files.WriteInLogFile(log_file_msg)
		if err != nil {
			fmt.Println("Error writing in log file!")
		}
	}

}
