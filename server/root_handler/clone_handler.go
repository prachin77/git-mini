package roothandler

import (
	"context"
	"fmt"

	"github.com/prachin77/pkr/pb"
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
 