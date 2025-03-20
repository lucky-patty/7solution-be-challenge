package main 

import (
  "fmt"
  "log"
  "net"

  "google.golang.org/grpc"
  "lucky-patty/7solution/challenge3/gen/go/beef"
  "lucky-patty/7solution/challenge3/internal/service"
)

func main() {
  listener, err := net.Listen("tcp", ":50051")
  if err != nil {
    log.Fatalf("Failed to listen: %v \n", err)
  }

  grpcServer := grpc.NewServer()
  beefpb.RegisterBeefServiceServer(grpcServer, &service.BeefServiceServer{})

  fmt.Println("gRPC server is running on port 50051")

  if err := grpcServer.Serve(listener); err != nil {
    log.Fatalf("Failed to serve: %v", err)
  }
}
