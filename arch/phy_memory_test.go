package arch

import (
	"testing"
)

func TestPhyMemory(t *testing.T) {
	eo := func(cond bool, s string, args ...interface{}) {
		if cond {
			t.Fatalf(s, args...)
		}
	}

	as := func(cond bool) {
		if !cond {
			t.Fatal("assertion failed")
		}
	}

	size := uint32(20 * PageSize)
	m := newPhyMemory(size)
	eo(m.Size() != size, "size mismatch")

	_, e := m.ReadU32(0)
	eo(e == nil, "read on nullptr should return memory out of range")

	w, e := m.ReadU32(PageSize + 4)
	eo(e != nil, "get an error for word reading")
	eo(w != 0, "page is not zeroed out")

	_, e = m.ReadU32(PageSize + 13)
	eo(e != errMisalign, "should have misalign error")
	_, e = m.ReadU32(size - 1)
	eo(e != errMisalign, "should have misalign error")

	_, e = m.ReadU32(size)
	eo(e.Code != ErrOutOfRange, "should have out of range error")
	eo(e.Arg != size, "expect arg to be set")

	off := uint32(56 + PageSize*2)
	as(m.WriteU8(off+0, 0x37) == nil)
	as(m.WriteU8(off+1, 0x21) == nil)
	as(m.WriteU8(off+2, 0x5a) == nil)
	as(m.WriteU8(off+3, 0x70) == nil)
	exp := uint32(0x705a2137)
	w, e = m.ReadU32(off)
	as(e == nil)
	eo(w != exp, "expect 0x%08x got 0x%08x", exp, w)
	b, e := m.ReadU8(off + 2)
	as(e == nil)
	eo(b != 0x5a, "expect 0x5a got %02x", b)
}
