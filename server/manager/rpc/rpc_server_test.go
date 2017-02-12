package rpc

// import (
// 	"github.com/oikomi/FishChatServer2/protocol/rpc"
// 	"golang.org/x/net/context"
// 	"google.golang.org/grpc"
// 	"testing"
// )

// func TestSync(t *testing.T) {
// 	// Set up a connection to the server.
// 	conn, err := grpc.Dial("127.0.0.1:24000", grpc.WithInsecure())
// 	if err != nil {
// 		t.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()
// 	m := rpc.NewManagerServerRPCClient(conn)
// 	// Contact the server and print out its response.
// 	r, err := m.Sync(context.Background(), &rpc.MGSyncMsgReq{UID: 123, CurrentID: 25, TotalID: 26})
// 	if err != nil {
// 		t.Fatalf("could not get next value: %v", err)
// 	}
// 	t.Log(r)
// }
