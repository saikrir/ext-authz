package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/saikrir/ext-authz/internal/authsvc"

	authPB "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"
)

const AuthHostEnvVar = "AuthHost"
const AuthApiUrl = "http://%s:9999/v1/auth/tokens"

func Run() error {

	authSvcHost := os.Getenv(AuthHostEnvVar)

	if len(authSvcHost) == 0 {
		return fmt.Errorf("failed to find %s variable ", AuthHostEnvVar)
	}

	addr, port := "0.0.0.0", 3001
	process := fmt.Sprintf("%s:%d", addr, port)

	listener, err := net.Listen("tcp", process)
	if err != nil {
		log.Printf("failed to start proccess on port %d, error: %#v ", port, err)
		return err
	}
	grpcServer := grpc.NewServer()

	authSvc := authsvc.NewAuthSvc(fmt.Sprintf(AuthApiUrl, authSvcHost))
	authPB.RegisterAuthorizationServer(grpcServer, authSvc)
	log.Println("will start grpcserver on", port)
	return grpcServer.Serve(listener)
}

func main() {
	if err := Run(); err != nil {
		log.Fatal("failed to start extauth ", err)
	}
}
