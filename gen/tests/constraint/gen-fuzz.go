// Code generated by trunnel. DO NOT EDIT.

// +build gofuzz

package constraint

func FuzzDate(data []byte) int {
	_, err := ParseDate(data)
	if err != nil {
		return 0
	}
	return 1
}
