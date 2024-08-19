package driver

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func NewQueue() *sqs.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err.Error())
	}

	return sqs.NewFromConfig(cfg)
}

func NewLocalQueue(url string) *sqs.Client {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						PartitionID:   "aws",
						SigningRegion: "ap-northeast-1",
						URL:           url,
					}, nil
				},
			),
		),
	)
	if err != nil {
		panic(err.Error())
	}

	return sqs.NewFromConfig(cfg)
}
