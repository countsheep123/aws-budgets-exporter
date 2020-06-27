package budgets

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type Client struct {
	cfg aws.Config
}

func New(roleArn, roleSessionName string) (*Client, error) {
	config, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	config.Credentials = credentialProvider{
		StsClient:       sts.New(config),
		RoleArn:         roleArn,
		RoleSessionName: roleSessionName,
	}

	return &Client{
		cfg: config,
	}, nil
}

type credentialProvider struct {
	StsClient       *sts.Client
	RoleArn         string
	RoleSessionName string
}

func (s credentialProvider) assumeRole(ctx context.Context) (*sts.Credentials, error) {
	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(s.RoleArn),
		RoleSessionName: aws.String(s.RoleSessionName),
	}

	out, err := s.StsClient.AssumeRoleRequest(input).Send(ctx)
	if err != nil {
		return nil, err
	}

	return out.Credentials, nil
}

func (s credentialProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	role, err := s.assumeRole(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}

	return aws.Credentials{
		AccessKeyID:     aws.StringValue(role.AccessKeyId),
		SecretAccessKey: aws.StringValue(role.SecretAccessKey),
		SessionToken:    aws.StringValue(role.SessionToken),
		Expires:         aws.TimeValue(role.Expiration),
	}, nil
}
