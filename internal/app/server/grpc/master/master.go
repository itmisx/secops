package master

import (
	"log"
	"net"

	"free-access/internal/protoc"

	"google.golang.org/grpc"
)

type CommonServer struct {
	protoc.UnimplementedCommonServer
}

// StartGRPCService
func StartGRPCService(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	protoc.RegisterCommonServer(s, &CommonServer{})
	if err := s.Serve(ln); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
