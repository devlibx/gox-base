package encryption

type noOpServiceFactoryImpl struct {
}

func (n noOpServiceFactoryImpl) GetEncryptorDecryptService(groupName string) (EncryptorDecryptService, error) {
	return &noOpEncryptorDecryptService{}, nil
}

type noOpEncryptorDecryptService struct {
}

func (n noOpEncryptorDecryptService) EncryptAndOutputBase64Ciphertext(data string) (string, error) {
	return data, nil
}

func (n noOpEncryptorDecryptService) DecryptFromBase64Ciphertext(data string) (string, error) {
	return data, nil
}
