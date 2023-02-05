package node

import (
	"context"
	"fmt"
	"time"

	"free-access/config"
	"free-access/internal/protoc"

	"google.golang.org/grpc/metadata"
)

func Heartbeat() {
	go func() {
		time.Sleep(time.Millisecond * 100)
		mdContext := metadata.AppendToOutgoingContext(context.Background(), "node_id", NodeID)
		hb, err := GRPCClient.Heartbeat(mdContext)
		if err != nil {
			fmt.Println(err)
		}
		for {
			cmd := &protoc.NodeInfo{}
			cmd.NodeID = NodeID
			for k, v := range config.Config.SSHService.Labels {
				label := &protoc.NodeInfo_NodeLabel{}
				label.Key = k
				label.Value = v
				cmd.Labels = append(cmd.Labels, label)
			}
			hb.Send(cmd)
			time.Sleep(time.Millisecond * 100)
		}
	}()
}
