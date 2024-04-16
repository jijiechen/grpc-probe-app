package main

import (
	"context"
	pb "github.com/jijiechen/grpc-probe-app/proto"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"time"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
	healthpb.UnimplementedHealthServer

	status healthpb.HealthCheckResponse_ServingStatus
	c      <-chan struct{}
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	log.Printf("handling probe for service: %s", req.GetService())
	return s.getStatusResp(), nil
}

func (s *server) Watch(req *healthpb.HealthCheckRequest, server healthpb.Health_WatchServer) error {
	log.Printf("handling probe for service: %s", req.GetService())
	if err := server.SendMsg(s.getStatusResp()); err != nil {
		return err
	}

	for {
		select {
		case <-s.c:
			if err := server.SendMsg(s.getStatusResp()); err != nil {
				return err
			}
		}
	}
}

func (s *server) startStatusTicker() {
	tickerChannel := make(chan struct{})
	s.c = tickerChannel
	s.status = healthpb.HealthCheckResponse_NOT_SERVING

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if s.status == healthpb.HealthCheckResponse_SERVING {
				s.status = healthpb.HealthCheckResponse_NOT_SERVING
			} else {
				s.status = healthpb.HealthCheckResponse_SERVING
			}

			// don't block the ticker...
			select {
			case tickerChannel <- struct{}{}:
			default:
			}
		}
	}
}

func (s *server) getStatusResp() *healthpb.HealthCheckResponse {
	return &healthpb.HealthCheckResponse{
		Status: s.status,
	}
}
