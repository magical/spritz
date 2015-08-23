package spritz

import (
	"bytes"
	"testing"
)

func TestSpritz(t *testing.T) {
	var tests = []struct {
		in   string
		want []byte
	}{
		{"ABC", []byte{0x77, 0x9a, 0x8e, 0x01, 0xf9, 0xe9, 0xcb, 0xc0}},
		{"spam", []byte{0xf0, 0x60, 0x9a, 0x1d, 0xf1, 0x43, 0xce, 0xbf}},
		{"arcfour", []byte{0x1a, 0xfa, 0x8b, 0x5e, 0xe3, 0x37, 0xdb, 0xc7}},
	}
	for _, tt := range tests {
		d := New()
		d.absorb([]byte(tt.in))
		var got [8]byte
		d.squeeze(got[:])
		if !bytes.Equal(got[:], tt.want) {
			t.Errorf("Spritz(%q) = % x, want % x", tt.in, got[:], tt.want)
		}
	}
}

func BenchmarkShuffle(b *testing.B) {
	d := New()
	for i := 0; i < b.N; i++ {
		d.shuffle()
	}
}

func BenchmarkAbsorb1kB(b *testing.B) {
	var in [1 << 10]byte
	d := New()
	b.SetBytes(int64(len(in)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Reset()
		d.absorb(in[:])
	}
}

func BenchmarkSqueeze1kB(b *testing.B) {
	var out [1 << 10]byte
	d := New()
	b.SetBytes(int64(len(out)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Reset()
		d.squeeze(out[:])
	}
}
