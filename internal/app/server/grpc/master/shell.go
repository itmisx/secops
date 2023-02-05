package master

import (
	"io"
	"sync"
	"time"

	"free-access/internal/protoc"

	"google.golang.org/grpc/metadata"
)

var (
	// map[sessionID]chan byte
	ShellCMDSendList sync.Map
	// map[sessionID]chan []byte
	ShellCMDRecvList sync.Map
	// map[nodeID]chan sessionID
	ShellSessionList sync.Map
)

// NewShell
// use redis pub/sub resolve distributed server
func (cs CommonServer) NewShell(in *protoc.Empty, ss protoc.Common_NewShellServer) error {
	nodeID := "1"
	md, _ := metadata.FromIncomingContext(ss.Context())
	if len(md["node_id"]) > 0 {
		nodeID = md["node_id"][0]
	}
	if nodeID == "" {
		return nil
	}
	ShellSessionList.LoadOrStore(nodeID, make(chan string, 100))
	v, _ := ShellSessionList.Load(nodeID)
	sessionIDChan := v.(chan string)
	for {
		select {
		case <-ss.Context().Done():
			return nil
		case sessionID := <-sessionIDChan:
			shellSession := &protoc.ShellSession{}
			shellSession.SessionID = sessionID
			ss.Send(shellSession)
		}
	}
}

// ShellCMD
// use redis pub/sub resolve distributed server
func (cs CommonServer) ShellCMD(ss protoc.Common_ShellCMDServer) error {
	sessionID := ""
	// recv
	go func() {
		for {
			select {
			case <-ss.Context().Done():
				return
			default:
				cmd, err := ss.Recv()
				if err == io.EOF {
					return
				}
				if err == nil {
					sessionID = cmd.SessionID
					v, loaded := ShellCMDRecvList.LoadOrStore(sessionID, make(chan []byte, 100))
					if loaded {
						if cmdChan, ok := v.(chan []byte); ok {
							cmdChan <- cmd.GetData()
						}
					}
				}
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()
	// send
	for {
		select {
		case <-ss.Context().Done():
			return nil
		default:
			if sessionID == "" {
				continue
			}
			v, loaded := ShellCMDSendList.LoadOrStore(sessionID, make(chan byte, 100))
			if loaded {
				if cmdChan, ok := v.(chan byte); ok {
					cmd := &protoc.ShellCMD{}
					for {
						if len(cmdChan) > 0 {
							cmd.Data = append(cmd.Data, <-cmdChan)
						} else {
							break
						}
					}
					if len(cmd.GetData()) > 0 {
						ss.Send(cmd)
					}
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
}
