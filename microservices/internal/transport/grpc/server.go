package grpc

import (
	"log"
	"net"
	"time"

	"github.com/Edbeer/microservices/cert"
	"github.com/Edbeer/microservices/internal/config"
	"github.com/Edbeer/microservices/internal/transport/grpc/interceptor"
	accountpb "github.com/Edbeer/microservices/proto/api/account/v1"
	examplepb "github.com/Edbeer/microservices/proto/api/example/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type Deps struct {
	Account     accountpb.AccountServiceServer
	Example     examplepb.ExampleServiceServer
	Interceptor *interceptor.AccountInterceptor
	Config      *config.Config
}

type GrpcServer struct {
	Deps Deps
	srv  *grpc.Server
}

func NewServer(deps Deps) *GrpcServer {
	// tlsCredentials
	tlsCredentials, err := cert.LoadTLSCredentialsServer()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}
	return &GrpcServer{
		srv: grpc.NewServer(
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle:     time.Duration(deps.Config.GrpsServer.MaxConnectionIdle) * time.Second,
				MaxConnectionAge:      time.Duration(deps.Config.GrpsServer.MaxConnectionAge) * time.Second,
				Time:                  time.Duration(deps.Config.GrpsServer.Time) * time.Second,
				Timeout:               time.Duration(deps.Config.GrpsServer.Timeout) * time.Second,
			}),
			grpc.Creds(tlsCredentials),
			grpc.UnaryInterceptor(deps.Interceptor.Unary()),
			grpc.StreamInterceptor(deps.Interceptor.Stream()),
		),
		Deps: deps,
	}
}

func (s *GrpcServer) ListenAndServe(port string) error {
	addr := ":" + port

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	// register services
	accountpb.RegisterAccountServiceServer(s.srv, s.Deps.Account)
	examplepb.RegisterExampleServiceServer(s.srv, s.Deps.Example)
	// reflection grpc server
	reflection.Register(s.srv)

	if err := s.srv.Serve(lis); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s *GrpcServer) Stop() {
	s.srv.GracefulStop()
}
