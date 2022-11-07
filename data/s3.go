package data

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/utils/icons"
)

type S3Client struct {
	ctx context.Context
	s3  *s3.Client
}

func NewS3Client(ctx context.Context, s3 *s3.Client) *S3Client {
	return &S3Client{
		ctx: ctx,
		s3:  s3,
	}
}

func (c *S3Client) GetBuckets() ([]table.Row, error) {
	input := s3.ListBucketsInput{}
	output, err := c.s3.ListBuckets(c.ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("error listing buckets: %w", err)
	}

	var rows []table.Row
	for _, b := range output.Buckets {
		rows = append(rows, table.Row{
			aws.ToString(b.Name),
			LOADING_ALIAS,
			formatTime(b.CreationDate),
		})
	}

	return rows, nil
}

func (c *S3Client) GetBucketRegion(bucket string) (string, error) {
	region, err := manager.GetBucketRegion(c.ctx, c.s3, bucket)
	if err != nil {
		return "", fmt.Errorf("error getting region for bucket %s: %w", bucket, err)
	}
	return region, nil
}

func (c *S3Client) GetBucketPolicy(bucket string, region string) (string, error) {
	input := s3.GetBucketPolicyInput{
		Bucket: aws.String(bucket),
	}
	output, err := c.s3.GetBucketPolicy(c.ctx, &input, func(options *s3.Options) { options.Region = region })
	if err != nil {
		return "", fmt.Errorf("error getting policy for bucket %s: %w", bucket, err)
	}

	return formatJson(*output.Policy), nil
}

func (c *S3Client) GetBucketTags(bucket string, region string) ([]table.Row, error) {
	input := s3.GetBucketTaggingInput{
		Bucket: aws.String(bucket),
	}
	output, err := c.s3.GetBucketTagging(c.ctx, &input, func(options *s3.Options) { options.Region = region })
	if err != nil {
		return nil, fmt.Errorf("error getting tags for bucket %s: %w", bucket, err)
	}

	var rows []table.Row
	for _, t := range output.TagSet {
		rows = append(rows, table.Row{aws.ToString(t.Key), aws.ToString(t.Value)})
	}

	return rows, nil
}

func (c *S3Client) GetObjects(bucket string, region string, prefix string, nextToken *string) ([]table.Row, *string, error) {
	input := s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Delimiter: aws.String("/"),
		Prefix:    aws.String(prefix),
	}
	// if nextToken != nil {
	//     input.ContinuationToken = nextToken
	// }

	output, err := c.s3.ListObjectsV2(c.ctx, &input, func(options *s3.Options) { options.Region = region })
	if err != nil {
		return nil, nil, fmt.Errorf("error listing objects for bucket %s with prefix %s: %w", bucket, prefix, err)
	}

	var rows []table.Row

	for _, o := range output.CommonPrefixes {
		rows = append(rows, table.Row{
			fmt.Sprintf("%s %s", icons.FOLDER, strings.Replace(aws.ToString(o.Prefix), prefix, "", 1)),
			"-",
			"-",
		})
	}

	for _, o := range output.Contents {
		rows = append(rows, table.Row{
			fmt.Sprintf("%s %s", icons.FILE, strings.Replace(aws.ToString(o.Key), prefix, "", 1)),
			formatTime(o.LastModified),
			strconv.FormatInt(o.Size, 10),
		})
	}

	// Don't return a next token, for large buckets it'll just kill the app with too much data
	// Will need to send the search query to the server instead
	return rows, nil, nil
}

func (c *S3Client) GetObjectProperties(bucket string, key string, region string) ([]table.Row, error) {
	rows := []table.Row{
		{"Bucket", bucket},
		{"Key", key},
		{"Region", region},
		{"ARN", fmt.Sprintf("arn:aws:s3:::%s/%s", bucket, key)},
	}

	headInput := s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	headOutput, err := c.s3.HeadObject(c.ctx, &headInput, func(options *s3.Options) { options.Region = region })
	if err != nil {
		return nil, fmt.Errorf("error getting properties for object s3://%s/%s: %w", bucket, key, err)
	}

	rows = append(rows, table.Row{"Last Modified", formatTime(headOutput.LastModified)})
	rows = append(rows, table.Row{"ETag", aws.ToString(headOutput.ETag)})

	aclInput := s3.GetObjectAclInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	aclOutput, err := c.s3.GetObjectAcl(c.ctx, &aclInput, func(options *s3.Options) { options.Region = region })
	if err != nil {
		return nil, fmt.Errorf("error getting ACL for object s3://%s/%s: %w", bucket, key, err)
	}

	rows = append(rows, table.Row{"Owner", aws.ToString(aclOutput.Owner.DisplayName)})

	return rows, nil
}
