package main 

import (
  "lucky-patty/7solution/challenge3/loader"
  "context"
  "log"
  "net"
  "net/http"
  "fmt"

  "google.golang.org/grpc"
  "google.golang.org/grpc/reflection"
  "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

  beefpb "lucky-patty/7solution/challenge3/gen/go/beef"
)

type Server struct {
  beefpb.UnimplementedBeefServiceServer
}

func (s *Server) GetSummary(ctx context.Context, req *beefpb.Empty) (*beefpb.BeefSummary, error) {
  res, err := loader.TextLoader()
  if err != nil {
    fmt.Println("Error: ", err)
    return nil, err 
  }
  return &beefpb.BeefSummary{
    Beef: map[string]int32 {
      "t-bone": res["t-bone"],
      "fatback": res["fatback"],
      "pastrami": res["pastrami"],
      "pork": res["pork"],
      "meatloaf": res["meatloaf"],
      "jowl": res["jowl"],
      "enim": res["enim"],
      "bresaola": res["bresaola"],
    },
  }, nil
}

func startGRPCServer() {
  listener, err := net.Listen("tcp", ":50051")
  if err != nil {
    log.Fatalf("Failed to listen: %v \n", err)
  }

  grpcServer := grpc.NewServer()
  beefpb.RegisterBeefServiceServer(grpcServer, &Server{})
  reflection.Register(grpcServer)

  fmt.Println("gRPC server is running on port 50051")

  if err := grpcServer.Serve(listener); err != nil {
    log.Fatalf("Failed to serve: %v", err)
  }
}

func startRESTServer() {
  ctx := context.Background()
  mux := runtime.NewServeMux()
  opts := []grpc.DialOption{grpc.WithInsecure()}

  err := beefpb.RegisterBeefServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
  if err != nil {
    log.Fatalf("Failed to start REST Server: %v ", err)
  }

  log.Println("REST API Listening on port 8080")
  if err := http.ListenAndServe(":8080", mux); err != nil {
    log.Fatalf("Failed to start REST Server: %v", err)
  }

}

func main() {
  go startRESTServer()
  startGRPCServer()
}

