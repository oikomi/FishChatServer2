package codec

import (
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
	}
	codec.closer, _ = rw.(io.Closer)
	return codec
}

type protobufCodec struct {
	p      *ProtobufProtocol
	closer io.Closer
}

func (c *protobufCodec) Receive() (interface{}, error) {
	var in gobMsg
	err := c.decoder.Decode(&in)
	if err != nil {
		return nil, err
	}
	var message interface{}
	if t, exists := c.p.types[in.Type]; exists {
		message = reflect.New(t).Interface()
	} else {
		return nil, ErrGobUnknow
	}
	err = c.decoder.Decode(message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (c *protobufCodec) Send(msg interface{}) error {
	var out gobMsg
	t := reflect.TypeOf(msg)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if name, exists := c.p.names[t]; exists {
		out.Type = name
	} else {
		return ErrGobUnknow
	}
	err := c.encoder.Encode(&out)
	if err != nil {
		return err
	}
	return c.encoder.Encode(msg)
}
func (c *protobufCodec) Close() error {
	if c.closer != nil {
		return c.closer.Close()
	}
	return nil
}
