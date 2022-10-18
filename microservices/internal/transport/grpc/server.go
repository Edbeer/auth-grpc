package grpc

import (
	"log"
	"net"

	accountpb "github.com/Edbeer/microservices/proto/api/account/v1"
	examplepb "github.com/Edbeer/microservices/proto/api/example/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/Edbeer/microservices/internal/transport/grpc/interceptor"
)

type Deps struct {
	Account accountpb.AccountServiceServer
	Example examplepb.ExampleServiceServer
	Interceptor *interceptor.AccountInterceptor
}

type Server struct {
	Deps Deps
	srv  *grpc.Server
}

func NewServer(deps Deps) *Server {
	return &Server{
		srv:  grpc.NewServer(
			grpc.UnaryInterceptor(deps.Interceptor.Unary()),
			grpc.StreamInterceptor(deps.Interceptor.Stream()),
		),
		Deps: deps,
	}
}

func (s *Server) ListenAndServe(port string) error {
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

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
