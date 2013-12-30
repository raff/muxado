package muxado

import (
	"time"

	"github.com/raff/muxado/proto"
	"github.com/raff/muxado/proto/frame"
)

// streamAdaptor recasts the types of some function calls by the proto/Stream implementation
// so that it satisfies the public interface
type streamAdaptor struct {
	proto.IStream
}

func (a *streamAdaptor) Id() StreamId {
	return StreamId(a.IStream.Id())
}

func (a *streamAdaptor) RelatedStreamId() StreamId {
	return StreamId(a.IStream.RelatedStreamId())
}

func (a *streamAdaptor) StreamInfo() StreamInfo {
	return StreamInfo(a.IStream.StreamInfo())
}

func (a *streamAdaptor) Session() Session {
	return &sessionAdaptor{a.IStream.Session()}
}

// sessionAdaptor recasts the types of some function calls by the proto/Session implementation
// so that it satisfies the public interface
type sessionAdaptor struct {
	proto.ISession
}

func (a *sessionAdaptor) Accept() (Stream, error) {
	str, err := a.ISession.Accept()
	return &streamAdaptor{str}, err
}

func (a *sessionAdaptor) Open() (Stream, error) {
	str, err := a.ISession.Open()
	return &streamAdaptor{str}, err
}

func (a *sessionAdaptor) OpenEx(info []byte) (Stream, error) {
	str, err := a.ISession.OpenEx(info)
	return &streamAdaptor{str}, err
}

func (a *sessionAdaptor) OpenStream(priority StreamPriority, related StreamId, fin bool, info []byte) (Stream, error) {
	str, err := a.ISession.OpenStream(frame.StreamPriority(priority), frame.StreamId(related), fin, info)
	return &streamAdaptor{str}, err
}

func (a *sessionAdaptor) GoAway(code ErrorCode, debug []byte) error {
	return a.ISession.GoAway(frame.ErrorCode(code), debug)
}

func (a *sessionAdaptor) Wait() (ErrorCode, error, []byte) {
	code, err, debug := a.ISession.Wait()
	return ErrorCode(code), err, debug
}

func (a *sessionAdaptor) SetInactivityTime(duration time.Duration) error {
	return a.ISession.SetInactivityTime(duration)
}
