package main

import (
	"context"
	"log"

	accountpb "github.com/Edbeer/proto/api/account/v1"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	transportOption := grpc.WithInsecure()
	cc1, err := grpc.Dial(":8080", transportOption)
	if err != nil {
		log.Fatal(err)
	}
	defer cc1.Close()

	client := accountpb.NewAccountServiceClient(cc1)
	res, err := client.SignIn(ctx, &accountpb.SignInRequest{
		Email: "edbeermtn@gmail.com",
		Password: "Password",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", res)
}