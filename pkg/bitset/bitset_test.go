package bitset

import (
	"testing"
)

func TestBitset32(t *testing.T) {
	bs32 := NewBitset32()

	if bs32.IsSet(0) {
		t.Error("Expected bit 0 to be unset")
	}

	bs32.Set(0)
	if !bs32.IsSet(0) {
		t.Error("Expected bit 0 to be set")
	}

	bs32.Set(1)
	if !bs32.IsSet(1) {
		t.Error("Expected bit 1 to be set")
	}

	bs32.Clear(0)
	if bs32.IsSet(0) {
		t.Error("Expected bit 0 to be unset after Clear")
	}

	bs32.Reset()
	for i := 0; i < 32; i++ {
		if bs32.IsSet(i) {
			t.Errorf("Expected bit %d to be unset after Reset", i)
		}
	}

	expectPanic(t, func() {
		bs32.Set(32)
	})
}

func TestBitset64(t *testing.T) {
	bs64 := NewBitset64()

	if bs64.IsSet(0) {
		t.Error("Expected bit 0 to be unset")
	}

	bs64.Set(0)
	if !bs64.IsSet(0) {
		t.Error("Expected bit 0 to be set")
	}

	bs64.Set(1)
	if !bs64.IsSet(1) {
		t.Error("Expected bit 1 to be set")
	}

	bs64.Clear(0)
	if bs64.IsSet(0) {
		t.Error("Expected bit 0 to be unset after Clear")
	}

	bs64.Reset()
	for i := 0; i < 64; i++ {
		if bs64.IsSet(i) {
			t.Errorf("Expected bit %d to be unset after Reset", i)
		}
	}

	expectPanic(t, func() {
		bs64.Set(64)
	})
}

func TestStringer(t *testing.T) {
	bs32 := NewBitset32()
	bs32.Set(0)
	bs32.Set(1)
	expectedString32 := "00000000000000000000000000000011"
	if bs32.String() != expectedString32 {
		t.Errorf("Expected %s, got %s", expectedString32, bs32.String())
	}

	bs64 := NewBitset64()
	bs64.Set(0)
	bs64.Set(63)
	expectedString64 := "1000000000000000000000000000000000000000000000000000000000000001"
	if bs64.String() != expectedString64 {
		t.Errorf("Expected %s, got %s", expectedString64, bs64.String())
	}
}

func expectPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic but did not occur")
		}
	}()
	f()
}
