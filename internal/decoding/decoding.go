package decoding

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"

	"github.com/andybalholm/brotli"

	"github.com/LexusEgorov/goMetrics/internal/dohsimpson"
)

type decoding struct{}
type algo int

const (
	gz algo = iota
	deflate
	br
)

type closeableReader interface {
	io.Reader
	Close() error
}

func (d decoding) decode(data []byte, algo algo) ([]byte, *dohsimpson.Error) {
	var r closeableReader
	var err error

	switch algo {
	case gz:
		r, err = gzip.NewReader(bytes.NewReader(data))
	case deflate:
		r = flate.NewReader(bytes.NewReader(data))
	default:
		return data, nil
	}

	if err != nil {
		return nil, dohsimpson.NewDoh(400, err.Error())
	}

	defer r.Close()

	var b bytes.Buffer

	_, err = b.ReadFrom(r)

	if err != nil {
		return nil, dohsimpson.NewDoh(400, err.Error())
	}

	return b.Bytes(), nil
}

func (d decoding) DecodeGz(data []byte) ([]byte, *dohsimpson.Error) {
	return d.decode(data, gz)
}

func (d decoding) DecodeDeflate(data []byte) ([]byte, *dohsimpson.Error) {
	return d.decode(data, deflate)
}

func (d decoding) DecodeBr(data []byte) ([]byte, *dohsimpson.Error) {
	r := brotli.NewReader(bytes.NewReader(data))
	var b bytes.Buffer

	_, err := b.ReadFrom(r)

	if err != nil {
		return nil, dohsimpson.NewDoh(400, err.Error())
	}

	return b.Bytes(), nil
}

func NewDecoding() decoding {
	return decoding{}
}
