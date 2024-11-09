package encoding

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"

	"github.com/andybalholm/brotli"

	"github.com/LexusEgorov/goMetrics/internal/dohsimpson"
)

type encoding struct{}
type algo int

const (
	gz algo = iota
	deflate
	br
)

type closeableWriter interface {
	io.Writer
	Close() error
}

func (e encoding) encode(data []byte, algo algo) ([]byte, *dohsimpson.Error) {
	var b bytes.Buffer

	var w closeableWriter
	var err error

	switch algo {
	case gz:
		w, err = gzip.NewWriterLevel(&b, gzip.BestCompression)
	case deflate:
		w, err = flate.NewWriter(&b, flate.BestCompression)
	case br:
		w = brotli.NewWriterLevel(&b, brotli.BestCompression)
	default:
		return data, nil
	}

	if err != nil {
		return nil, dohsimpson.NewDoh(500, err.Error())
	}

	_, err = w.Write(data)

	if err != nil {
		return nil, dohsimpson.NewDoh(500, err.Error())
	}

	err = w.Close()

	if err != nil {
		return nil, dohsimpson.NewDoh(500, err.Error())
	}

	return b.Bytes(), nil
}

func (e encoding) EncodeGz(data []byte) ([]byte, *dohsimpson.Error) {
	return e.encode(data, gz)
}

func (e encoding) EncodeDeflate(data []byte) ([]byte, *dohsimpson.Error) {
	return e.encode(data, deflate)
}

func (e encoding) EncodeBr(data []byte) ([]byte, *dohsimpson.Error) {
	return e.encode(data, br)
}

func NewEncoding() encoding {
	return encoding{}
}
