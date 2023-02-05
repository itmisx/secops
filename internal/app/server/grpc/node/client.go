package node

import (
	"log"
	"time"

	"free-access/internal/protoc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var (
	GRPCClient protoc.CommonClient
	NodeID     string
)

func init() {
	// NodeID = helper.RandString(32)
	NodeID = "1"
}

// CommonClientInit
func CommonClientInit(addr string) error {
	ln, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Second * 10,
			Timeout:             time.Second * 10,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
		return err
	}
	GRPCClient = protoc.NewCommonClient(ln)
	return nil
}

func StartGRPCClient(addr string) {
	CommonClientInit(addr)
	Heartbeat()
	NewShell()
}
