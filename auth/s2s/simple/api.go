package goxAuthSimpleS2S

import "fmt"

type AuthConfig struct {
	Disabled   bool                     `yaml:"disabled" json:"disabled"`
	SecretSalt string                   `yaml:"secret_salt" json:"secret_salt"`
	Clients    map[string]*ClientConfig `yaml:"clients" json:"clients"`
}

type ClientConfig struct {
	ClientName     string   `yaml:"client_name" json:"client_name"`
	ClientSecret   string   `yaml:"client_secret" json:"client_secret"`
	AllowedActions []string `yaml:"allowed_actions" json:"allowed_actions"`
}

type AuthError struct {
	Err      error
	Reason   string
	ClientId string
}

func (e *AuthError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("auth failed for client_id=%s with reason=%s and error=%s", e.ClientId, e.Reason, e.Err.Error())
	}
	return fmt.Sprintf("auth failed for client_id=%s with reason=%s", e.ClientId, e.Reason)
}
