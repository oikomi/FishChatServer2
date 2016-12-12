package rpc

import (
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/idgen/conf"
	"github.com/oikomi/FishChatServer2/server/idgen/dao"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type RPCServer struct {
	dao        *dao.Dao
	machine_id uint64 // 10-bit machine id
	// client_pool chan etcd.KeysAPI
	ch_proc chan chan uint64
}

const (
	SERVICE        = "[SNOWFLAKE]"
	ENV_MACHINE_ID = "MACHINE_ID" // specific machine id
	PATH           = "/seqs/"
	UUID_KEY       = "/seqs/snowflake-uuid"
	BACKOFF        = 100  // max backoff delay millisecond
	CONCURRENT     = 128  // max concurrent connections to etcd
	UUID_QUEUE     = 1024 // uuid process queue
)

const (
	TS_MASK         = 0x1FFFFFFFFFF // 41bit
	SN_MASK         = 0xFFF         // 12bit
	MACHINE_ID_MASK = 0x3FF         // 10bit
)

func (s *RPCServer) init() {
	// s.client_pool = make(chan etcd.KeysAPI, CONCURRENT)
	s.ch_proc = make(chan chan uint64, UUID_QUEUE)
	// init client pool
	// for i := 0; i < CONCURRENT; i++ {
	// 	s.client_pool <- etcdclient.KeysAPI()
	// }

	// check if user specified machine id is set
	// if env := os.Getenv(ENV_MACHINE_ID); env != "" {
	// 	if id, err := strconv.Atoi(env); err == nil {
	// 		s.machine_id = (uint64(id) & MACHINE_ID_MASK) << 12
	// 		log.Info("machine id specified:", id)
	// 	} else {
	// 		log.Panic(err)
	// 		os.Exit(-1)
	// 	}
	// } else {
	s.initMachineID()
	// }
	go s.uuid_task()
}

func (s *RPCServer) initMachineID() {
	// client := <-s.client_pool
	// defer func() { s.client_pool <- client }()
	var prevIndex int64
	var prevValue int
	for {
		// get the key
		resp, err := s.dao.Etcd.EtcCli.Get(context.Background(), UUID_KEY)
		if err != nil {
			glog.Error(err)
			panic(err)
		}
		for _, value := range resp.Kvs {
			prevValue, err = strconv.Atoi(string(value.Value))
			if err != nil {
				glog.Error(err)
				return
			}
			prevIndex = value.ModRevision
			glog.Info(prevValue)
			glog.Info(prevIndex)
		}
		// get prevValue & prevIndex
		// prevValue, err := strconv.Atoi(resp.Node.Value)
		// if err != nil {
		// 	log.Panic(err)
		// 	os.Exit(-1)
		// }
		// prevIndex := resp.Node.ModifiedIndex

		// CompareAndSwap
		_, err = s.dao.Etcd.EtcCli.Put(context.Background(), UUID_KEY, fmt.Sprint(prevValue+1), clientv3.WithRev(prevIndex))
		if err != nil {
			cas_delay()
			continue
		}
		// record serial number of this service, already shifted
		s.machine_id = (uint64(prevValue+1) & MACHINE_ID_MASK) << 12
		return
	}
}

// get next value of a key, like auto-increment in mysql
func (s *RPCServer) Next(ctx context.Context, in *rpc.Snowflake_Key) (sfRes *rpc.Snowflake_Value, err error) {
	// client := <-s.client_pool
	// defer func() { s.client_pool <- client }()
	var prevIndex int64
	var prevValue int
	var resp *clientv3.GetResponse
	key := PATH + in.Name
	for {
		glog.Info(s.dao.Etcd.EtcCli)
		if s.dao.Etcd.EtcCli == nil {
			glog.Error("s.dao.Etcd.EtcCli == nil")
			return
		}
		// get the key
		resp, err = s.dao.Etcd.EtcCli.Get(context.Background(), key)
		if err != nil {
			glog.Error(err)
			return nil, errors.New("Key not exists, need to create first")
		}
		// get prevValue & prevIndex
		for _, value := range resp.Kvs {
			prevValue, err = strconv.Atoi(string(value.Value))
			if err != nil {
				glog.Error(err)
				return nil, errors.New("marlformed value")
			}
			prevIndex = value.ModRevision
		}
		// prevValue, err := strconv.Atoi(resp.Node.Value)
		// if err != nil {
		// 	glog.Error(err)
		// 	return nil, errors.New("marlformed value")
		// }
		// prevIndex := resp.Node.ModifiedIndex
		// CompareAndSwap
		_, err = s.dao.Etcd.EtcCli.Put(context.Background(), key, fmt.Sprint(prevValue+1), clientv3.WithRev(prevIndex))
		if err != nil {
			cas_delay()
			continue
		}
		sfRes = &rpc.Snowflake_Value{int64(prevValue + 1)}
		return
	}
}

// generate an unique uuid
func (s *RPCServer) GetUUID(context.Context, *rpc.Snowflake_NullRequest) (*rpc.Snowflake_UUID, error) {
	req := make(chan uint64, 1)
	s.ch_proc <- req
	return &rpc.Snowflake_UUID{<-req}, nil
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
			glog.Error("clock shift happened, waiting until the clock moving to the next millisecond.")
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
	glog.Info("[idgen] rpc server init: ", conf.Conf.RPCServer.Addr)
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	s := grpc.NewServer()
	d, err := dao.NewDao()
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	rpcServer := &RPCServer{
		dao: d,
	}
	rpcServer.init()
	rpc.RegisterIDGenServerRPCServer(s, rpcServer)
	s.Serve(lis)
}
