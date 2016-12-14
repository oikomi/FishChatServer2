package main

import (
	"flag"
	"fmt"
	wonaming "github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"github.com/oikomi/FishChatServer2/service_discovery/example/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	// "time"
)

var (
	serv = flag.String("service", "hello_service", "service name")
	// reg  = flag.String("reg", "127.0.0.1:8500", "register address")
	reg = flag.String("reg", "http://127.0.0.1:2379", "register address")
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

func main() {
	flag.Parse()
	fmt.Println("----1----")
	r := wonaming.NewResolver(*serv)
	fmt.Println("----2----")
	b := grpc.RoundRobin(r)
	fmt.Println("----3----")
	conn, err := grpc.Dial(*reg, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		panic(err)
	}
	fmt.Println("----4----")
	client := pb.NewHelloServiceClient(conn)
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Greeting: "world"})
	if err == nil {
		fmt.Printf("Reply is %s\n", resp.Reply)
	}

	// ticker := time.NewTicker(2 * time.Second)
	// for t := range ticker.C {
	// 	fmt.Println("----1----")
	// 	client := pb.NewHelloServiceClient(conn)
	// 	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Greeting: "world"})
	// 	if err == nil {
	// 		fmt.Printf("%v: Reply is %s\n", t, resp.Reply)
	// 	}
	// }
}
