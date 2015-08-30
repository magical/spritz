package spritz

import (
	"bytes"
	"testing"
)

var spongeTests = []struct {
	in   string
	want []byte
}{
	// Vectors from the Spritz paper
	{"ABC", []byte{0x77, 0x9a, 0x8e, 0x01, 0xf9, 0xe9, 0xcb, 0xc0}},
	{"spam", []byte{0xf0, 0x60, 0x9a, 0x1d, 0xf1, 0x43, 0xce, 0xbf}},
	{"arcfour", []byte{0x1a, 0xfa, 0x8b, 0x5e, 0xe3, 0x37, 0xdb, 0xc7}},
}

func TestSpritz(t *testing.T) {
	q := NewSponge()
	for _, tt := range spongeTests {
		q.Reset()
		q.Write([]byte(tt.in))
		var got [8]byte
		q.Read(got[:])
		if !bytes.Equal(got[:], tt.want) {
			t.Errorf("Spritz(%q) = % x, want % x", tt.in, got[:], tt.want)
		}
	}
}

func TestXORKeystream(t *testing.T) {
	q := NewSponge()
	for _, tt := range spongeTests {
		q.Reset()
		q.Write([]byte(tt.in))
		var got [8]byte
		// xor keystream with the zero buffer
		// should be equivalent to Read
		q.XORKeyStream(got[:], got[:])
		if !bytes.Equal(got[:], tt.want) {
			t.Errorf("Spritz(%q) = % x, want % x", tt.in, got[:], tt.want)
		}
	}
}

func TestSpritz_MultipleWrites(t *testing.T) {
	// Check that writing one byte at a time is equivalent to writing many bytes at once
	q := NewSponge()
	q.Write([]byte{'A'})
	q.Write([]byte{'B'})
	q.Write([]byte{'C'})
	var got [8]byte
	q.Read(got[:])
	want := []byte{0x77, 0x9a, 0x8e, 0x01, 0xf9, 0xe9, 0xcb, 0xc0}
	if !bytes.Equal(got[:], want) {
		t.Errorf("got % x, want % x", got[:], want)
	}
}

func TestSpritz_MultipleReads(t *testing.T) {
	// Check that reading one byte at a time is equivalent to reading many bytes at once
	q := NewSponge()
	q.Write([]byte("ABC"))
	var got [8]byte
	for i := range got {
		q.Read(got[i : i+1])
	}
	want := []byte{0x77, 0x9a, 0x8e, 0x01, 0xf9, 0xe9, 0xcb, 0xc0}
	if !bytes.Equal(got[:], want) {
		t.Errorf("got % x, want % x", got, want)
	}
}

func TestSpritz_MultipleXORs(t *testing.T) {
	// Check that encrypting one byte at a time is equivalent to encrypting many bytes at once
	q := NewSponge()
	q.Write([]byte("ABC"))
	var got [8]byte
	for i := range got {
		q.XORKeyStream(got[i:i+1], got[i:i+1])
	}
	want := []byte{0x77, 0x9a, 0x8e, 0x01, 0xf9, 0xe9, 0xcb, 0xc0}
	if !bytes.Equal(got[:], want) {
		t.Errorf("got % x, want % x", got, want)
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
