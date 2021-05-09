package goxaws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/devlibx/gox-base"
	errors2 "github.com/devlibx/gox-base/errors"
)

type awsSessionContext struct {
	session *session.Session
	gox.CrossFunction
}

func (a *awsSessionContext) GetSession() *session.Session {
	return a.session
}

func NewAwsContext(config Config) (ctx AwsContext, err error) {
	_ctx := &awsSessionContext{CrossFunction: gox.NewNoOpCrossFunction()}

	// Setup AWS session
	_ctx.session, err = session.NewSession(
		&aws.Config{
			Region:   aws.String(config.Region),
			Endpoint: aws.String(config.Endpoint),
		})
	if err != nil {
		return nil, errors2.Wrap(err, "failed to create aws session: endpoint=%s, region=%s", config.Endpoint, config.Region)
	}

	return _ctx, err
}
