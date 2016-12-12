package rpc

import (
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"testing"
)

const (
	address  = "localhost:25000"
	test_key = "test-uuid"
)

func TestSnowflake(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := rpc.NewIDGenServerRPCClient(conn)
	// Contact the server and print out its response.
	r, err := c.Next(context.Background(), &rpc.Snowflake_Key{Name: test_key})
	if err != nil {
		t.Fatalf("could not get next value: %v", err)
	}
	t.Log(r.Value)
}

func TestSnowflakeUUID(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := rpc.NewIDGenServerRPCClient(conn)

	// Contact the server and print out its response.
	r, err := c.GetUUID(context.Background(), &rpc.Snowflake_NullRequest{})
	if err != nil {
		t.Fatalf("could not get next value: %v", err)
	}
	t.Logf("%b", r.Uuid)
}
