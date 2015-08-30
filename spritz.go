// Package spritz implements the Spritz hash function and Spritz-xor
// stream cipher, as specified in the paper
//
//      Spritz—a spongy RC4-like stream cipher and hash function.
//      Ronald L. Rivest and Jacob C. N. Schuldt
//      https://people.csail.mit.edu/rivest/pubs/RS14.pdf
//
// Spritz is an evolution of RC4 and, like RC4, is rather slow.
// Standard ciphers like AES are faster and have been more carefully analyzed.
// Please consider this before using Spritz.
package spritz

const size = 256

// Sponge implements the Spritz sponge-like function.
type Sponge struct {
	s       [size]uint8
	a       int   // number of bytes absorbed
	i, j, k uint8 // state pointers
	w       uint8
	z       uint8 // last output
}

// NewSponge returns a new Sponge.
func NewSponge() *Sponge {
	q := new(Sponge)
	q.Reset()
	return q
}

// Reset sets the Sponge to its initial state.
func (q *Sponge) Reset() {
	q.i = 0
	q.j = 0
	q.k = 0
	q.z = 0
	q.a = 0
	q.w = 1
	for i := range q.s {
		q.s[i] = byte(i)
	}
}

func (q *Sponge) swap(i, j int) {
	q.s[i], q.s[j] = q.s[j], q.s[i]
}

// Write adds data to the Sponge state.
// It never returns an error.
//
// This is the Absorb operation from the Spritz paper.
func (q *Sponge) Write(b []byte) (n int, err error) {
	for _, v := range b {
		if q.a == size/2 {
			q.shuffle()
		}
		q.swap(q.a, size/2+int(v%16))
		q.a++

		if q.a == size/2 {
			q.shuffle()
		}
		q.swap(q.a, size/2+int(v/16))
		q.a++
	}
	return len(b), nil
}

// WriteStop writes a "stop symbol" to the Sponge state.
//
// This is the AbsorbStop operation from the Spritz paper.
func (q *Sponge) WriteStop() {
	if q.a == size/2 {
		q.shuffle()
	}
	q.a++
}

func (q *Sponge) shuffle() {
	q.whip()
	q.crush()
	q.whip()
	q.crush()
	q.whip()
	q.a = 0
}

func (q *Sponge) whip() {
	i := q.i
	j := q.j
	k := q.k
	w := q.w
	s := &q.s
	for r := 0; r < size*2; r++ {
		i += w
		j = k + s[j+s[i]]
		k = i + k + s[j]
		s[i], s[j] = s[j], s[i]
	}
	q.i = i
	q.j = j
	q.k = k
	q.w += 2
}

func (q *Sponge) crush() {
	for i := 0; i < size/2; i++ {
		// TODO: make constant-time
		if q.s[i] > q.s[size-1-i] {
			q.swap(i, size-1-i)
		}
	}
}

// Read fills b with pseudorandom bytes.
// It never returns an error.
//
// This is the Squeeze operation from the Spritz paper.
func (q *Sponge) Read(b []byte) (n int, err error) {
	if q.a > 0 {
		q.shuffle()
	}
	i := q.i
	j := q.j
	k := q.k
	w := q.w
	z := q.z
	s := q.s
	for ii := range b {
		i += w
		j = k + s[j+s[i]]
		k = i + k + s[j]
		s[i], s[j] = s[j], s[i]
		z = s[j+s[i+s[z+k]]]
		b[ii] = z
	}
	q.i = i
	q.j = j
	q.k = k
	q.z = z
	return len(b), nil
}

// XORKeyStream implements cipher.Stream:
// it copies src to dst, XORing each byte with a byte of the key stream.
//
// Src and dst may point at the same memory, but may not otherwise overlap.
//
// XORKeyStream panics if dst is shorter than src.
func (q *Sponge) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic("spritz: destination buffer too small")
	}
	if q.a > 0 {
		q.shuffle()
	}
	i := q.i
	j := q.j
	k := q.k
	w := q.w
	z := q.z
	s := q.s
	for ii, v := range src {
		i += w
		j = k + s[j+s[i]]
		k = i + k + s[j]
		s[i], s[j] = s[j], s[i]
		z = s[j+s[i+s[z+k]]]
		dst[ii] = v ^ z
	}
	q.i = i
	q.j = j
	q.k = k
	q.z = z
}
