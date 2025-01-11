package encryption

import (
	"encoding/base64"
	"fmt"
	"github.com/devlibx/gox-base/v2"
	goxJsonUtils "github.com/devlibx/gox-base/v2/serialization/utils/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAes(t *testing.T) {
	key, err := GenerateAESKey(32)
	assert.NoError(t, err)
	fmt.Println(base64.StdEncoding.EncodeToString(key))

	encConfig := &EncryptDecryptConfigs{
		Group: map[string]*EncryptDecryptConfig{
			"test": {
				Algo: "aes_32",
				AesConfig: &AesConfig{
					Base64CodedKey: base64.StdEncoding.EncodeToString(key),
				},
			},
		},
	}
	fmt.Println(goxJsonUtils.ObjectToStringSuppressError(encConfig))

	f, err := NewServiceFactory(gox.NewNoOpCrossFunction(), encConfig)
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
