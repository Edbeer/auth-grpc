package grpc

import (
	"log"
	"net"

	accountpb "github.com/Edbeer/proto/api/account/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Deps struct {
	Account accountpb.AccountServiceServer
}

type Server struct {
	Deps Deps
	srv  *grpc.Server
}

func NewServer(deps Deps) *Server {
	return &Server{
		srv:  grpc.NewServer(
			// grpc.UnaryInterceptor(deps.Interceptor.Unary()),
			// grpc.StreamInterceptor(deps.Interceptor.Stream()),
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
