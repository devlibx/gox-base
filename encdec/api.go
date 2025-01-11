package encryption

// EncryptDecryptConfigs is the top-level configuration for encryption/decryption.
//
// We can provide one ore more group of encryption/decryption configurations.
type EncryptDecryptConfigs struct {
	Disabled bool                             `json:"disabled"`
	Group    map[string]*EncryptDecryptConfig `json:"group"`
}

// EncryptDecryptConfig is the configuration for a single group of encryption/decryption.
type EncryptDecryptConfig struct {
	Disabled  bool       `json:"disabled"`
	Algo      string     `json:"algorithm"`
	AesConfig *AesConfig `json:"aes_config"`
}

// AesConfig is the configuration for AES encryption/decryption.
type AesConfig struct {
	Base64CodedKey string `json:"base_64_coded_key"`
}

// ServiceFactory is the interface for creating a new encryption/decryption service.
type ServiceFactory interface {
	GetEncryptorDecryptService(groupName string) (EncryptorDecryptService, error)
}

// EncryptorDecryptService is the interface for encrypting and decrypting data.
type EncryptorDecryptService interface {

	// EncryptAndOutputBase64Ciphertext encrypts the input data and returns the base64-encoded ciphertext.
	// The output is a base64-encoded string that contains the nonce and the ciphertext.
	EncryptAndOutputBase64Ciphertext(data string) (string, error)

	// DecryptFromBase64Ciphertext decrypts the base64-encoded ciphertext and returns the original data.
	// The input is a base64-encoded string that contains the nonce and the ciphertext.
	DecryptFromBase64Ciphertext(data string) (string, error)
}
