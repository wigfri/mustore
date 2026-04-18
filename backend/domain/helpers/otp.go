package helpers

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTPCode(digits int) (string, error) {
	if digits <= 0 || digits > 12 {
		return "", fmt.Errorf("invalid digits length")
	}
	buf := make([]byte, digits)
	for i := range buf {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		buf[i] = '0' + byte(n.Int64())
	}
	return string(buf), nil
}
