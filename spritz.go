package spritz

const size = 256

type digest struct {
	i, j, k uint8
	w       uint8
	z       byte // last output byte
	a       int  // number of bytes absorbed
	s       [size]byte
}

func newDigest() *digest {
	d := new(digest)
	d.reset()
	return d
}

func (d *digest) reset() {
	d.i = 0
	d.j = 0
	d.k = 0
	d.z = 0
	d.a = 0
	d.w = 1
	for i := range d.s {
		d.s[i] = byte(i)
	}
}

func (d *digest) swap(i, j int) {
	d.s[i], d.s[j] = d.s[j], d.s[i]
}

func (d *digest) absorb(b []byte) {
	for _, v := range b {
		if d.a == size/2 {
			d.shuffle()
		}
		d.swap(d.a, size/2+int(v%16))
		d.a += 1

		if d.a == size/2 {
			d.shuffle()
		}
		d.swap(d.a, size/2+int(v/16))
		d.a += 1
	}
}

func (d *digest) stop() {
	if d.a == size/2 {
		d.shuffle()
	}
	d.a += 1
}

func (d *digest) shuffle() {
	d.whip()
	d.crush()
	d.whip()
	d.crush()
	d.whip()
	d.a = 0
}

func (d *digest) whip() {
	i := d.i
	j := d.j
	k := d.k
	w := d.w
	for r := 0; r < size*2; r++ {
		i += w
		j = k + d.s[j+d.s[i]]
		k = i + k + d.s[j]
		d.s[i], d.s[j] = d.s[j], d.s[i]
	}
	d.i = i
	d.j = j
	d.k = k
	d.w += 2
}

func (d *digest) crush() {
	for i := 0; i < size/2; i++ {
		// TODO: make constant-time
		if d.s[i] > d.s[size-1-i] {
			d.swap(i, size-1-i)
		}
	}
}

func (d *digest) squeeze(b []byte) {
	if d.a > 0 {
		d.shuffle()
	}
	i := d.i
	j := d.j
	k := d.k
	w := d.w
	z := d.z
	for ii := range b {
		i += w
		j = k + d.s[j+d.s[i]]
		k = i + k + d.s[j]
		d.s[i], d.s[j] = d.s[j], d.s[i]
		z = d.s[j+d.s[i+d.s[z+k]]]
		b[ii] = z
	}
	d.i = i
	d.j = j
	d.k = k
	d.z = z
}
