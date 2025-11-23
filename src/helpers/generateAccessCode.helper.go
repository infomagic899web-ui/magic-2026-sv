package helpers

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// GenerateAccessCode generates a random alphanumeric access code
func GenerateAccessCode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var code strings.Builder
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic("Failed to generate random access code: " + err.Error())
		}
		code.WriteByte(charset[index.Int64()])
	}
	return code.String()
}
