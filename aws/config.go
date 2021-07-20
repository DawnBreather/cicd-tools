package aws

import (
	"context"
	"github.com/DawnBreather/go-commons/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var _logger = logger.New()

func newConfig(region string) aws.Config {
	var cfg aws.Config
	var err error
	if region == "" {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	}
	if err != nil {
		_logger.Fatalf("unable to load SDK config: %v", err)
	}

	return cfg
}
