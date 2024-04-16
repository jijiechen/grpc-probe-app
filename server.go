package main

import (
	"context"
	pb "github.com/jijiechen/grpc-probe-app/proto"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"sync"
	"time"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
	healthpb.UnimplementedHealthServer

	status        healthpb.HealthCheckResponse_ServingStatus
	streamCounter int
	watchers      map[int]chan struct{}
	lock          sync.Mutex
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

	sId := s.newStream()

	for {
		select {
		case <-s.watchers[sId]:
			if err := server.SendMsg(s.getStatusResp()); err != nil {
				return err
			}
		case <-server.Context().Done():
			s.destroyStream(sId)
			return nil
		}
	}
}

func (s *server) getStatusResp() *healthpb.HealthCheckResponse {
	return &healthpb.HealthCheckResponse{
		Status: s.status,
	}
}

func (s *server) newStream() int {
	defer s.lock.Unlock()
	s.lock.Lock()
	s.streamCounter++

	if s.watchers == nil {
		s.watchers = make(map[int]chan struct{})
	}
	s.watchers[s.streamCounter] = make(chan struct{})

	return s.streamCounter
}

func (s *server) destroyStream(sid int) {
	defer s.lock.Unlock()
	s.lock.Lock()

	delete(s.watchers, sid)
}

func (s *server) notifyWatchers() {
	defer s.lock.Unlock()
	s.lock.Lock()
	for _, ch := range s.watchers {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

func (s *server) startStatusTicker() {
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
			go s.notifyWatchers()
		}
	}
}
