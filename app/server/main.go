package main

import (
	"free-access/config"
	"free-access/internal/app/server/grpc/master"
	"free-access/internal/app/server/grpc/node"
	"free-access/internal/app/server/web"
)

func main() {
	config.ParseConfig()
	if config.Config.Common.RunMode == "master" {
		// grpc server
		go master.StartGRPCService(config.Config.Common.TunnelListenPort)
		// web server
		go web.StartHttpServer(config.Config.Common.WebListenPort)
	}
	// init grpc client
	node.StartGRPCClient(config.Config.Common.MasterAddr)
	<-make(chan int)
}
