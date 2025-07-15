//从pion官方拿的，需要辅助fmtp和pion_extend.go才能工作，目的在于重写和优化部分逻辑 20250711 陈浩 初始化

// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

//go:build !js
// +build !js

package pion

import (
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
	"strings"
	"sync"
)

// trackBinding is a single bind for a Track
// Bind can be called multiple times, this stores the
// result for a single bind call so that it can be used when writing.
type trackBinding struct {
	id                          string
	ssrc, ssrcRTX, ssrcFEC      webrtc.SSRC
	payloadType, payloadTypeRTX webrtc.PayloadType
	writeStream                 webrtc.TrackLocalWriter
}

// TrackLocalStatic 是一个本地轨道，主要用来接收和转发RTP包，如果存在rid属性，则说明是一个simulcast流
// 详细调整参考文献：https://zhuanlan.zhihu.com/p/12087620711
type TrackLocalStatic struct {
	mu                sync.RWMutex
	bindings          []trackBinding
	codec             webrtc.RTPCodecCapability
	payloader         func(webrtc.RTPCodecCapability) (rtp.Payloader, error)
	id, rid, streamID string
	rtpTimestamp      *uint32
	sequenceNumber    uint16
	timestamp         uint32
}

// NewTrackLocalStatic returns a TrackLocalStatic.
func NewTrackLocalStatic(
	c webrtc.RTPCodecCapability,
	id, streamID string,
	options ...func(*TrackLocalStatic),
) (*TrackLocalStatic, error) {
	t := &TrackLocalStatic{
		codec:    c,
		bindings: []trackBinding{},
		id:       id,
		streamID: streamID,
	}

	for _, option := range options {
		option(t)
	}

	return t, nil
}

// WithRTPStreamID sets the RTP stream ID for this TrackLocalStatic.
func WithRTPStreamID(rid string) func(*TrackLocalStatic) {
	return func(t *TrackLocalStatic) {
		t.rid = rid
	}
}

// WithPayloader allows the user to override the Payloader.
func WithPayloader(h func(webrtc.RTPCodecCapability) (rtp.Payloader, error)) func(*TrackLocalStatic) {
	return func(s *TrackLocalStatic) {
		s.payloader = h
	}
}

// WithRTPTimestamp set the initial RTP timestamp for the track.
func WithRTPTimestamp(timestamp uint32) func(*TrackLocalStatic) {
	return func(s *TrackLocalStatic) {
		s.rtpTimestamp = &timestamp
	}
}

// Bind is called by the PeerConnection after negotiation is complete
// This asserts that the code requested is supported by the remote peer.
// If so it sets up all the state (SSRC and PayloadType) to have a call.
func (s *TrackLocalStatic) Bind(trackContext webrtc.TrackLocalContext) (webrtc.RTPCodecParameters, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	parameters := webrtc.RTPCodecParameters{RTPCodecCapability: s.codec}
	if codec, matchType := codecParametersFuzzySearch(
		parameters,
		trackContext.CodecParameters(),
	); matchType != codecMatchNone {
		s.bindings = append(s.bindings, trackBinding{
			ssrc:           trackContext.SSRC(),
			ssrcRTX:        trackContext.SSRCRetransmission(),
			ssrcFEC:        trackContext.SSRCForwardErrorCorrection(),
			payloadType:    codec.PayloadType,
			payloadTypeRTX: findRTXPayloadType(codec.PayloadType, trackContext.CodecParameters()),
			writeStream:    trackContext.WriteStream(),
			id:             trackContext.ID(),
		})

		return codec, nil
	}

	return webrtc.RTPCodecParameters{}, webrtc.ErrUnsupportedCodec
}

// Unbind implements the teardown logic when the track is no longer needed. This happens
// because a track has been stopped.
func (s *TrackLocalStatic) Unbind(t webrtc.TrackLocalContext) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.bindings {
		if s.bindings[i].id == t.ID() {
			s.bindings[i] = s.bindings[len(s.bindings)-1]
			s.bindings = s.bindings[:len(s.bindings)-1]

			return nil
		}
	}

	return webrtc.ErrUnbindFailed
}

// ID is the unique identifier for this Track. This should be unique for the
// stream, but doesn't have to globally unique. A common example would be 'audio' or 'video'
// and StreamID would be 'desktop' or 'webcam'.
func (s *TrackLocalStatic) ID() string { return s.id }

