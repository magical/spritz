package spritz

import (
	"bytes"
	"encoding/hex"
	"hash"
	"testing"
)

// Test that Hash implements hash.Hash
var _ hash.Hash = new(Hash)

func TestHash(t *testing.T) {
	var tests = []struct {
		in   string
		want []byte
	}{
		// vectors from the spritz paper
		{"ABC", []byte{0x02, 0x8f, 0xa2, 0xb4, 0x8b, 0x93, 0x4a, 0x18}},
		{"spam", []byte{0xac, 0xbb, 0xa0, 0x81, 0x3f, 0x30, 0x0d, 0x3a}},
		{"arcfour", []byte{0xff, 0x8c, 0xf2, 0x68, 0x09, 0x4c, 0x87, 0xb9}},
	}
	h := NewHash(32)
	for _, tt := range tests {
		h.Reset()
		h.Write([]byte(tt.in))
		got := h.Sum(nil)
		if !bytes.Equal(got[:8], tt.want) {
			t.Errorf("Hash(%q, len=32) = %x, want %x...", tt.in, got, tt.want)
		}
	}
}

func TestHashSizes(t *testing.T) {
	var tests = []struct {
		size int
		want string
	}{
		{-1, ""},
		{0, ""},
		{1, "70"},
		{2, "ebb7"},
		{3, "9e30de"},
		{4, "b13ae345"},
		{5, "2cefdbe378"},
		{6, "7a5514943e2e"},
		{7, "80c2f2579fcc3b"},
		{8, "847bc542d05705ef"},
		{9, "48fac5b0d5112f8e4b"},
		{10, "77a7b4e6b055d972f6b4"},
		{11, "b2820bba2faeee45626793"},
		{12, "b58f212e18581383fabe51f7"},
		{13, "b915245ab6343187ff651502ab"},
		{14, "555b1a914ad52463da7b30940591"},
		{15, "9165a7fb14d3c043fed6563c7ed5d6"},
		{16, "dee2b6e00fada570e614d81921289202"},
		{17, "2b154a4cf8e25682a1726763672e4fb415"},
		{18, "24674ba9c23dba08e4cb452b7a454373b16d"},
		{19, "59c2c6c6435052ebd512aff89007de0ab9837e"},
		{20, "5f3869492e591e0d4235b0e604ec2652fc746661"},
		{21, "90194fa9deb2083552176677c503969e1e695c8244"},
		{22, "a639f5c2d3303ebdfe4d9242bb697fa40c67f52fdbd4"},
		{23, "ff0a2f6a66dc031c5d5c163d04ecad5b5cfc1248688223"},
		{24, "178d819c0186bb338c8c735bfb531efa882c23efdf029f9a"},
		{25, "09ecfde6499ed36acd9e0619b9d2db5dc9763c7e81a4f23f7c"},
		{26, "642ac7224f079f8635fe17b81a65a956982fb80d1ab3046d41eb"},
		{27, "98dfb92d2dd4f69209a2a27e0e7ba6178d0407cf5ff0c37687041d"},
		{28, "b1d9dea220516168f5e62aa5c8741843cf0df6a3808887d8737a92b0"},
		{29, "fa0475aa5b23c1558aace751270e076efc827fb154dab072743daa7178"},
		{30, "f97731f36c8944c4cd798d87a200508c792f5be65410de934c5735316f19"},
		{31, "1e47c291bbea6dc071a87faf7325bc79c2f97654bc059a40973b056944cbe7"},
		{32, "eddbfc9e608c1a73eb8d1311c483626104b8ea762d3075768af586838ffb0381"},
		{256, "5a55bc67cc48a54473e04eb2baf184872da6d25652cfa611bf399727a87ce71575e7e9d68a218df0cd273cb4c1464e363f35c94d58b015d1f0566fed6edb508c2b683b63ee54ff85dea1bcf7f9a7fd3cc5b5426dc4441ea0a96d17fb419600e54266d1ff2de07853325b3f1be827ff72a27df2ddd3db63c8f5ec0e5c98612a72919c98ae5ee05592681dbb2c43cd652baef8f899db8f4d260352529436b1bf8870fcab2fac857370c1f31c8aa7448b372ed058c87d0b79ef57a93f62fe94296861ccfa71f96d423101aafca596f7a74222d6bf1259259954f3eded938a667ab0973bf532b1239b23860fa2c307a75fa16e199a46d6b205295785528129e56028"},
		{257, "e4dceb98a05f8d6336d908c1a5d77fdf16f329c40d63e60e310efa3d6394b7efcc7c7f631ee718dea0266bba3a5c2bc2f0373654a2d4181ccf19737bd79ad04cd7d43a60524446ed944919d1220922cbd958ad151a2c06b6c72dd17b9d80c4ecddf774c906abbab8a528cad3bdecb3bffcee090c5c660e0b0b7d31329a6de8dc2866975c061684ea6fa5a39404a2ed182585a2cb2ed0937abab087875e348567465b0110b7bf9e099f60b2611e8ebb7f70ede22e0a1c668ea12e2f95ae6a742a7ce18ba366c535687ae88b9cebd854242d0f1d5cb213a11400f63ec6f3b7c2d2f5504bca05d0370b4ad4c86c553b9e9f4718e88cb3e5ee63087662ef7cda24a30a"},
	}
	for _, tt := range tests {
		h := NewHash(tt.size)
		got := h.Sum(nil)
		want, _ := hex.DecodeString(tt.want)
		if !bytes.Equal(got, want) {
			t.Errorf("Hash(\"\", len=%d) = %x, want %x", tt.size, got, want)
		}
	}
}
