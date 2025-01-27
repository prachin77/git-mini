package roothandler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/prachin77/pkr/models"
	"github.com/prachin77/pkr/pb"
	"github.com/prachin77/pkr/root/files"
	"github.com/prachin77/pkr/security"
	"github.com/prachin77/pkr/utils"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
)

// type background_service_server struct {
// 	pb.UnimplementedBackgroundServiceServer
// }

type BackgroundServiceServer struct {
	pb.UnimplementedBackgroundServiceServer
}

const (
	Client_publicKey_Filepath = models.ROOT_DIR + "\\clientPublicKey.pem"
)

func (s *BackgroundServiceServer) GetHostPcPublicKey(ctx context.Context, req *emptypb.Empty) (*pb.PublicKey, error) {
	public_key_data, host_public_key_file_path, err := utils.GetHostPublicKey()
	if err != nil {
		fmt.Println("Error retrieving host public key:", err)
		return nil, err
	}
	if public_key_data == "" {
		fmt.Println("Public key data is empty.")
		return nil, fmt.Errorf("public key data is empty")
	}
	p, ok := peer.FromContext(ctx)
	if !ok {
		fmt.Println("Unable to retrieve client IP address.")
	} else {
		fmt.Println("Client IP address : ", p.Addr.String())
	}

	fmt.Println("public key retrived successfully ...")
	return &pb.PublicKey{
		PublicKey:         []byte(public_key_data),
		PublicKeyFilepath: host_public_key_file_path,
	}, nil
}

func (s *BackgroundServiceServer) InitWorkspaceConnWithPort(ctx context.Context, req *pb.InitRequest) (*pb.InitResponse, error) {
	// 1. decrypt workspace password using host PC public key , which is provided in request
	// 2. search workspace with the help of password & name
	// 3. return worskspace workspace path , port & username

	if req.WorkspaceName == "" || req.WorkspacePassword == "" || req.Port == "" || req.WorkspaceIp == "" {
		return nil, errors.New("missing required fields in the request")
	}

	// Decrypt workspace password using host PC private key
	decrypted_password, err := security.DecryptWorkspacePassword(req.WorkspacePassword)
	if err != nil {
		fmt.Println("error decrypting workspace password: ", err)
		return nil, errors.New("failed to decrypt workspace password")
	}
	fmt.Println("Decrypted workspace password: ", decrypted_password)

	workspace_path, workspace_name, workspace_hosted_date, workspace_hosted_port, username := files.GetHostWorkspaceInfo(decrypted_password, req.WorkspaceName)
	if workspace_path == "" || workspace_hosted_date == "" || workspace_hosted_port == "" {
		fmt.Println("error retrieving workspace path, hosted date, or port from host PC!")
		return nil, errors.New("workspace not found! Invalid credentials")
	}

	// store client public key to host PC
	// bcz when sending files from host PC to client
	// we can send encrypted data using client public key
	clientPublicKey, err := utils.ParseBytesToPublicKey(req.PublicKey)
	if err != nil {
		fmt.Println("error parsing client public key: ", err)
		return nil, errors.New("invalid public key provided")
	}

	err = security.StorePublicKeys(clientPublicKey, Client_publicKey_Filepath)
	if err != nil {
		fmt.Println("error saving client public key: ", err)
		return nil, errors.New("failed to store client public key")
	}
	fmt.Println("Client public key successfully stored at Host PC : ", Client_publicKey_Filepath)

	fmt.Printf("New connection established:\nClient IP: %s\nWorkspace: %s\n", req.WorkspaceIp, req.WorkspaceName)

	return &pb.InitResponse{
		WorkspacePath:       workspace_path,
		WorkspaceName:       workspace_name,
		Port:                workspace_hosted_port,
		Username:            username,
		WorkspaceHostedDate: workspace_hosted_date,
	}, nil
}

