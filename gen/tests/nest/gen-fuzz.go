// Code generated by trunnel. DO NOT EDIT.

// +build gofuzz

package nest

func FuzzPoint(data []byte) int {
	_, err := ParsePoint(data)
	if err != nil {
		return 0
	}
	return 1
}

func FuzzRect(data []byte) int {
	_, err := ParseRect(data)
	if err != nil {
		return 0
	}
	return 1
}
