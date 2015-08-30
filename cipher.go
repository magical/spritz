package spritz

// Cipher implements cipher.Stream
type Cipher struct {
	q Sponge
}

// NewCipher returns a new Cipher with the given key.
func NewCipher(key []byte) *Cipher {
	c := new(Cipher)
	c.q.Reset()
	c.q.Write(key)
	return c
}

func (c *Cipher) XORKeyStream(dst, src []byte) {
	c.q.XORKeyStream(dst, src)
}
