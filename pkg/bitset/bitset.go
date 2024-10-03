package bitset

import "fmt"

type Bitset interface {
	Set(bit int)
	Clear(bit int)
	IsSet(bit int) bool
	Reset()
	fmt.Stringer
}

type Bitset32 struct {
	value uint32
}

func NewBitset32() *Bitset32 {
	return &Bitset32{value: 0}
}

func (b *Bitset32) Get32() uint32 {
	return b.value
}

func (b *Bitset32) Set(bit int) {
	b.validate(bit)
	b.value |= (1 << bit)
}

func (b *Bitset32) Clear(bit int) {
	b.validate(bit)
	b.value &^= (1 << bit)
}

func (b *Bitset32) IsSet(bit int) bool {
	b.validate(bit)
	return b.value&(1<<bit) != 0
}

func (b *Bitset32) Reset() {
	b.value = 0
}

func (b *Bitset32) validate(bit int) {
	if bit < 0 || bit >= 32 {
		panic("bit index out of range")
	}
}

func (b *Bitset32) String() string {
	return fmt.Sprintf("%032b", b.value)
}

type Bitset64 struct {
	value uint64
}

func NewBitset64() *Bitset64 {
	return &Bitset64{
		value: 0,
	}
}

func (b *Bitset64) Get64() uint64 {
	return b.value
}

func (b *Bitset64) Set(bit int) {
	b.validate(bit)
	b.value |= (1 << bit)
}

func (b *Bitset64) Clear(bit int) {
	b.validate(bit)
	b.value &^= (1 << bit)
}

func (b *Bitset64) IsSet(bit int) bool {
	b.validate(bit)
	return b.value&(1<<bit) != 0
}

func (b *Bitset64) Reset() {
	b.value = 0
}

func (b *Bitset64) validate(bit int) {
	if bit < 0 || bit >= 64 {
		panic("bit index out of range")
	}
}

func (b *Bitset64) String() string {
	return fmt.Sprintf("%064b", b.value)
}
