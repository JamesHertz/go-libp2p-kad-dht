package encrypt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testKey = []byte("AES256Key-32Characters1234567890")

func TestAES(t *testing.T) {
	plaintext := []byte("nootwashere")
	ciphertext, err := EncryptAES(plaintext, testKey)
	require.NoError(t, err)
	plaintextRes, err := DecryptAES(ciphertext, testKey)
	require.NoError(t, err)
	require.Equal(t, plaintext, plaintextRes)
}
