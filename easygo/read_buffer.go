package easygo

type IBuffer interface {
	ReadableBytes() int
	MakeSpace(needSize int)
	HasWritten(length int)
	Retrieve(length int)
	PeekRead(size int) []byte
	PeekWrite(size ...int) []byte
}

type Buffer struct {
	Me                IBuffer
	ByteSlice         []byte
	ReadIdx, WriteIdx int
}

func NewBuffer(initLen ...int) *Buffer {
	len := append(initLen, 1024)[0]
	p := &Buffer{}
	p.Init(p, len)
	return p
}

func (self *Buffer) Init(me IBuffer, initLen ...int) {
	self.Me = me
	len := append(initLen, 1024)[0]
	if len <= 0 {
		panic("滚,长度至少为1")
	}
	self.ByteSlice = make([]byte, len)
	self.ReadIdx, self.WriteIdx = 0, 0 // [前闭,后开)
}

func (self *Buffer) ReadableBytes() int {
	return self.WriteIdx - self.ReadIdx
}

func (self *Buffer) MakeSpace(needSize int) { // 挪动或扩充,得到一个连续的 needSize 可写空间.(只会扩充得比需要的多,不会少)
	readable := self.WriteIdx - self.ReadIdx
	length := len(self.ByteSlice)

	for length-readable < needSize || length-readable < length/3.0 { // 可写空间不够 needSize 或不足 1/3,要搞大他
		length *= 2
	}

	if length != len(self.ByteSlice) {
		newSlice := make([]byte, length) // 另搞一个对象扩充空间
		dst := newSlice[0:readable]
		src := self.ByteSlice[self.ReadIdx : self.ReadIdx+readable]
		copy(dst, src)
		self.ByteSlice = newSlice
	} else {
		dst := self.ByteSlice[0:readable]
		src := self.ByteSlice[self.ReadIdx : self.ReadIdx+readable]
		copy(dst, src)
	}
	self.ReadIdx, self.WriteIdx = 0, readable
}

func (self *Buffer) HasWritten(length int) {
	if length > len(self.ByteSlice)-self.WriteIdx {
		panic("都没有这么多字节数")
	}
	self.WriteIdx += length
}

func (self *Buffer) Retrieve(length int) {
	if length > self.WriteIdx-self.ReadIdx {
		panic("没有这么多字节数")
	}
	if length < self.WriteIdx-self.ReadIdx {
		self.ReadIdx += length
	} else {
		self.ReadIdx, self.WriteIdx = 0, 0
	}
}

func (self *Buffer) PeekWrite(sizes ...int) []byte { // 0 表示拿全部
	size := append(sizes, 0)[0]
	if size == 0 {
		if self.WriteIdx == len(self.ByteSlice) { // 没有空间了
			self.MakeSpace(1024)
		}
		return self.ByteSlice[self.WriteIdx:]
	} else {
		if self.WriteIdx+size > len(self.ByteSlice) { // 空间不足
			self.MakeSpace(size)
		}
		return self.ByteSlice[self.WriteIdx : self.WriteIdx+size]
	}
}

func (self *Buffer) PeekRead(size int) []byte {
	return self.ByteSlice[self.ReadIdx : self.ReadIdx+size]
}
