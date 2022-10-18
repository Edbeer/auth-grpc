package main

import (
	"log"

	"github.com/Edbeer/client/internal/interceptor"
	"github.com/Edbeer/client/internal/services"
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

func streamWorld(exampleClient *services.ExampleClient) {
	exampleClient.StreamWorld()
}

const (
	password = "password"
	email = "edbeermtn@gmail.com"
)

func main() {
	transportOption := grpc.WithInsecure()
	// cc1
	cc1, err := grpc.Dial(":8080", transportOption)
	if err != nil {
		log.Fatal(err)
	}
	defer cc1.Close()

	account := services.NewAccClient(cc1, password, email)
	interceptor, err := interceptor.NewAccInterceptor(account)
	if err != nil {
		log.Fatal(err)
	}
	interceptor.AuthMethods("/example.v1.ExampleService/Hello", true)
	interceptor.AuthMethods("/example.v1.ExampleService/World", true)
	interceptor.AuthMethods("/example.v1.ExampleService/StreamWorld", true)
	// cc2
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
	streamWorld(example)
	// account.SignOut(interceptor.RefreshToken)
}