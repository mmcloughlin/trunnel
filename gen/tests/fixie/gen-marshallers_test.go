// Code generated by trunnel. DO NOT EDIT.

package fixie

import "testing"

func TestColorCorpus(t *testing.T) {
	cases := []struct {
		Data []byte
	}{
		{
			Data: []byte{0xfd, 0x94, 0x1},
		},
	}
	for _, c := range cases {
		_, err := ParseColor(c.Data)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFixieDemoCorpus(t *testing.T) {
	cases := []struct {
		Data []byte
	}{
		{
			Data: []byte{0x51, 0xe1, 0xec, 0x85, 0xe2, 0x2e, 0x89, 0x78, 0xa1, 0xd4, 0xf, 0x85, 0x51, 0x46, 0x39, 0xd, 0xe0, 0xb1, 0xa7, 0x9e, 0xaf, 0x48, 0x18, 0xd, 0x2d, 0xb, 0x75, 0xfb, 0x2a, 0xbd, 0xf4, 0x4a, 0x4f, 0xf9, 0x5f, 0xf6, 0x62, 0xa5, 0xee, 0xe8, 0xd3, 0xff, 0x12, 0x4, 0x5b, 0x73, 0xc8, 0x6e, 0x41, 0xc0, 0xfc, 0x2f, 0xfa, 0xc2},
		},
	}
	for _, c := range cases {
		_, err := ParseFixieDemo(c.Data)
		if err != nil {
			t.Fatal(err)
		}
	}
}
