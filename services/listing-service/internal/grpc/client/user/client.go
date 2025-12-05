package user

import (
	"context"
	"fmt"
	userpb "github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/pb/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api  userpb.UserServiceClient
	conn *grpc.ClientConn
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	return &Client{
		api:  userpb.NewUserServiceClient(conn),
		conn: conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetUserByID(ctx context.Context, userID string) (*userpb.User, error) {
	resp, err := c.api.GetUserByID(ctx, &userpb.GetUserByIDRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return resp.User, nil
}

func (c *Client) ValidateUser(ctx context.Context, userID string) (bool, error) {
	resp, err := c.api.ValidateUser(ctx, &userpb.ValidateUserRequest{
		UserId: userID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to validate user: %w", err)
	}
	return resp.Exists && resp.IsActive, nil
}
