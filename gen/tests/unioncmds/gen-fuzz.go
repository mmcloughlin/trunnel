// Code generated by trunnel. DO NOT EDIT.

// +build gofuzz

package unioncmds

func FuzzUnionCmds(data []byte) int {
	_, err := ParseUnionCmds(data)
	if err != nil {
		return 0
	}
	return 1
}