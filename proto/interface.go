package proto

import (
	"github.com/raff/muxado/proto/frame"
	"net"
	"time"
)

type IStream interface {
	Write([]byte) (int, error)
	Read([]byte) (int, error)
	Close() error
	SetDeadline(time.Time) error
	SetReadDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
	HalfClose([]byte) (int, error)
	Id() frame.StreamId
	RelatedStreamId() frame.StreamId
	StreamInfo() frame.StreamInfo
	Session() ISession
	RemoteAddr() net.Addr
	LocalAddr() net.Addr
}

type ISession interface {
	Open() (IStream, error)
	OpenEx([]byte) (IStream, error)
	OpenStream(frame.StreamPriority, frame.StreamId, bool, []byte) (IStream, error)
	Accept() (IStream, error)
	Kill() error
	GoAway(frame.ErrorCode, []byte) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close() error
	Wait() (frame.ErrorCode, error, []byte)
	NetListener() net.Listener
	NetDial(network, addr string) (net.Conn, error)
}
