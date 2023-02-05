package node

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"free-access/internal/protoc"

	"github.com/creack/pty"
	"google.golang.org/grpc/metadata"
)

// NewShell
func NewShell() {
	mdContext := metadata.AppendToOutgoingContext(context.Background(), "node_id", NodeID)
	ns, err := GRPCClient.NewShell(mdContext, &protoc.Empty{})
	if err != nil {
		fmt.Println(err)
	}
	for {
		sessionData, err := ns.Recv()
		if err == nil {
			shell(sessionData.SessionID)
		}
		time.Sleep(time.Millisecond * 100)
	}
}

// Shell
// use redis pub/sub resolve distributed server
func shell(sessionID string) {
	var ptyf *os.File
	go func() {
		// Create arbitrary command.
		c := exec.Command("bash")
		c.Env = []string{"TERM=xterm"}
		// Start the command with a pty.
		ptyf, _ = pty.Start(c)
		// Make sure to close the pty at the end.
		defer func() {
			_ = ptyf.Close()
		}()

		mdContext := metadata.AppendToOutgoingContext(context.Background(), "node_id", NodeID)
		conn, err := GRPCClient.ShellCMD(mdContext)
		if err != nil {
			fmt.Println(err)
		}
		// exec cmd
		go func() {
			for {
				send := make([]byte, 4096)
				ptyf.SetDeadline(time.Now().Add(time.Second * 1000))
				size, _ := ptyf.Read(send)
				if size > 0 {
					cmd := &protoc.ShellCMD{}
					cmd.SessionID = sessionID
					cmd.Data = send[0:size]
					conn.Send(cmd)
				}
				time.Sleep(time.Millisecond * 100)
			}
		}()
		// return exec res
		for {
			if cmd, err := conn.Recv(); err == nil {
				ptyf.Write(cmd.GetData())
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
}
