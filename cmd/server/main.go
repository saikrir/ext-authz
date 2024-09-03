package main

import (
	"fmt"
	"log"
	"net"

	"github.com/saikrir/ext-authz/internal/authsvc"

	authPB "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"
)

func main() {
	addr, port := "0.0.0.0", 3001
	process := fmt.Sprintf("%s:%d", addr, port)

	listener, err := net.Listen("tcp", process)
	if err != nil {
		log.Fatalf("failed to start proccess on port %d, error: %#v ", port, err)
	}
	grpcServer := grpc.NewServer()

	authSvc := authsvc.NewAuthSvc("http://localhost:9999/v1/auth/tokens")
	authPB.RegisterAuthorizationServer(grpcServer, authSvc)
	log.Println("will start grpcserver on", port)
	grpcServer.Serve(listener)
}
