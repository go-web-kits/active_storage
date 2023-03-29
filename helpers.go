package active_storage

import (
	"github.com/awnumar/fastrand"
)

func generateBase36Key() string {
	base36Alphabet := []byte("0123456789abcdefghijklmnopqrstuvwxyz")
	result := []byte{}
	for _, b := range fastrand.Bytes(24) {
		idx := int(b % 64)
		if idx >= 36 {
			idx = fastrand.Intn(36)
		}
		result = append(result, base36Alphabet[idx])
	}
	return string(result)
}
