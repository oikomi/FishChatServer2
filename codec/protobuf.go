package codec

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/libnet"
	"io"
)

type ProtobufProtocol struct {
}

func Protobuf() *ProtobufProtocol {
	return &ProtobufProtocol{}
}

func (p *ProtobufProtocol) NewCodec(rw io.ReadWriter) libnet.Codec {
	codec := &protobufCodec{
		p: p,
		w: NewWriter(rw),
		r: NewReader(rw),
	}
	codec.closer, _ = rw.(io.Closer)
	return codec
}

type protobufCodec struct {
	p      *ProtobufProtocol
	w      *Writer
	r      *Reader
	closer io.Closer
}

func (c *protobufCodec) Receive() (interface{}, error) {
	data := c.r.ReadPacket(SplitByUint16BE)
	if data != nil {
		glog.Info(string(data))
	}

	return data, nil
}

func (c *protobufCodec) Send(msg interface{}) error {
	data, err := proto.Marshal(msg.(proto.Message))
	if err != nil {
		glog.Error(err)
	}
	c.w.WritePacket(data, SplitByUint16BE)

	return nil
}
func (c *protobufCodec) Close() error {
	if c.closer != nil {
		return c.closer.Close()
	}
	return nil
}
