package services

import (
	"context"
	"log"
	"time"

	examplepb "github.com/Edbeer/proto/api/example/v1"
	"google.golang.org/grpc"
)

type ExampleClient struct {
	service examplepb.ExampleServiceClient
}

func NewExampleClient(cc *grpc.ClientConn) *ExampleClient {
	service := examplepb.NewExampleServiceClient(cc)
	return &ExampleClient{service: service}
}

func (e *ExampleClient) Hello(req *examplepb.HelloRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 *time.Second)
	defer cancel()
	hello, err := e.service.Hello(ctx, req)
	if err != nil {
		log.Fatal("Bad hello")
	}
	log.Printf("Hello: %s", hello.Hello)
}

func (e *ExampleClient) World(req *examplepb.WorldRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 *time.Second)
	defer cancel()
	world, err := e.service.World(ctx, req)
	if err != nil {
		log.Fatal("Bad world")
	}
	log.Printf("World: %s", world.World)
}