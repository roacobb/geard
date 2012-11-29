// Copyright (c) 2012 Graeme Connell. All rights reserved.
// Copyright (c) 2009-2012 Andreas Krennmair. All rights reserved.

package gopacket

import (
	"errors"
)

type decodeResult struct {
	// An error encountered in this decode call.  If this is set, everything else
	// will be ignored.
	err error
	// The layer we've created with this decode call
	layer Layer
	// The next decoder to call
	next decoder
	// The bytes that are left to be decoded
	left []byte
}

// decoder decodes the next layer in a packet.  It returns a set of useful
// information, which is used by the packet decoding logic to update packet
// state.  Optionally, the decode function may set any of the specificLayer
// pointers to point to the new layer it has created.
//
// This decoder interface is the internal interface used by gopacket to store
// the next method to use for decoding the rest of the data available in the
// packet.  It should exhibit the following behavior:
// * if there's an error, set decodeResult.err.  All other fields will be
//   ignored and a DecodeError layer will be created with that error.
// * if there's NOT an error, set layer to the layer created by this decoder,
//   next to the next decoder to run, and left to the bytes not yet processed.
//   if either decoder is nil or left is empty, this packet's decoding is
//   considered complete and nothing else is done.
//
// If the decoded layer is one of the specific layers in specificLayers, the
// function should set specificLayers' pointer to the new layer.  For example,
// note how decodeIp4 sets specificLayers' network pointer to the newly created
// IPv4 layer object.
type decoder interface {
	decode([]byte, *specificLayers) decodeResult
}

// decoderFunc is an implementation of decoder that's a simple function.
type decoderFunc func([]byte, *specificLayers) decodeResult

func (d decoderFunc) decode(data []byte, s *specificLayers) decodeResult {
	// function, call thyself.
	return d(data, s)
}

// DecodeMethod tells gopacket how to decode a packet.
type DecodeMethod bool

const (
	// Lazy decoding decodes the minimum number of layers needed to return data
	// for a packet at each function call.  Be careful using this with concurrent
	// packet processors, as each call to packet.* could mutate the packet, and
	// two concurrent function calls could interact poorly.
	Lazy DecodeMethod = true
	// Eager decoding decodes all layers of a packet immediately.  Slower than
	// lazy decoding, but better if the packet is expected to be used concurrently
	// at a later date, since after an eager Decode, the packet is guaranteed to
	// not mutate itself on packet.* function calls.
	Eager DecodeMethod = false
)

// PacketDecoder provides the functionality to decode a set of bytes into a
// packet, and decode that packet into one or more layers.
type PacketDecoder interface {
	Decode(data []byte, method DecodeMethod) Packet
}

// DecodeFailure is a packet layer created if decoding of the packet data failed
// for some reason.  It implements ErrorLayer.
type DecodeFailure struct {
	data []byte
	err  error
}

// Returns the entire payload which failed to be decoded.
func (d *DecodeFailure) Payload() []byte { return d.data }

// Returns the error encountered during decoding.
func (d *DecodeFailure) Error() error { return d.err }

// Returns TYPE_DECODE_FAILURE
func (d *DecodeFailure) LayerType() LayerType { return TYPE_DECODE_FAILURE }

// decodeUnknown "decodes" unsupported data types by returning an error.
// This decoder will thus always return a DecodeFailure layer.
var decodeUnknown decoderFunc = func(data []byte, _ *specificLayers) (out decodeResult) {
	out.err = errors.New("Link type not currently supported")
	return
}

// decodePayload decodes data by returning it all in a Payload layer.
var decodePayload decoderFunc = func(data []byte, s *specificLayers) (out decodeResult) {
	payload := &Payload{Data: data}
	out.layer = payload
	s.application = payload
	return
}