package encryption

// NoOpServiceFactory is a factory to create encryptor decryptor service
// This is a no-op implementation of the service factory
type noOpServiceFactoryImpl struct {
}

// GetEncryptorDecryptService returns a new instance of No-OP EncryptorDecryptService
func (n noOpServiceFactoryImpl) GetEncryptorDecryptService(groupName string) (EncryptorDecryptService, error) {
	return &noOpEncryptorDecryptService{}, nil
}

// noOpEncryptorDecryptService a no-op implementation of EncryptorDecryptService
type noOpEncryptorDecryptService struct {
}

// EncryptAndOutputBase64Ciphertext no-op, just return the input as output
func (n noOpEncryptorDecryptService) EncryptAndOutputBase64Ciphertext(data string) (string, error) {
	return data, nil
}

// DecryptFromBase64Ciphertext no-op, just return the input as output
func (n noOpEncryptorDecryptService) DecryptFromBase64Ciphertext(data string) (string, error) {
	return data, nil
}
