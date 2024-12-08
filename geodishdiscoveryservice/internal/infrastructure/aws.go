package infrastructure

import (
	"context"
	"fmt"
	"geodishdiscoveryservice/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// InitAWS initializes the AWS configuration
func initAWS(cfg config.AWSConfig) (aws.Config, error) {
	awsConfig, err := awsConfig.LoadDefaultConfig(context.Background(),
	awsConfig.WithRegion(cfg.Region()),
	awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.Credentials().AccessKey(), cfg.Credentials().SecretKey(), ""),
		),
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS config: %w", err)
	}
	return awsConfig, nil
}
