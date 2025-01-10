package roothandler

import (
	"context"
	"errors"
	"fmt"

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

func (s *BackgroundServiceServer) GetHostPcPublicKey(ctx context.Context, req *emptypb.Empty) (*pb.PublicKey, error) {
	public_key_data, host_public_key_file_path , err := utils.GetHostPublicKey()
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
		PublicKey: []byte(public_key_data),
		PublicKeyFilepath: host_public_key_file_path,
	}, nil
}

func (s *BackgroundServiceServer) InitWorkspaceConnWithPort(ctx context.Context , req *pb.InitRequest) (*pb.InitResponse , error) {
	// 1. decrypt workspace password using host PC public key , which is provided in request
	// 2. search workspace with the help of password & name
	// 3. return worskspace workspace path , port & username 

	decrypted_password , err := security.DecryptWorkspacePassword(req.WorkspacePassword)
	if err != nil{
		fmt.Println("error decrypting workspace password : ",err)
	}
	fmt.Println("decrypted workspace password : ",decrypted_password)

	workspace_path , workspace_hosted_date , workspace_hosted_port , username := files.GetHostWorkspaceInfo(decrypted_password , req.WorkspaceName)
	if workspace_path == "" || workspace_hosted_date == "" || workspace_hosted_port == ""{
		fmt.Println("error retrieving workspace path , hosted date or port from host PC !")
		return nil , errors.New("workspace not found , invalid credentials !")
	}

	return &pb.InitResponse{
		WorkspacePath: workspace_path,
		Port : workspace_hosted_port,
		Username: username,
		WorkspaceHostedDate: workspace_hosted_date,
	},nil
}
 