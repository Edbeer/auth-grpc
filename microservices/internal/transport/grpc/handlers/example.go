package handlers

import (
	"context"
	"fmt"
	"io"

	examplepb "github.com/Edbeer/microservices/proto/api/example/v1"
)

type Example struct {
	examplepb.UnimplementedExampleServiceServer
}

func NewExample() *Example {
	return &Example{}
}

func (e *Example) Hello(ctx context.Context, req *examplepb.HelloRequest) (*examplepb.HelloResponse, error) {
	return &examplepb.HelloResponse{
		Hello: req.Hello,
	}, nil
}

func (e *Example) World(ctx context.Context, req *examplepb.WorldRequest) (*examplepb.WorldResponse, error) {
	return &examplepb.WorldResponse{
		World: req.World,
	}, nil
}

func (e *Example) StreamWorld(stream examplepb.ExampleService_StreamWorldServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Println("hello:", req.GetHello())
		
		if err := stream.Send(&examplepb.StreamWorldResponse{
			World: "World",
		}); err != nil {
			return err
		}
	}
}