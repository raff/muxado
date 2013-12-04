package frame

import (
	"encoding/binary"
)

const (
	// syn frames are actually longer, but they are variable length
	synFrameSize = headerSize
)

type RStreamSyn struct {
	Header

	relatedStreamId StreamId
	streamPriority  StreamPriority
}

// RelatedStreamId returns the related stream's id.
// A zero value means no related stream was specified
func (f *RStreamSyn) RelatedStreamId() StreamId {
	return f.relatedStreamId
}

// StreamPriority returns the stream priority set on this frame
func (f *RStreamSyn) StreamPriority() StreamPriority {
	return f.streamPriority
}

func (f *RStreamSyn) readFrom(d deserializer) (err error) {
	if f.Flags().IsSet(flagRelatedStream) {
		if err := binary.Read(d, order, &f.relatedStreamId); err != nil {
			return err
		}
	}

	if f.Flags().IsSet(flagStreamPriority) {
		if err := binary.Read(d, order, &f.streamPriority); err != nil {
			return err
		}
	}

	return
}

type WStreamSyn struct {
	Header
	fixed   [synFrameSize]byte
	toWrite []byte // when writing, you just pass a byte slice to write
}

func (f *WStreamSyn) writeTo(s serializer) (err error) {
	if _, err = s.Write(f.fixed[:]); err != nil {
		return err
	}

	if _, err = s.Write(f.toWrite); err != nil {
		return err
	}

	return
}

func (f *WStreamSyn) Set(streamId, relatedStreamId StreamId, streamPriority StreamPriority, fin bool, info []byte) (err error) {
	var (
		flags  flagsType
		length int
	)

	// set fin bit
	if fin {
		flags.Set(flagFin)
	}

	// validate the related stream
	if relatedStreamId != 0 {
		if relatedStreamId > streamMask {
			err = protoError("Related stream id %d is out of range", relatedStreamId)
			return
		}

		flags.Set(flagRelatedStream)
		length += 4
	}

	// validate the stream priority
	if streamPriority != 0 {
		if streamPriority > priorityMask {
			err = protoError("Priority %d is out of range", streamPriority)
			return
		}

		flags.Set(flagStreamPriority)
		length += 4
	}

	if length > 0 {
		f.toWrite = make([]byte, length)
		p := 0

		if flags.IsSet(flagRelatedStream) {
			order.PutUint32(f.toWrite[p:], uint32(relatedStreamId))
			p += 4
		}

		if flags.IsSet(flagStreamPriority) {
			order.PutUint32(f.toWrite[p:], uint32(streamPriority))
			p += 4
		}
	}

	// make the frame
	if err = f.Header.SetAll(TypeStreamSyn, length, streamId, flags); err != nil {
		return
	}

	return
}

func NewWStreamSyn() (f *WStreamSyn) {
	f = new(WStreamSyn)
	f.Header = Header(f.fixed[:headerSize])
	return
}
