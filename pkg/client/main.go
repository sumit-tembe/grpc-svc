package main

import (
	"context"
	"fmt"

	pb "github.com/sumit-tembe/grpc-svc/pkg/grpc/user"
	"google.golang.org/grpc"
)

type usersClient struct {
	Client pb.UsersClient
}

func initUsersClient(url string) usersClient {
	cc, err := grpc.Dial(url, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return usersClient{
		Client: pb.NewUsersClient(cc),
	}
}

func (c *usersClient) GetUsers(ctx context.Context, ids []int64) (*pb.GetUsersResponse, error) {
	req := &pb.GetUsersRequest{
		Ids: ids,
	}
	return c.Client.GetUsers(ctx, req)
}

func main() {
	client := initUsersClient("localhost:8080")
	fmt.Println(client.GetUsers(context.Background(), nil))
}
