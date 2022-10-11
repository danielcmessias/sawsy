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
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/danielcmessias/lfq/ui/components/table"
	"github.com/danielcmessias/lfq/utils/icons"
)


type S3Client struct {
	ctx  context.Context
    conn *s3.Client
}

func NewS3Client(cfg aws.Config, ctx context.Context) *S3Client {
	return &S3Client{
		ctx:  ctx,
		conn: s3.NewFromConfig(cfg),
	}
}

func (c *S3Client) GetBucketsRows(nextToken *string) ([]table.Row, *string, error) {
	input := s3.ListBucketsInput{}
	output, err := c.conn.ListBuckets(c.ctx, &input)
    if err != nil {
        log.Fatalf("unable to get databases: %v", err)
    }

	var wg sync.WaitGroup
	regionsCh := make(chan struct {int; string}, len(output.Buckets))

	var rows []table.Row
	for i, b := range output.Buckets {
		wg.Add(1)
		go func(rowIndex int, bucket string, ) {
			defer wg.Done()

			region, err := manager.GetBucketRegion(c.ctx, c.conn, bucket)
			regionsCh <- struct {int; string}{rowIndex, region}

			if err != nil {
				fmt.Println("Error getting regio for", bucket)
			}
		}(i, aws.ToString(b.Name))

		rows = append(rows, table.Row{
			aws.ToString(b.Name),
			"...",
			b.CreationDate.String(),
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


func (c *S3Client) GetObjectsRows(bucket string, region string, prefix string, nextToken *string) ([]table.Row, *string, error) {
	input := s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Delimiter: aws.String("/"),
		EncodingType: types.EncodingTypeUrl,
		Prefix: aws.String(prefix),
	}
	// if nextToken != nil {
    //     input.ContinuationToken = nextToken
    // }

	output, err := c.conn.ListObjectsV2(c.ctx, &input, func(options *s3.Options) { options.Region = region })
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
			o.LastModified.String(),
			strconv.FormatInt(o.Size, 10),
		})
	}

	// Don't return a next token, for large buckets it'll just kill the app with too much data
	// Will need to send the search query to the server instead
	return rows, nil, nil
}

