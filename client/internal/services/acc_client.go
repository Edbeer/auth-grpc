package services

import (
	"context"
	"log"

	"time"

	accountpb "github.com/Edbeer/proto/api/account/v1"
	"google.golang.org/grpc"
)

type AccClient struct {
	service  accountpb.AccountServiceClient
	email    string
	password string
}

func NewAccClient(cc *grpc.ClientConn, password, email string) *AccClient {
	service := accountpb.NewAccountServiceClient(cc)
	return &AccClient{
		service:  service,
		password: password,
		email:    email,
	}
}

func (client *AccClient) SignIn() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	req := &accountpb.SignInRequest{
		Email: client.email,
		Password: client.password,
	}
	
	res, err := client.service.SignIn(ctx, req)
	if err != nil {
		return nil, err
	}
	tokens := []string{res.GetAccessToken(), res.GetRefreshToken()}
	return tokens, nil
}

func (client *AccClient) RefreshTokens(token string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	req := &accountpb.RefreshTokensRequest{
		RefreshToken: token,
	}

	res, err := client.service.RefreshTokens(ctx, req)
	if err != nil {
		return nil, err
	}
	tokens := []string{res.GetAccessToken(), res.GetRefreshToken()}
	return tokens, nil
}

func (client *AccClient) SignOut(token string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	req := &accountpb.SignOutRequest{
		RefreshToken: token,
	}

	res, err := client.service.SignOut(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("acc exist: ", res)
}