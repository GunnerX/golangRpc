package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

// gob的编码解码器
type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder // 编码器	写操作
	enc  *gob.Encoder // 解码器	读操作
}

var _ Codec = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn), // 解码器 从conn中读取解码信息
		enc:  gob.NewEncoder(buf),  // 编码器 编码信息并写入conn(通过缓冲buf)
	}
}

// GobCodec 实现 Codec接口的4个方法

// 读取并解码header
func (c *GobCodec) ReadHeader(header *Header) error {
	return c.dec.Decode(header)
}

// 读取并解码body
func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

// 编码并写入header body
func (c *GobCodec) Write(header *Header, body interface{}) (err error) {
	defer func() {
		c.buf.Flush()
		if err != nil {
			c.Close()
		}
	}()

	if err = c.enc.Encode(header); err != nil {
		log.Println("rpc codec: gob error encoding header: ", err)
		return err
	}
	if err := c.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body: ", err)
		return err
	}
	return nil
}

// 关闭连接
func (c *GobCodec) Close() error {
	return c.conn.Close()
}
