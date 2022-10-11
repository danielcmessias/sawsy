package data

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

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

func (c *S3Client) GetBuckets(nextToken *string) ([]table.Row, *string, error) {
	input := s3.ListBucketsInput{}
	output, err := c.s3.ListBuckets(c.ctx, &input)
	if err != nil {
		log.Fatalf("unable to get databases: %v", err)
	}

	var wg sync.WaitGroup
	regionsCh := make(chan struct {
		int
		string
	}, len(output.Buckets))

	var rows []table.Row
	for i, b := range output.Buckets {
		wg.Add(1)
		go func(rowIndex int, bucket string) {
			defer wg.Done()

			region, err := manager.GetBucketRegion(c.ctx, c.s3, bucket)
			regionsCh <- struct {
				int
				string
			}{rowIndex, region}

			if err != nil {
				log.Fatalf("Error getting region for bucket %s, %v", bucket, err)
			}
		}(i, aws.ToString(b.Name))

		rows = append(rows, table.Row{
			aws.ToString(b.Name),
			"...",
			formatTime(b.CreationDate),
		})
	}

	go func() {
		wg.Wait()
		close(regionsCh)
	}()
	for item := range regionsCh {
		rows[item.int][1] = item.string
	}

	return rows, nil, nil
}

func (c *S3Client) GetBucketPolicy(bucket string, region string) (string, error) {
	input := s3.GetBucketPolicyInput{
		Bucket: aws.String(bucket),
	}
	output, err := c.s3.GetBucketPolicy(c.ctx, &input, func(options *s3.Options) { options.Region = region })
	if err != nil {
		return "", err
	}

	return formatJson(*output.Policy), nil
}

func (c *S3Client) GetBucketTags(bucket string, region string) ([]table.Row, error) {
	input := s3.GetBucketTaggingInput{
		Bucket: aws.String(bucket),
	}
	output, err := c.s3.GetBucketTagging(c.ctx, &input, func(options *s3.Options) { options.Region = region })
	if err != nil {
		log.Fatalf("unable to get bucket tags: %v", err)
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
		log.Fatalf("unable to get objects: %v", err)
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
		log.Fatalf("unable to get object properties: %v", err)
	}

	rows = append(rows, table.Row{"Last Modified", formatTime(headOutput.LastModified)})
	rows = append(rows, table.Row{"ETag", aws.ToString(headOutput.ETag)})

	aclInput := s3.GetObjectAclInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	aclOutput, err := c.s3.GetObjectAcl(c.ctx, &aclInput, func(options *s3.Options) { options.Region = region })
	if err != nil {
		log.Fatalf("unable to get object acl: %v", err)
	}

	rows = append(rows, table.Row{"Owner", aws.ToString(aclOutput.Owner.DisplayName)})

	return rows, nil
}
