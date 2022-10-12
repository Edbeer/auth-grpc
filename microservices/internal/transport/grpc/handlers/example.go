package handlers

import (
	"context"

	examplepb "github.com/Edbeer/proto/api/example/v1"
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