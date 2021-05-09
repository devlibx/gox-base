package goxaws

import "github.com/aws/aws-sdk-go/aws/session"

// Dynamo DB setup configuration
type Config struct {
	Region   string `json:"region,string" yaml:"region,string"`
	Endpoint string `json:"endpoint,string" yaml:"endpoint,string"`
}

type AwsContext interface {
	GetSession() *session.Session
}
