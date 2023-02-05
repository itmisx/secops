package master

import (
	"io"
	"time"

	"free-access/internal/define"
	"free-access/internal/protoc"

	"github.com/itmisx/go-cache"
)

// NodeRegister
func (cs CommonServer) Heartbeat(ss protoc.Common_HeartbeatServer) error {
	for {
		select {
		case <-ss.Context().Done():
			return nil
		default:
			nodeInfo, err := ss.Recv()
			if err == io.EOF {
				return nil
			}
			if err == nil {
				// cache node info
				cache.HSet(define.NodeList, nodeInfo.NodeID, nodeInfo, time.Second*3, nil)
				time.Sleep(time.Millisecond * 100)
			}
		}
	}
}
