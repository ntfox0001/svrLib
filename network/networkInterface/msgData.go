package networkInterface

import "github.com/golang/protobuf/proto"

type IMsgData interface {
	Marshal() (dAtA []byte, err error)
	MarshalTo(dAtA []byte) (int, error)
	Unmarshal(dAtA []byte) error
	Size() (n int)
	Reset()
	String() string
	ProtoMessage()
}
type RawMsgData struct {
	name string
	raw  []byte
}

func NewRawMsgData(name string, buf []byte) *RawMsgData {
	return &RawMsgData{
		name: name,
		raw:  buf,
	}
}
func (d *RawMsgData) Name() string {
	return d.name
}
func (d *RawMsgData) Unmarshal(pm proto.Message) error {
	return proto.Unmarshal(d.raw, pm)
}