func (s *BackgroundServiceServer) GetFiles(req *pb.CloneRequest, stream pb.BackgroundService_GetFilesServer) error {
	// 1. Zip folder to be cloned except for folders & files -> config, exe files
	fmt.Println("Host PC workspace path:", req.WorkspacePath)
	fmt.Println("Host PC workspace name:", req.WorkspaceName)

	// Step 1: Create a zip file
	zip_filePath, err := files.ZipData(req.WorkspacePath, req.WorkspaceName)
	if err != nil {
		fmt.Println("Error zipping files inside folder:", err)
		return err
	}

	// Step 2: Generate AES key and nonce
	AES_KEY, err := security.GenerateAESKeys()
	if err != nil {
		fmt.Println("Error generating AES key:", err)
		return err
	}

	nonce, err := security.GenerateNonce()
	if err != nil {
		fmt.Println("Error generating nonce for AES key:", err)
		return err
	}

	// Step 3: Encrypt the zip file
	encrypted_zip_FilePath := strings.Replace(zip_filePath, ".zip", ".enc", 1)
	err = security.AESEncryptZipFile(zip_filePath, encrypted_zip_FilePath, AES_KEY, nonce)
	if err != nil {
		fmt.Println("Error encrypting zip file:", err)
		return err
	}

	// Step 4: Read client public key
	client_publicKey, err := os.ReadFile(Client_publicKey_Filepath)
	if err != nil {
		fmt.Println("Error accessing client public key from host PC:", err)
		return err
	}

	// Step 5: Encrypt AES key and nonce with client's public key
	encrypt_key, err := security.EncryptZipFile(string(AES_KEY), string(client_publicKey))
	if err != nil {
		fmt.Println("Error encrypting AES key:", err)
		return err
	}

	encrypt_nonce, err := security.EncryptZipFile(string(nonce), string(client_publicKey))
	if err != nil {
		fmt.Println("Error encrypting nonce:", err)
		return err
	}

	// Step 6: Open the encrypted zip file
	encrypted_zipFile_data, err := os.Open(encrypted_zip_FilePath)
	if err != nil {
		fmt.Println("Error opening encrypted zip file!")
		return err
	}
	defer encrypted_zipFile_data.Close()

	// Send metadata (key and nonce) to the client
	metadata := &pb.Files{
		FileName:    "metadata.enc",
		FileContent: []byte(fmt.Sprintf("Key:%s|Nonce:%s", encrypt_key, encrypt_nonce)),
	}
	err = stream.Send(metadata)
	if err != nil {
		fmt.Println("Error sending metadata:", err)
		return err
	}

	// Step 7: Send encrypted file in chunks
	buff := make([]byte, 2048)
	chunkNumber := 0
	totalBytesSent := 0
	const maxRetries = 3

	for {
		// Read a chunk from the file
		num, err := encrypted_zipFile_data.Read(buff)
		if err == io.EOF {
			fmt.Println("File transfer complete.")
			break
		}
		if err != nil {
			fmt.Printf("Error reading encrypted zip file at chunk %d: %v\n", chunkNumber, err)
			return err
		}

		chunkData := buff[:num]
		chunkNumber++
		totalBytesSent += num

		// Retry logic for sending each chunk
		retryCount := 0
		for {
			err = stream.Send(&pb.Files{
				FileName:    fmt.Sprintf("workspace_chunk_%d.enc", chunkNumber),
				FileContent: chunkData,
			})
			if err == nil {
				fmt.Printf("Successfully sent chunk %d (%d bytes).\n", chunkNumber, num)
				break
			}

			retryCount++
			if retryCount >= maxRetries {
				fmt.Printf("Failed to send chunk %d after %d retries. Error: %v\n", chunkNumber, retryCount, err)
				return err
			}

			fmt.Printf("Retrying chunk %d... Attempt %d of %d\n", chunkNumber, retryCount, maxRetries)
			time.Sleep(500 * time.Millisecond) // Add a short delay before retrying
		}
	}

	// Final status logs
	fmt.Printf("Total chunks sent: %d\n", chunkNumber)
	fmt.Printf("Total bytes sent: %d\n", totalBytesSent)
	fmt.Println("Encrypted zip file transfer completed.")

	return nil
}
