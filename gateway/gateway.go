package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/gateway/conf"
	"github.com/oikomi/FishChatServer2/libnet"
	"io"
	// "io/ioutil"
	// "github.com/golang/protobuf/proto"
	"encoding/binary"
	"io/ioutil"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		glog.Error("conf.Init() error: ", err)
		panic(err)
	}
	server, err := libnet.Serve("tcp", conf.Conf.Server.Addr, libnet.Packet(2, 1024*1024, 4096, binary.LittleEndian, TestCodec{}))
	if err != nil {
		glog.Error("libnet.Serve error: ", err)
		panic(err)
	}
	glog.Info("server start: ", server.Listener().Addr().String())
	for {
		session, err := server.Accept()
		if err != nil {
			glog.Error("server.Accept error: ", err)
			break
		}
		go func() {
			addr := session.Conn().RemoteAddr().String()
			glog.Info("client ", addr, " connected")
			for {
				var msg []byte
				if err = session.Receive(&msg); err != nil {
					glog.Error("session.Receive error: ", err)
					break
				}
				glog.Info("receive msg : ", string(msg))
				if err = session.Send(msg); err != nil {
					glog.Error("session.Send error: ", err)
					break
				}
			}
			//println("client", addr, "closed")
		}()
	}
}

type TestCodec struct {
}

type TestEncoder struct {
	w io.Writer
}

type TestDecoder struct {
	r io.Reader
}

func (codec TestCodec) NewEncoder(w io.Writer) libnet.Encoder {
	return &TestEncoder{w}
}

func (codec TestCodec) NewDecoder(r io.Reader) libnet.Decoder {
	return &TestDecoder{r}
}

func (encoder *TestEncoder) Encode(msg interface{}) error {
	_, err := encoder.w.Write(msg.([]byte))
	return err
}

func (decoder *TestDecoder) Decode(msg interface{}) error {
	// We use ReadAll() here because we know the reader is a buffer object not a real net.Conn
	d, err := ioutil.ReadAll(decoder.r)
	if err != nil {
		return err
	}
	*(msg.(*[]byte)) = d
	return nil
}
