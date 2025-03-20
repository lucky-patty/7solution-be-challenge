package main

import (
  "context"
  "encoding/json"
  "os"
  "fmt"
  "io/ioutil"
  "net"
  "net/http"
  "net/http/httptest"
  "testing"

  "lucky-patty/7solution/challenge3/loader"
  beefpb "lucky-patty/7solution/challenge3/gen/go/beef"
  "google.golang.org/grpc"
  "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
  "google.golang.org/grpc/test/bufconn"
  
  "github.com/stretchr/testify/assert"
)

type BeefSummary struct {
  Beef map[string]int32 `json:"beef"`
}

type TestServer struct {
  beefpb.UnimplementedBeefServiceServer
}

const bufSize = 1024 * 1024 

var list *bufconn.Listener

func InitGRPCServer(t *testing.T, register func(*grpc.Server)) *grpc.ClientConn {
  list = bufconn.Listen(bufSize)
  s := grpc.NewServer()
  register(s)

  go func() {
    if err := s.Serve(list); err != nil {
      t.Fatalf("Server failed: %v", err)
    }
  }()

  conn, err := grpc.DialContext(context.Background(), "bufnet", 
    grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
      return list.Dial()
    }),
    grpc.WithInsecure(),
  )

  if err != nil {
      t.Fatalf("Failed to dial bufnet: %v", err)
  }

  return conn
}

func loadExpectedResponse(filename string) (BeefSummary, error) {
  file, err := os.Open(filename)
  if err != nil {
    return BeefSummary{}, err 
  }
  defer file.Close()

  bytes, err := ioutil.ReadAll(file)
  if err != nil {
    return BeefSummary{}, err 
  }

  var expectedResponse BeefSummary 
  err = json.Unmarshal(bytes, &expectedResponse) 
  if err != nil {
    return BeefSummary{}, err
  }

  return expectedResponse, nil 
}

func (s *TestServer) GetSummary(ctx context.Context, req *beefpb.Empty) (*beefpb.BeefSummary, error) {
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



func TestGRPC_GetSummary(t *testing.T) {
  // Load 
  expectedResponse, err := loadExpectedResponse("test/expected_beef_summary.json")
  assert.NoError(t, err)

  conn := InitGRPCServer(t, func(s *grpc.Server) {
    beefpb.RegisterBeefServiceServer(s, &TestServer{})
  })

  // defer conn.Close()
  if err != nil {
    t.Fatalf("Failed to connect: %v", err)
  }
  defer conn.Close()

  client := beefpb.NewBeefServiceClient(conn)
  resp, err := client.GetSummary(context.Background(), &beefpb.Empty{})
  assert.NoError(t, err)

  assert.Equal(t, expectedResponse.Beef, resp.Beef)

  t.Logf("Response: %+v \n", resp.Beef)
}


func TestREST_GetSummary(t *testing.T) {
  // Load 
  expectedResponse, err := loadExpectedResponse("test/expected_beef_summary.json")
  assert.NoError(t, err)

  conn := InitGRPCServer(t, func(s *grpc.Server) {
    beefpb.RegisterBeefServiceServer(s, &TestServer{})
  })

  // defer conn.Close()
  if err != nil {
    t.Fatalf("Failed to connect: %v", err)
  }
  defer conn.Close()

  // Create REST 
  ctx := context.Background()
  mux := runtime.NewServeMux()
  err = beefpb.RegisterBeefServiceHandler(ctx, mux, conn)
  assert.NoError(t, err)


  // Start HTTP 
  server := httptest.NewServer(mux)
  defer server.Close() 

  // Make GET request 
  resp, err := http.Get(server.URL + "/beef/summary")
  assert.NoError(t, err) 
  defer resp.Body.Close() 

  // Read resp 
  body, err := ioutil.ReadAll(resp.Body)
  assert.NoError(t, err)

  var actualResponse BeefSummary 
  err = json.Unmarshal(body, &actualResponse)
  assert.NoError(t, err)

  // Compare 
  assert.Equal(t, expectedResponse.Beef, actualResponse.Beef)

  t.Logf("REST Response: %+v \n", actualResponse.Beef)
}
