package services

import (
	"context"
	"io"
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

func (e *ExampleClient) StreamWorld() {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()

	stream, err := e.service.StreamWorld(ctx)
	if err != nil {
		log.Fatal("Stream", err)
	}
	for i := 0; i < 3; i++ {
		if err := stream.Send(&examplepb.StreamWorldRequest{
			Hello: "Hello",
		}); err != nil {
			log.Fatal("send", err)
		}
		log.Println("Hello")
	}

	if err := stream.CloseSend(); err != nil {
		log.Fatal("close stream", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("recv", err)
		}
		log.Printf("%s\n", res)
	}
}