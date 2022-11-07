package data

import (
	"context"
	"fmt"
	"time"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/danielcmessias/sawsy/ui/components/table"
)

type LambdaClient struct {
	ctx        context.Context
	lambda     *lambda.Client
	cloudwatch *cloudwatch.Client
}

func NewLambdaClient(ctx context.Context, lambda *lambda.Client, cloudwatch *cloudwatch.Client) *LambdaClient {
	return &LambdaClient{
		ctx:        ctx,
		lambda:     lambda,
		cloudwatch: cloudwatch,
	}
}

func (c *LambdaClient) GetFunctions(nextToken *string) ([]table.Row, *string, error) {
	input := lambda.ListFunctionsInput{
		Marker: nextToken,
	}
	output, err := c.lambda.ListFunctions(c.ctx, &input)
	if err != nil {
		return nil, nil, fmt.Errorf("error listing lambda functions: %w", err)
	}

	var rows []table.Row
	for _, f := range output.Functions {
		lastModifiedTime, _ := time.Parse(ISO_8601, aws.ToString(f.LastModified))

		rows = append(rows, table.Row{
			aws.ToString(f.FunctionName),
			string(f.Runtime),
			formatTime(&lastModifiedTime),
			aws.ToString(f.Description),
			aws.ToString(f.FunctionArn),
		})
	}

	return rows, output.NextMarker, nil
}

func (c *LambdaClient) GetFunctionDetails(functionName string) ([]table.Row, error) {
	input := lambda.GetFunctionInput{
		FunctionName: aws.String(functionName),
	}
	output, err := c.lambda.GetFunction(c.ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("error getting lambda function %s details: %w", functionName, err)
	}

	cfg := output.Configuration
	lastModifiedTime, _ := time.Parse(ISO_8601, aws.ToString(cfg.LastModified))
	var vpcId string
	if cfg.VpcConfig != nil {
		vpcId = aws.ToString(cfg.VpcConfig.VpcId)
	}

	rows := []table.Row{
		{"Function name", aws.ToString(cfg.FunctionName)},
		{"Runtime", aws.ToString(cfg.FunctionName)},
		{"Memory", fmt.Sprint(aws.ToInt32(cfg.MemorySize))},
		{"Last modified", formatTime(&lastModifiedTime)},
		{"Description", aws.ToString(cfg.Description)},
		{"ARN", aws.ToString(cfg.FunctionArn)},
		{"Handler", aws.ToString(cfg.Handler)},
		{"Role", aws.ToString(cfg.Role)},
		{"VPC Id", vpcId},
	}

	return rows, nil
}

func (c *LambdaClient) GetMetric(functionName string, metricName string, statistic types.Statistic) ([]float64, error) {
	endTime := time.Now()
	startTime := endTime.Add(time.Hour * -24)

	input := cloudwatch.GetMetricStatisticsInput{
		MetricName: aws.String(metricName),
		Namespace:  aws.String("AWS/Lambda"),
		Period:     aws.Int32(300),
		StartTime:  aws.Time(startTime),
		EndTime:    aws.Time(endTime),
		Statistics: []types.Statistic{
			statistic,
		},
		Dimensions: []types.Dimension{
			{
				Name:  aws.String("FunctionName"),
				Value: aws.String(functionName),
			},
		},
	}
	output, err := c.cloudwatch.GetMetricStatistics(c.ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("error getting metric %s for lambda function %s: %w", metricName, functionName, err)
	}

	datapoints := make([]float64, len(output.Datapoints))
	for i, d := range output.Datapoints {
		datapoints[i] = aws.ToFloat64(statisticOfDatapoint(d, statistic))
	}

	return datapoints, nil
}
