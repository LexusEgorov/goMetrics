package encoding

type encoding struct{}

func (e encoding) EncodeGz() {}

func (e encoding) DecodeGz() {}

func (e encoding) EncodeCompress() {}

func (e encoding) DecodeCompress() {}

func (e encoding) EncodeDeflate() {}

func (e encoding) DecodeDeflate() {}

func (e encoding) EncodeBr() {}

func (e encoding) DecodeBr() {}

func NewEncoding() encoding {
	return encoding{}
}
