package codec

import "io"

type Header struct {
	ServiceMethod string // 类似 "Service.Method"
	Seq           uint64
	Error         string
}

// 对信息体进行编解码的接口Codec, 编码解码器
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gab"
	JsonType Type = "application/json"
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
