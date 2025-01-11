package encryption

import (
	"encoding/base64"
	"fmt"
	"github.com/devlibx/gox-base/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAes(t *testing.T) {
	key, err := GenerateAESKey(32)
	assert.NoError(t, err)

	f, err := NewServiceFactory(gox.NewNoOpCrossFunction(), &EncryptDecryptConfigs{
		Group: map[string]*EncryptDecryptConfig{
			"test": {
				Algo: "aes_32",
				AesConfig: &AesConfig{
					Base64CodedKey: base64.StdEncoding.EncodeToString(key),
				},
			},
		},
	})
	assert.NoError(t, err)

	s, err := f.GetEncryptorDecryptService("test")
	assert.NoError(t, err)

	base64EncryptedData, err := s.EncryptAndOutputBase64Ciphertext("My Name Is Abcd")
	assert.NoError(t, err)

	out, err := s.DecryptFromBase64Ciphertext(base64EncryptedData)
	assert.NoError(t, err)

	fmt.Println(out)
	assert.Equal(t, "My Name Is Abcd", out)
}
