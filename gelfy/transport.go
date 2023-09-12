package gelfy

import (
	"bytes"
	"compress/zlib"
	"crypto/rand"
	"encoding/json"
	"net"
)

// Options describes options for Transport.
// If serialzied message data size greater than CompressSize, message will be compress using zlib.
// If serialized message data size greater than ChunkSize, message will be splitted into multiply chunks.
// For more details see https://www.graylog.org/resources/gelf-2/
type Options struct {
	CompressSize int
	ChunkSize    int
}

// DefaultOptions describes default options set for transport.
var DefaultOptions = Options{
	CompressSize: 1000,
	ChunkSize:    8000,
}

// Transport describes UDP connection to graylog server.
type Transport struct {
	opts Options
	conn net.Conn
}

// NewTransport returns new UDP transport to grayloag server.
// Server ip and port must be specified in addr.
// Transport will use opts as options set.
// If opts.CompressSize == 0 or opts.ChunkSize == 0, values from DefaultOptions will be used.
func NewTransport(addr string, opts Options) (*Transport, error) {
	if opts.CompressSize == 0 {
		opts.CompressSize = DefaultOptions.CompressSize
	}

	if opts.ChunkSize == 0 {
		opts.ChunkSize = DefaultOptions.ChunkSize
	}

	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}

	return &Transport{opts: opts, conn: conn}, nil
}

// Send send GELF message to server.
func (t *Transport) Send(m *Message) error {
	data, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return t.SendRaw(data)
}

// SendRaw send serialized GELF message data to server.
func (t *Transport) SendRaw(data []byte) error {
	data = t.tryCompress(data)
	return t.send(data)
}

func (t *Transport) tryCompress(data []byte) []byte {
	if len(data) < t.opts.CompressSize {
		return data
	}

	buf := bytes.Buffer{}
	zw := zlib.NewWriter(&buf)
	if _, err := zw.Write(data); err != nil {
		return data
	}
	zw.Close()

	if buf.Len() < len(data) {
		return buf.Bytes()
	}

	return data
}

func (t *Transport) send(data []byte) error {
	if len(data) < t.opts.ChunkSize {
		return t.sendSingle(data)
	}
	return t.sendChunked(data)
}

func (t *Transport) sendSingle(data []byte) error {
	_, err := t.conn.Write(data)
	return err
}

var magicChunked = []byte{0x1e, 0x0f}

func (t *Transport) sendChunked(data []byte) error {
	id := make([]byte, 8)
	rand.Read(id)

	nChunks := uint8(len(data)/t.opts.ChunkSize + 1)

	in := bytes.NewBuffer(data)
	out := bytes.Buffer{}
	for i := uint8(0); i < nChunks; i++ {
		out.Reset()
		out.Write(magicChunked)
		out.Write(id)
		out.WriteByte(i)
		out.WriteByte(nChunks)
		out.Write(in.Next(t.opts.ChunkSize))

		_, err := t.conn.Write(out.Bytes())
		if err != nil {
			return err
		}
	}

	return nil
}
