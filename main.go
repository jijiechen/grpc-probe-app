// Copied from https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_server/main.go
package main

import (
	"flag"
	"fmt"
	pb "github.com/jijiechen/grpc-probe-app/proto"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
)

var (
	port = flag.Int("port", 5085, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcS := grpc.NewServer()

	s := &server{}
	go s.startStatusTicker()
	pb.RegisterGreeterServer(grpcS, s)
	healthpb.RegisterHealthServer(grpcS, s)

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcS.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
