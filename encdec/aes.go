package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/devlibx/gox-base/v2"
	"github.com/devlibx/gox-base/v2/errors"
	"io"
)

// EncryptorDecryptService is an interface to encrypt and decrypt data using AES
type aesEncryptorDecryptService struct {
	cf     gox.CrossFunction
	config *EncryptDecryptConfig
	gcm    cipher.AEAD
	nonce  []byte
}

func (n *aesEncryptorDecryptService) EncryptAndOutputBase64Ciphertext(data string) (string, error) {
	cipherText := n.gcm.Seal(nil, n.nonce, []byte(data), nil)
	result := append(n.nonce, cipherText...)
	return base64.StdEncoding.EncodeToString(result), nil
}

func (n *aesEncryptorDecryptService) DecryptFromBase64Ciphertext(data string) (string, error) {
	cipherData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode base64 ciphertext: input=%s", data)
	}

	nonce, cipherText := cipherData[:12], cipherData[12:]
	plainText, err := n.gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to decrypt ciphertext: input=%s", data)
	}
	return string(plainText), nil

}

func NewAesEncryptorDecryptService(cf gox.CrossFunction, config *EncryptDecryptConfig) (EncryptorDecryptService, error) {
	a := &aesEncryptorDecryptService{
		cf:     cf,
		config: config,
	}

	key, err := base64.StdEncoding.DecodeString(config.AesConfig.Base64CodedKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode base64 coded key to generate aes key for encryption/decryption")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cipher block")
	}

	// Generate a random nonce
	a.nonce = make([]byte, 12) // 12 bytes for GCM
	if _, err := io.ReadFull(rand.Reader, a.nonce); err != nil {
		return nil, errors.Wrap(err, "failed to generate nonce")
	}

	// Create a GCM cipher mode instance
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create GCM cipher")
	}
	a.gcm = aesGCM

	return a, nil
}

// GenerateAESKey generates a random AES key of the specified length (16, 24, or 32 bytes).
//
// We can use this as a helper to generate a random key for AES encryption/decryption.
func GenerateAESKey(length int) ([]byte, error) {
	if length != 16 && length != 24 && length != 32 {
		return nil, fmt.Errorf("invalid key length: %d (must be 16, 24, or 32 bytes)", length)
	}

	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random key: %v", err)
	}

	return key, nil
}
