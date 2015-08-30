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
		// Vectors from the Spritz paper
		{"ABC", []byte{0x77, 0x9a, 0x8e, 0x01, 0xf9, 0xe9, 0xcb, 0xc0}},
		{"spam", []byte{0xf0, 0x60, 0x9a, 0x1d, 0xf1, 0x43, 0xce, 0xbf}},
		{"arcfour", []byte{0x1a, 0xfa, 0x8b, 0x5e, 0xe3, 0x37, 0xdb, 0xc7}},
	}
	q := NewSponge()
	for _, tt := range tests {
		q.Reset()
		q.Write([]byte(tt.in))
		var got [8]byte
		q.Read(got[:])
		if !bytes.Equal(got[:], tt.want) {
			t.Errorf("Spritz(%q) = % x, want % x", tt.in, got[:], tt.want)
		}
	}
}

func BenchmarkShuffle(b *testing.B) {
	q := NewSponge()
	for i := 0; i < b.N; i++ {
		q.shuffle()
	}
}

func BenchmarkAbsorb1kB(b *testing.B) {
	var in [1 << 10]byte
	q := NewSponge()
	b.SetBytes(int64(len(in)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Reset()
		q.Write(in[:])
	}
}

func BenchmarkSqueeze1kB(b *testing.B) {
	var out [1 << 10]byte
	q := NewSponge()
	b.SetBytes(int64(len(out)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Reset()
		q.Read(out[:])
	}
}
