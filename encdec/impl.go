package encryption

import (
	"github.com/devlibx/gox-base/v2"
	"github.com/devlibx/gox-base/v2/errors"
)

type serviceFactoryImpl struct {
	cf                         gox.CrossFunction
	encryptDecryptConfigs      *EncryptDecryptConfigs
	encryptorDecryptServiceMap map[string]EncryptorDecryptService
}

func (s serviceFactoryImpl) GetEncryptorDecryptService(groupName string) (EncryptorDecryptService, error) {
	if out, ok := s.encryptorDecryptServiceMap[groupName]; ok {
		return out, nil
	}
	return &noOpEncryptorDecryptService{}, nil
}

func NewServiceFactory(cf gox.CrossFunction, encryptDecryptConfigs *EncryptDecryptConfigs) (ServiceFactory, error) {
	sf := &serviceFactoryImpl{
		cf:                         cf,
		encryptDecryptConfigs:      encryptDecryptConfigs,
		encryptorDecryptServiceMap: map[string]EncryptorDecryptService{},
	}
	if encryptDecryptConfigs.Disabled {
		return &noOpServiceFactoryImpl{}, nil
	}

	for group, config := range encryptDecryptConfigs.Group {

		// If it is disabled then set a no op service
		if config.Disabled {
			sf.encryptorDecryptServiceMap[group] = &noOpEncryptorDecryptService{}
			continue
		}

		switch config.Algo {
		case "aes_32":
			if t, err := NewAesEncryptorDecryptService(cf, config); err != nil {
				return nil, errors.Wrap(err, "failed to create aes encryptor/decryptor service")
			} else {
				sf.encryptorDecryptServiceMap[group] = t
			}
		default:
			return nil, errors.New("invalid algo - only aes_32 is supported: algorithm=%s", config.Algo)
		}
	}
	return sf, nil
}
