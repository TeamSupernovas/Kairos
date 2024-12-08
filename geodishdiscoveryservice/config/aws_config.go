package config

import "os"

type AWSConfig struct {
	s3Bucket           string
	region             string
	locationPlaceIndex string
	credentials        AWSCredentials
}

type AWSCredentials struct {
	accessKey string
	secretKey string
}

func loadAWSConfig() AWSConfig {
	return AWSConfig{
		s3Bucket:           os.Getenv("AWS_S3_BUCKET"),
		region:             os.Getenv("AWS_REGION"),
		locationPlaceIndex: os.Getenv("AWS_LOCATION_PLACE_INDEX_NAME"),
		credentials: AWSCredentials{
			accessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
			secretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
	}
}

func (ac AWSConfig) S3Bucket() string {
	return ac.s3Bucket
}

func (ac AWSConfig) Region() string {
	return ac.region
}

func (ac AWSConfig) LocationPlaceIndex() string {
	return ac.locationPlaceIndex
}

func (ac AWSConfig) Credentials() AWSCredentials {
	return ac.credentials
}

// AWSCredentials getters
func (cred AWSCredentials) AccessKey() string {
	return cred.accessKey
}

func (cred AWSCredentials) SecretKey() string {
	return cred.secretKey
}
