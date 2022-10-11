package data

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lakeformation"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

const MAX_RESULTS = 32

type Client struct {
	ctx context.Context
	sts *sts.Client

	Glue          *GlueClient
	IAM           *IAMClient
	LakeFormation *LakeFormationClient
	Lambda        *LambdaClient
	RDS           *RDSClient
	S3            *S3Client
}

func NewClient() Client {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	cloudwatch := cloudwatch.NewFromConfig(cfg)
	glue := glue.NewFromConfig(cfg)
	iam := iam.NewFromConfig(cfg)
	lakeformation := lakeformation.NewFromConfig(cfg)
	lambda := lambda.NewFromConfig(cfg)
	rds := rds.NewFromConfig(cfg)
	s3 := s3.NewFromConfig(cfg)

	return Client{
		ctx: ctx,
		sts: sts.NewFromConfig(cfg),

		Glue:          NewGlueClient(ctx, glue, s3),
		IAM:           NewIAMClient(ctx, iam),
		LakeFormation: NewLakeFormationClient(ctx, lakeformation, glue),
		Lambda:        NewLambdaClient(ctx, lambda, cloudwatch),
		RDS:           NewRDSClient(ctx, rds, cloudwatch),
		S3:            NewS3Client(ctx, s3),
	}
}

func (c *Client) GetCurrentAWSAccountId() string {
	input := sts.GetCallerIdentityInput{}
	output, err := c.sts.GetCallerIdentity(c.ctx, &input)
	if err != nil {
		log.Fatalf("unable to get caller identity: %v", err)
	}

	return *output.Account
}
