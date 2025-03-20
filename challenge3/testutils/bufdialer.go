package testutils 

import (
  "context"
  "net"
  "testing"

  beefpb "lucky-patty/7solution/challenge3/gen/go/beef"
  
  "google.golang.org/grpc"
  "google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024 

var list *bufconn.Listener

func InitGRPCServer(t *testing.T, server func(*grpc.Server)) *grpc.ClientConn {
  lis = bufconn.Listen(bufSize)
  s := grpc.NewServer()

  server(s)

  go func() {
    if err := s.Serve(lis); err != nil {
      t.Fatalf("Server failed: %v", err)
    }
  }()

  conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
    return lis.Dial()
  }),
  grpc.WithInsecure(),
  )
  if err != nil {
    t.Failf("Failed to dial bufnet: %v", err)
  }

  return conn
}


func BufDialer(context.Context, string) (grpc.ClientConnInterface, error) {
  return grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
    return lis.Dial()
  }), grpc.WithInsecure())
}
