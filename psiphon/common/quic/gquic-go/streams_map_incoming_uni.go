// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/bassosimone/psiphon-tunnel-core/internal/github.com/cheekybits/genny

package gquic

import (
	"fmt"
	"sync"

	"github.com/bassosimone/psiphon-tunnel-core/psiphon/common/quic/gquic-go/internal/protocol"
	"github.com/bassosimone/psiphon-tunnel-core/psiphon/common/quic/gquic-go/internal/wire"
)

type incomingUniStreamsMap struct {
	mutex sync.RWMutex
	cond  sync.Cond

	streams map[protocol.StreamID]receiveStreamI

	nextStream    protocol.StreamID // the next stream that will be returned by AcceptStream()
	highestStream protocol.StreamID // the highest stream that the peer openend
	maxStream     protocol.StreamID // the highest stream that the peer is allowed to open
	maxNumStreams int               // maximum number of streams

	newStream        func(protocol.StreamID) receiveStreamI
	queueMaxStreamID func(*wire.MaxStreamIDFrame)

	closeErr error
}

func newIncomingUniStreamsMap(
	nextStream protocol.StreamID,
	initialMaxStreamID protocol.StreamID,
	maxNumStreams int,
	queueControlFrame func(wire.Frame),
	newStream func(protocol.StreamID) receiveStreamI,
) *incomingUniStreamsMap {
	m := &incomingUniStreamsMap{
		streams:          make(map[protocol.StreamID]receiveStreamI),
		nextStream:       nextStream,
		maxStream:        initialMaxStreamID,
		maxNumStreams:    maxNumStreams,
		newStream:        newStream,
		queueMaxStreamID: func(f *wire.MaxStreamIDFrame) { queueControlFrame(f) },
	}
	m.cond.L = &m.mutex
	return m
}

func (m *incomingUniStreamsMap) AcceptStream() (receiveStreamI, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var str receiveStreamI
	for {
		var ok bool
		if m.closeErr != nil {
			return nil, m.closeErr
		}
		str, ok = m.streams[m.nextStream]
		if ok {
			break
		}
		m.cond.Wait()
	}
	m.nextStream += 4
	return str, nil
}

func (m *incomingUniStreamsMap) GetOrOpenStream(id protocol.StreamID) (receiveStreamI, error) {
	m.mutex.RLock()
	if id > m.maxStream {
		m.mutex.RUnlock()
		return nil, fmt.Errorf("peer tried to open stream %d (current limit: %d)", id, m.maxStream)
	}
	// if the id is smaller than the highest we accepted
	// * this stream exists in the map, and we can return it, or
	// * this stream was already closed, then we can return the nil
	if id <= m.highestStream {
		s := m.streams[id]
		m.mutex.RUnlock()
		return s, nil
	}
	m.mutex.RUnlock()

	m.mutex.Lock()
	// no need to check the two error conditions from above again
	// * maxStream can only increase, so if the id was valid before, it definitely is valid now
	// * highestStream is only modified by this function
	var start protocol.StreamID
	if m.highestStream == 0 {
		start = m.nextStream
	} else {
		start = m.highestStream + 4
	}
	for newID := start; newID <= id; newID += 4 {
		m.streams[newID] = m.newStream(newID)
		m.cond.Signal()
	}
	m.highestStream = id
	s := m.streams[id]
	m.mutex.Unlock()
	return s, nil
}

func (m *incomingUniStreamsMap) DeleteStream(id protocol.StreamID) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.streams[id]; !ok {
		return fmt.Errorf("Tried to delete unknown stream %d", id)
	}
	delete(m.streams, id)
	// queue a MAX_STREAM_ID frame, giving the peer the option to open a new stream
	if numNewStreams := m.maxNumStreams - len(m.streams); numNewStreams > 0 {
		m.maxStream = m.highestStream + protocol.StreamID(numNewStreams*4)
		m.queueMaxStreamID(&wire.MaxStreamIDFrame{StreamID: m.maxStream})
	}
	return nil
}

func (m *incomingUniStreamsMap) CloseWithError(err error) {
	m.mutex.Lock()
	m.closeErr = err
	for _, str := range m.streams {
		str.closeForShutdown(err)
	}
	m.mutex.Unlock()
	m.cond.Broadcast()
}
