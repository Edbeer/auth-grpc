package main

import (
	"log"

	"github.com/Edbeer/client/internal/services"
	"github.com/Edbeer/client/internal/interceptor"
	examplepb "github.com/Edbeer/proto/api/example/v1"
	"google.golang.org/grpc"
)

func hello(exampleClient *services.ExampleClient) {
	exampleClient.Hello(&examplepb.HelloRequest{
		Hello: "hello",
	})
}

func world(exampleClient *services.ExampleClient) {
	exampleClient.World(&examplepb.WorldRequest{
		World: "world",
	})
}

const (
	password = "password"
	email = "edbeermtn@gmail.com"
)

func authMethods() map[string]bool {
	const examplePath = "/example.v1.ExampleService/"
	return map[string]bool{
		examplePath + "Hello": true,
		examplePath + "World": true,
	}
}

func main() {
	transportOption := grpc.WithInsecure()
	cc1, err := grpc.Dial(":8080", transportOption)
	if err != nil {
		log.Fatal(err)
	}
	defer cc1.Close()

	account := services.NewAccClient(cc1, password, email)
	tokens, err := account.SignIn()
	if err != nil {
		log.Fatal(err)
	}
	interceptor, err := interceptor.NewAccInterceptor(account, authMethods(), tokens[0])
	if err != nil {
		log.Fatal(err)
	}
	t, err := account.RefreshTokens(tokens[1])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("tokens:", t)
	cc2, err := grpc.Dial(
		":8080", 
		transportOption,
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),

	)
	if err != nil {
		log.Fatal(err)
	}
	defer cc2.Close()

	example := services.NewExampleClient(cc2)
	hello(example)
	world(example)
}