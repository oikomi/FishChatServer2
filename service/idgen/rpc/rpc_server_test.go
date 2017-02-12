package rpc

// import (
// 	"fmt"
// 	"github.com/oikomi/FishChatServer2/protocol/rpc"
// 	"golang.org/x/net/context"
// 	"google.golang.org/grpc"
// 	"testing"
// )

// const (
// 	address  = "localhost:31000"
// 	test_key = "test-uuid"
// )

// func TestSnowflake(t *testing.T) {
// 	// Set up a connection to the server.
// 	conn, err := grpc.Dial(address, grpc.WithInsecure())
// 	if err != nil {
// 		t.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()
// 	c := rpc.NewIDGenServerRPCClient(conn)
// 	// Contact the server and print out its response.
// 	r, err := c.Next(context.Background(), &rpc.Snowflake_Key{Name: test_key})
// 	if err != nil {
// 		t.Fatalf("could not get next value: %v", err)
// 	}
// 	fmt.Println(r.Value)
// 	t.Log(r.Value)
// }

// func BenchmarkSnowflake(b *testing.B) {
// 	// Set up a connection to the server.
// 	conn, err := grpc.Dial(address, grpc.WithInsecure())
// 	if err != nil {
// 		b.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()
// 	c := rpc.NewIDGenServerRPCClient(conn)
// 	for i := 0; i < b.N; i++ {
// 		// Contact the server and print out its response.
// 		r, err := c.Next(context.Background(), &rpc.Snowflake_Key{Name: test_key})
// 		if err != nil {
// 			b.Fatalf("could not get next value: %v", err)
// 		}
// 		fmt.Println(r.Value)
// 	}
// }

// func TestSnowflakeUUID(t *testing.T) {
// 	// Set up a connection to the server.
// 	conn, err := grpc.Dial(address, grpc.WithInsecure())
// 	if err != nil {
// 		t.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()
// 	c := rpc.NewIDGenServerRPCClient(conn)

// 	// Contact the server and print out its response.
// 	r, err := c.GetUUID(context.Background(), &rpc.Snowflake_NullRequest{})
// 	if err != nil {
// 		t.Fatalf("could not get next value: %v", err)
// 	}
// 	fmt.Println(r.Uuid)
// 	t.Logf("%b", r.Uuid)
// }

// func BenchmarkSnowflakeUUID(b *testing.B) {
// 	// Set up a connection to the server.
// 	conn, err := grpc.Dial(address, grpc.WithInsecure())
// 	if err != nil {
// 		b.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()
// 	c := rpc.NewIDGenServerRPCClient(conn)
// 	for i := 0; i < b.N; i++ {
// 		// Contact the server and print out its response.
// 		r, err := c.GetUUID(context.Background(), &rpc.Snowflake_NullRequest{})
// 		if err != nil {
// 			b.Fatalf("could not get uuid: %v", err)
// 		}
// 		fmt.Println(r.Uuid)
// 	}
// }
