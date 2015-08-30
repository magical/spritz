package spritz

// Hash implements hash.Hash
type Hash struct {
	Sponge     *Sponge // state of the underlying sponge
	DigestSize int     // number of bytes returned by Sum
}

// NewHash returns a new Hash with the specified digest size.
func NewHash(size int) *Hash {
	return &Hash{
		Sponge:     NewSponge(),
		DigestSize: size,
	}
}

func (h *Hash) Reset() {
	h.Sponge.Reset()
}

func (h *Hash) Size() int {
	return h.DigestSize
}

func (h *Hash) BlockSize() int {
	return size
}

// Write adds more data to the hash. It never returns an error.
func (h *Hash) Write(p []byte) (n int, err error) {
	h.Sponge.Write(p)
	return len(p), nil
}

// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
func (h *Hash) Sum(b []byte) []byte {
	if h.DigestSize < 0 {
		return b
	}
	q := *h.Sponge // make a copy
	q.WriteStop()
	q.absorbUint(uint(h.DigestSize))
	i := len(b)
	if i+h.DigestSize > cap(b) {
		b = append(b[:i], make([]byte, h.DigestSize)...)
	}
	q.Read(b[i : i+h.DigestSize])
	return b[0 : i+h.DigestSize]
}

func (q *Sponge) absorbUint(v uint) {
	// 2.3: For definiteness, assume that r is represented as a base-N integer, high-order byte first, with no leading zeros.
	var b [8]byte
	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)
	n := 0
	for n < 8 && b[n] == 0 {
		n++
	}
	q.Write(b[n:])
}
