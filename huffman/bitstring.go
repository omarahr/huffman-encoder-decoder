package huffman

type BitString struct {
	bits []uint8
	size int
}

func NewBitString() *BitString {
	return &BitString{
		bits: make([]uint8, 0),
		size: 0,
	}
}

func (b *BitString) AddBit(bit bool) {
	if b.size%8 == 0 {
		b.bits = append(b.bits, 0)
	}

	// TODO won't handle large number might overflow
	// that's why shifting happens occasionally with the correct max size
	if bit {
		b.bits[b.size/8] |= 1 << (7 - uint(b.size%8))
	}

	b.size++
}

func (b *BitString) AddBits(bits []bool) {
	for _, bit := range bits {
		b.AddBit(bit)
	}
}

func (b *BitString) GetBytes() []byte {
	return b.bits
}

func (b *BitString) GetTrailingSize() int {
	if b.size == 0 {
		return 0
	}

	reminder := b.size % 8
	if reminder == 0 {
		return 8
	}

	return reminder
}

func (b *BitString) Size() int {
	return b.size
}

func (b *BitString) GetReadyBytes() []byte {
	readyBytesSize := b.size / 8
	readyBytes := append([]byte{}, b.bits[:readyBytesSize]...)
	b.size = b.size % 8
	b.bits = append([]byte{}, b.bits[readyBytesSize:]...)
	return readyBytes
}
