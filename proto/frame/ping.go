package frame

import "io"

const (
	pingBodySize  = 8
	pingFrameSize = headerSize + pingBodySize
)

var (
        defaultData [pingBodySize]byte
)

// Ping message
type RPing struct {
	Header
	body  [pingBodySize]byte
}

func (f *RPing) Body() []byte {
	return f.body[:]
}

func (f *RPing) readFrom(d deserializer) (err error) {
	if f.Length() != pingBodySize {
		return protoError("PING length must be %d, got %d", pingBodySize, f.Length())
	}

	_, err = io.ReadFull(d, f.body[:])
	return
}

type WPing struct {
	Header
        fixed [pingFrameSize]byte
}

func NewWPing() (f *WPing) {
	f = new(WPing)
	f.Header = f.fixed[:headerSize]
	return
}

func (f *WPing) writeTo(s serializer) (err error) {
	_, err = s.Write(f.fixed[:])
        return
}

func (f *WPing) Set(streamId StreamId, data []byte, ack bool) (err error) {
	var flags flagsType
	if ack {
		flags.Set(flagAck)
	}

        if len(data) != 0 {
            if len(data) != pingBodySize {
                return protoError("PING length must be %d, got %d", pingBodySize, len(data))
            }

	    copy(f.fixed[headerSize:], data)
        }

	if err = f.Header.SetAll(TypePing, pingBodySize, streamId, flags); err != nil {
		return
	}

	return
}
