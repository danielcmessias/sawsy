package data

import (
	"context"
	"fmt"
	"log"
	"time"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/danielcmessias/sawsy/ui/components/table"
)

type RDSClient struct {
	ctx        context.Context
	rds        *rds.Client
	cloudwatch *cloudwatch.Client
}

func NewRDSClient(ctx context.Context, rds *rds.Client, cloudwatch *cloudwatch.Client) *RDSClient {
	return &RDSClient{
		ctx:        ctx,
		rds:        rds,
		cloudwatch: cloudwatch,
	}
}

func (c *RDSClient) GetDBInstances(nextToken *string) ([]table.Row, *string, error) {
	input := rds.DescribeDBInstancesInput{
		Marker: nextToken,
	}
	output, err := c.rds.DescribeDBInstances(c.ctx, &input)
	if err != nil {
		log.Fatalf("unable to list db instances: %v", err)
	}

	var rows []table.Row
	for _, i := range output.DBInstances {
		rows = append(rows, table.Row{
			aws.ToString(i.DBInstanceIdentifier),
			aws.ToString(i.Engine),
			aws.ToString(i.AvailabilityZone),
			aws.ToString(i.DBInstanceClass),
			aws.ToString(i.DBInstanceStatus),
			aws.ToString(i.DBSubnetGroup.VpcId),
			formatBool(i.MultiAZ),
		})
	}
	return rows, output.Marker, nil
}

func (c *RDSClient) GetInstanceDetails(instance string) ([]table.Row, error) {
	input := rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(instance),
	}
	output, err := c.rds.DescribeDBInstances(c.ctx, &input)
	if err != nil {
		log.Fatalf("unable to get db instance: %v", err)
	}

	i := output.DBInstances[0]
	rows := []table.Row{
		{"Identifier", aws.ToString(i.DBInstanceIdentifier)},
		{"Engine", aws.ToString(i.Engine)},
		{"Region & AZ", aws.ToString(i.AvailabilityZone)},
		{"Size", aws.ToString(i.DBInstanceClass)},
		{"Status", aws.ToString(i.DBInstanceStatus)},
		{"Endpoint", aws.ToString(i.Endpoint.Address)},
		{"Port", fmt.Sprintf("%d", i.Endpoint.Port)},
		{"VPC", aws.ToString(i.DBSubnetGroup.VpcId)},
		{"Multi-AZ", formatBool(i.MultiAZ)},
		{"ARN", aws.ToString(i.DBInstanceArn)},
		{"Created on", formatTime(i.InstanceCreateTime)},
		{"Storage", fmt.Sprintf("%d", i.AllocatedStorage)},
		{"Storage type", aws.ToString(i.StorageType)},
	}
	return rows, nil
}

func (c *RDSClient) GetInstanceTags(instance string) ([]table.Row, error) {
	input := rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(instance),
	}
	output, err := c.rds.DescribeDBInstances(c.ctx, &input)
	if err != nil {
		log.Fatalf("unable to get db instance tags: %v", err)
	}

	var rows []table.Row
	for _, tag := range output.DBInstances[0].TagList {
		rows = append(rows, table.Row{
			aws.ToString(tag.Key),
			aws.ToString(tag.Value),
		})
	}
	return rows, nil
}

func (c *RDSClient) GetMetric(instance string, metricName string) ([]float64, error) {
	endTime := time.Now()
	startTime := endTime.Add(time.Hour * -3)

	input := cloudwatch.GetMetricStatisticsInput{
		MetricName: aws.String(metricName),
		Namespace:  aws.String("AWS/RDS"),
		Period:     aws.Int32(60),
		StartTime:  aws.Time(startTime),
		EndTime:    aws.Time(endTime),
		Statistics: []types.Statistic{
			types.StatisticAverage,
		},
		Dimensions: []types.Dimension{
			{
				Name:  aws.String("DBInstanceIdentifier"),
				Value: aws.String(instance),
			},
		},
	}
	output, err := c.cloudwatch.GetMetricStatistics(c.ctx, &input)
	if err != nil {
		log.Fatalf("unable to get metric: %v", err)
	}

	datapoints := make([]float64, len(output.Datapoints))
	for i, d := range output.Datapoints {
		datapoints[i] = aws.ToFloat64(d.Average)
	}

	return datapoints, nil
}
