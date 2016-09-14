package codec

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/libnet"
	"io"
	"reflect"
)

type ProtobufProtocol struct {
	types map[string]reflect.Type
	names map[reflect.Type]string
}

func Protobuf() *ProtobufProtocol {
	return &ProtobufProtocol{
		types: make(map[string]reflect.Type),
		names: make(map[reflect.Type]string),
	}
}

func (j *ProtobufProtocol) Register(t interface{}) {
	rt := reflect.TypeOf(t)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	name := rt.PkgPath() + "/" + rt.Name()
	j.types[name] = rt
	j.names[rt] = name
}

func (j *ProtobufProtocol) RegisterName(name string, t interface{}) {
	rt := reflect.TypeOf(t)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	j.types[name] = rt
	j.names[rt] = name
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

func (c *protobufCodec) Receive() ([]byte, error) {
	data := c.r.ReadPacket(SplitByUint16BE)
	return data, nil
}

func (c *protobufCodec) Send(msg interface{}) error {
	data, err := proto.Marshal(msg.(proto.Message))
	if err != nil {
		glog.Error(err)
		return err
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
