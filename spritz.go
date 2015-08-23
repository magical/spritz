package spritz

const size = 256

type digest struct {
    i, j, k, w uint8
    z byte // last output byte
    a int // number of bytes absorbed
    s [size]byte
}

func New() *digest {
    d := new(digest)
    d.Reset()
    return d
}

func (d *digest) Reset() {
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
        d.swap(d.a, size/2 + int(v%16))
        d.a += 1

        if d.a == size/2 { 
            d.shuffle()
        }
        d.swap(d.a, size/2 + int(v/16))
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
    for i := 0; i < size*2; i++ {
        d.update()
    }
    d.w += 2
}

func (d *digest) crush() {
    for i := 0; i < size/2; i++ {
        // TODO: make constant-time
        if d.s[i] > d.s[size - 1 - i] {
            d.swap(i, size - 1 - i)
        }
    }
}

func (d *digest) squeeze(b []byte) {
    if d.a > 0 {
        d.shuffle()
    }
    for i := range b {
        b[i] = d.drip()
    }
}

func (d *digest) drip() byte {
    if d.a > 0 {
        d.shuffle()
    }
    d.update()
    return d.output()
}

func (d *digest) update() {
    d.i += d.w
    d.j = d.k + d.s[d.j + d.s[d.i]]
    d.k = d.i + d.k + d.s[d.j]
    d.swap(int(d.i), int(d.j))
}

func (d *digest) output() byte {
    d.z = d.s[d.j + d.s[d.i + d.s[d.z + d.k]]]
    return d.z
}
