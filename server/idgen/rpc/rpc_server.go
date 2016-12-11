package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/idgen/conf"
	"github.com/oikomi/FishChatServer2/server/idgen/dao"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
	dao *dao.Dao
}

// get next value of a key, like auto-increment in mysql
func (s *RPCServer) Next(ctx context.Context, in *rpc.Snowflake_Key) (*rpc.Snowflake_Value, error) {
	client := <-s.client_pool
	defer func() { s.client_pool <- client }()
	key := PATH + in.Name
	for {
		// get the key
		resp, err := client.Get(context.Background(), key, nil)
		if err != nil {
			log.Error(err)
			return nil, errors.New("Key not exists, need to create first")
		}
		// get prevValue & prevIndex
		prevValue, err := strconv.Atoi(resp.Node.Value)
		if err != nil {
			log.Error(err)
			return nil, errors.New("marlformed value")
		}
		prevIndex := resp.Node.ModifiedIndex
		// CompareAndSwap
		resp, err = client.Set(context.Background(), key, fmt.Sprint(prevValue+1), &etcd.SetOptions{PrevIndex: prevIndex})
		if err != nil {
			cas_delay()
			continue
		}
		return &pb.Snowflake_Value{int64(prevValue + 1)}, nil
	}
}

// generate an unique uuid
func (s *RPCServer) GetUUID(context.Context, *rpc.Snowflake_NullRequest) (*rpc.Snowflake_UUID, error) {
	req := make(chan uint64, 1)
	s.ch_proc <- req
	return &pb.Snowflake_UUID{<-req}, nil
}

// uuid generator
func (s *RPCServer) uuid_task() {
	var sn uint64     // 12-bit serial no
	var last_ts int64 // last timestamp
	for {
		ret := <-s.ch_proc
		// get a correct serial number
		t := ts()
		if t < last_ts { // clock shift backward
			log.Error("clock shift happened, waiting until the clock moving to the next millisecond.")
			t = s.wait_ms(last_ts)
		}
		if last_ts == t { // same millisecond
			sn = (sn + 1) & SN_MASK
			if sn == 0 { // serial number overflows, wait until next ms
				t = s.wait_ms(last_ts)
			}
		} else { // new millsecond, reset serial number to 0
			sn = 0
		}
		// remember last timestamp
		last_ts = t
		// generate uuid, format:
		//
		// 0		0.................0		0..............0	0........0
		// 1-bit	41bit timestamp			10bit machine-id	12bit sn
		var uuid uint64
		uuid |= (uint64(t) & TS_MASK) << 22
		uuid |= s.machine_id
		uuid |= sn
		ret <- uuid
	}
}

// wait_ms will spin wait till next millisecond.
func (s *RPCServer) wait_ms(last_ts int64) int64 {
	t := ts()
	for t <= last_ts {
		t = ts()
	}
	return t
}

////////////////////////////////////////////////////////////////////////////////
// random delay
func cas_delay() {
	<-time.After(time.Duration(rand.Int63n(BACKOFF)) * time.Millisecond)
}

// get timestamp
func ts() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func RPCServerInit() {
	glog.Info("[manager] rpc server init: ", conf.Conf.RPCServer.Addr)
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	s := grpc.NewServer()
	rpcServer := &RPCServer{
		dao: dao.NewDao(),
	}
	rpc.RegisterIDGenServerRPCServer(s, rpcServer)
	s.Serve(lis)
}
