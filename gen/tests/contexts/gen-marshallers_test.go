// Code generated by trunnel. DO NOT EDIT.

package contexts

import "testing"

func TestPointCorpus(t *testing.T) {
	cases := []struct {
		Data []byte
	}{
		{
			Data: []byte{0x25, 0x1},
		},
	}
	for _, c := range cases {
		_, err := ParsePoint(c.Data)
		if err != nil {
			t.Fatal(err)
		}
	}
}
