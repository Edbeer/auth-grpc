package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"guthub.com/Edbeer/microservices/internal/core"
	"guthub.com/Edbeer/microservices/pkg/jwt"
)

// Token Manager interface
type Manager interface {
	GenerateJWTToken(user *core.User) (string, error)
	Parse(accessToken string) (*jwt.Claims, error)
}

type AccountInterceptor struct {
	manager    Manager
	accessRole map[string][]string
}

func NewAccountInterceptor(manager Manager, accessRole map[string][]string) *AccountInterceptor {
	return &AccountInterceptor{
		manager:    manager,
		accessRole: accessRole,
	}
}

// Unary returns a server interceptor function to authenticate and
// authorize unary rpc
func (a *AccountInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("----> unary interseptor: ", info.FullMethod)

		err := a.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// Stream returns a server interceptor function to authenticate and
// authorize stream rpc
func (a *AccountInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Println("----> stream interseptor: ", info.FullMethod)

		err := a.authorize(ss.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, ss)
	}
}

func (a *AccountInterceptor) authorize(ctx context.Context, method string) error {
	accessRole, ok := a.accessRole[method]
	if !ok {
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := a.manager.Parse(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	for _, role := range accessRole {
		if role == claims.Role {
			return nil
		}
	}

	return status.Errorf(codes.PermissionDenied, "no permission to access this RPC")
}