// StreamID is the group this track belongs too. This must be unique.
func (s *TrackLocalStatic) StreamID() string { return s.streamID }

// RID is the RTP stream identifier.
func (s *TrackLocalStatic) RID() string { return s.rid }

// Kind controls if this TrackLocal is audio or video.
func (s *TrackLocalStatic) Kind() webrtc.RTPCodecType {
	switch {
	case strings.HasPrefix(s.codec.MimeType, "audio/"):
		return webrtc.RTPCodecTypeAudio
	case strings.HasPrefix(s.codec.MimeType, "video/"):
		return webrtc.RTPCodecTypeVideo
	default:
		return webrtc.RTPCodecTypeUnknown
	}
}

// Codec gets the Codec of the track.
func (s *TrackLocalStatic) Codec() webrtc.RTPCodecCapability {
	return s.codec
}

// packetPool is a pool of packets used by WriteRTP and Write below
// nolint:gochecknoglobals
var rtpPacketPool = sync.Pool{
	New: func() any {
		return &rtp.Packet{}
	},
}

func resetPacketPoolAllocation(localPacket *rtp.Packet) {
	*localPacket = rtp.Packet{}
	rtpPacketPool.Put(localPacket)
}

func getPacketAllocationFromPool() *rtp.Packet {
	ipacket := rtpPacketPool.Get()

	return ipacket.(*rtp.Packet) //nolint:forcetypeassert
}

// WriteRTP writes a RTP Packet to the TrackLocalStatic
// If one PeerConnection fails the packets will still be sent to
// all PeerConnections. The error message will contain the ID of the failed
// PeerConnections so you can remove them.
func (s *TrackLocalStatic) WriteRTP(p *rtp.Packet) error {
	packet := getPacketAllocationFromPool()

	defer resetPacketPoolAllocation(packet)

	*packet = *p

	return s.writeRTP(packet)
}

func (s *TrackLocalStatic) SetRid(rid string) {
	s.rid = rid
}

func (s *TrackLocalStatic) SetSequenceNumber(seq uint16) {
	s.sequenceNumber = seq
}

func (s *TrackLocalStatic) SetTimestamp(timestamp uint32) {
	s.timestamp = timestamp
}

// writeRTP is like WriteRTP, except that it may modify the packet p.
func (s *TrackLocalStatic) writeRTP(packet *rtp.Packet) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	writeErrs := []error{}
	//==========================================
	/* 此部分内容是我们修改了官方逻辑的地方，以支持simulcast
	主要参考文献：https://zhuanlan.zhihu.com/p/12087620711
	*/
	if len(s.rid) > 0 { //如果是有大小流标记 我们才进行处理，这里处理前提条件是，写入的rtp包时间为时间增量，而不是绝对时间！！
		packet.SequenceNumber = s.sequenceNumber
		s.timestamp += packet.Timestamp
		packet.Timestamp = s.timestamp
		s.sequenceNumber++ //超过65535后，再累加会变成0开始，我们不需要特殊处理
	}
	//=========================================
	for _, b := range s.bindings {
		packet.Header.SSRC = uint32(b.ssrc)
		packet.Header.PayloadType = uint8(b.payloadType)
		// b.writeStream.WriteRTP below expects header and payload separately, so value of Packet.PaddingSize
		// would be lost. Copy it to Packet.Header.PaddingSize to avoid that problem.
		if packet.PaddingSize != 0 && packet.Header.PaddingSize == 0 {
			packet.Header.PaddingSize = packet.PaddingSize
		}
		if _, err := b.writeStream.WriteRTP(&packet.Header, packet.Payload); err != nil {
			writeErrs = append(writeErrs, err)
		}
	}

	return FlattenErrs(writeErrs)
}

// Write writes a RTP Packet as a buffer to the TrackLocalStatic
// If one PeerConnection fails the packets will still be sent to
// all PeerConnections. The error message will contain the ID of the failed
// PeerConnections so you can remove them.
func (s *TrackLocalStatic) Write(b []byte) (n int, err error) {
	packet := getPacketAllocationFromPool()

	defer resetPacketPoolAllocation(packet)

	if err = packet.Unmarshal(b); err != nil {
		return 0, err
	}

	return len(b), s.writeRTP(packet)
}
