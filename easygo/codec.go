package easygo

import (
	"encoding/binary"
	"fmt"
)

const HEADER_SIZE = 4

//=======================================================

type IEncoder interface {
	Encode(packet []byte) []byte
}
type Encoder struct {
	Me IEncoder
}

func NewEncoder() *Encoder {
	p := &Encoder{}
	p.Init(p)
	return p
}

func (self *Encoder) Init(me IEncoder) {
	self.Me = me
}

func (self *Encoder) Encode(packet []byte) []byte {
	size := len(packet)
	bytes := make([]byte, HEADER_SIZE+size)
	binary.BigEndian.PutUint32(bytes, uint32(size))

	dst, src := bytes[HEADER_SIZE:], packet
	copy(dst, src)
	return bytes
}

//=======================================================
// 解码器
type IDecoder interface {
	Reset()
	Decode(buffer IBuffer) ([][]byte, string)
}
type Decoder struct {
	Me                 IDecoder
	PacketLen          int
	MaxInputBufferSize int
}

func NewDecoder(maxInputBufferSize ...int) *Decoder {
	size := append(maxInputBufferSize, 3*1024*1024)[0]
	p := &Decoder{}
	p.Init(p, size)
	return p
}

func (self *Decoder) Init(me IDecoder, MaxInputBufferSize ...int) {
	self.Me = me
	self.PacketLen = 0
	self.MaxInputBufferSize = append(MaxInputBufferSize, 3*1024*1024)[0]
}

func (self *Decoder) Reset() {
	self.PacketLen = 0
}

func (self *Decoder) Decode(buffer IBuffer) ([][]byte, string) {
	var packets [][]byte
	header_size := HEADER_SIZE

	for {
		if self.PacketLen == 0 {
			if buffer.ReadableBytes() < header_size { // 连一个头都不够
				return packets, "" // 跳出,需要再读多一点数据
			}

			size := buffer.PeekRead(header_size)
			self.PacketLen = int(binary.BigEndian.Uint32(size))

			buffer.Retrieve(header_size)
			if self.PacketLen <= 0 {
				s := fmt.Sprintf("包大小是 %d,不得小于等于 0", self.PacketLen)
				panic(s)
			}
			if int(self.PacketLen) > self.MaxInputBufferSize { // 防止恶意攻击,恶意的客户端
				s := fmt.Sprintf("接收到 4 byte 的包头，指示包体大小是 %d,超过了 %d", self.PacketLen, self.MaxInputBufferSize)
				// panic(s)
				return packets, s
			}
		}

		if buffer.ReadableBytes() < self.PacketLen { // 不够一个逻辑包
			return packets, "" // 跳出,需要再读多一点数据
		}

		packet := buffer.PeekRead(self.PacketLen)
		buffer.Retrieve(self.PacketLen)
		self.PacketLen = 0
		packets = append(packets, packet)

	} // for
}
