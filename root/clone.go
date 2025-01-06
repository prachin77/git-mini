package root

import (
	"context"
	"fmt"
	"time"

	"github.com/prachin77/pkr/pb"
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

	

	fmt.Println("response : ",res)
}
