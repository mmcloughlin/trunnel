// Code generated by trunnel. DO NOT EDIT.

// +build gofuzz

package rem

func FuzzRem(data []byte) int {
	_, err := ParseRem(data)
	if err != nil {
		return 0
	}
	return 1
}