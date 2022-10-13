package interceptor

import (
	"context"
	"log"


	"github.com/Edbeer/client/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AccInterceptor struct {
	accClient   *services.AccClient
	accMethods  map[string]bool
	accessToken string
}

func NewAccInterceptor(
	accClient *services.AccClient,
	accMethods map[string]bool,
) (*AccInterceptor, error) {
	interceptor := &AccInterceptor{
		accClient:  accClient,
		accMethods: accMethods,
	}
	err := interceptor.refreshTokens()
	if err != nil {
		log.Fatal(err)
	}
	return interceptor, nil
}

func (interceptor *AccInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log.Printf("---> unary interceptor: %s", method)
		if interceptor.accMethods[method] {
			return invoker(
				interceptor.attachToken(ctx),
				method,
				req, reply,
				cc, opts...,
			)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (interceptor *AccInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		log.Printf("---> stream interceptor: %s", method)

		if interceptor.accMethods[method] {
			return streamer(
				interceptor.attachToken(ctx),
				desc,
				cc,
				method,
				opts...,
			)
		}
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (interceptor *AccInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.accessToken)
}

func (interceptor *AccInterceptor) refreshTokens() error {
	tokens, err := interceptor.accClient.SignIn()
	if err != nil {
		return err
	}
	interceptor.accessToken = tokens[0]
	log.Printf("token refresh: %v", tokens[0])
	return nil
}

