// Code generated by trunnel. DO NOT EDIT.

package rem

import "testing"

func TestRemCorpus(t *testing.T) {
	cases := []struct {
		Data []byte
	}{
		{
			Data: []byte{0xa5, 0xee, 0xe8, 0x2a, 0x62, 0xf6, 0x5f, 0xf9, 0x4f, 0x6e, 0xc8, 0x73, 0x5b, 0x4, 0x12, 0xff, 0xd3, 0x41, 0xc0},
		},
	}
	for _, c := range cases {
		_, err := ParseRem(c.Data)
		if err != nil {
			t.Fatal(err)
		}
	}
}