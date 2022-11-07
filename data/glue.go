package data

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/url"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/danielcmessias/sawsy/ui/components/table"
)

type GlueClient struct {
	ctx  context.Context
	glue *glue.Client
	s3   *s3.Client
}

func NewGlueClient(ctx context.Context, glue *glue.Client, s3 *s3.Client) *GlueClient {
	return &GlueClient{
		ctx:  ctx,
		glue: glue,
		s3:   s3,
	}
}

func (c *GlueClient) GetJobsRows(nextToken *string) ([]table.Row, *string, error) {
	listInput := glue.ListJobsInput{
		NextToken: nextToken,
	}
	listOutput, err := c.glue.ListJobs(c.ctx, &listInput)
	if err != nil {
		log.Fatalf("unable to get Glue jobs: %v", err)
	}

	if len(listOutput.JobNames) == 0 {
		return nil, nil, nil
	}

	getInput := glue.BatchGetJobsInput{
		JobNames: listOutput.JobNames,
	}
	getOutput, err := c.glue.BatchGetJobs(c.ctx, &getInput)
	if err != nil {
		log.Fatalf("unable to get Glue jobs: %v", err)
	}

	var rows []table.Row
	for _, j := range getOutput.Jobs {
		rows = append(rows, table.Row{
			aws.ToString(j.Name),
			aws.ToString(j.Command.Name),
			formatTime(j.LastModifiedOn),
			aws.ToString(j.GlueVersion),
			string(j.WorkerType),
			fmt.Sprintf("%d", aws.ToInt32(j.NumberOfWorkers)),
		})
	}

	return rows, listOutput.NextToken, nil
}

func (c *GlueClient) GetCrawlersRows(nextToken *string) ([]table.Row, *string, error) {
	listInput := glue.ListCrawlersInput{
		NextToken: nextToken,
	}
	listOutput, err := c.glue.ListCrawlers(c.ctx, &listInput)
	if err != nil {
		log.Fatalf("unable to get Glue jobs: %v", err)
	}

	if len(listOutput.CrawlerNames) == 0 {
		return nil, nil, nil
	}

	getInput := glue.BatchGetCrawlersInput{
		CrawlerNames: listOutput.CrawlerNames,
	}
	getOutput, err := c.glue.BatchGetCrawlers(c.ctx, &getInput)
	if err != nil {
		log.Fatalf("unable to get Glue crawlers: %v", err)
	}

	metricsInput := glue.GetCrawlerMetricsInput{
		CrawlerNameList: listOutput.CrawlerNames,
	}
	metricsOutput, err := c.glue.GetCrawlerMetrics(c.ctx, &metricsInput)
	if err != nil {
		log.Fatalf("unable to get Glue crawlers: %v", err)
	}

	if len(getOutput.Crawlers) != len(metricsOutput.CrawlerMetricsList) {
		log.Fatalf("GetCrawlers and GetCrawlerMetrics have different lengths.")
	}

	var rows []table.Row
	for i, c := range getOutput.Crawlers {
		metrics := metricsOutput.CrawlerMetricsList[i]

		var schedule string
		if c.Schedule != nil {
			schedule = aws.ToString(c.Schedule.ScheduleExpression)
		}

		rows = append(rows, table.Row{
			aws.ToString(c.Name),
			schedule,
			string(c.State),
			formatSeconds(int(metrics.LastRuntimeSeconds)),
			formatSeconds(int(metrics.MedianRuntimeSeconds)),
		})
	}

	return rows, listOutput.NextToken, nil
}

func (c *GlueClient) GetJobDetails(jobName string) ([]table.Row, error) {
	input := glue.GetJobInput{
		JobName: aws.String(jobName),
	}
	output, err := c.glue.GetJob(c.ctx, &input)
	if err != nil {
		log.Fatalf("unable to get Glue job details: %v", err)
	}

	job := output.Job
	rows := []table.Row{
		{"Job name", aws.ToString(job.Name)},
		{"Description", aws.ToString(job.Description)},
		{"Job Type", aws.ToString(job.Command.Name)},
		{"Role", aws.ToString(job.Role)},
		{"Glue Version", aws.ToString(job.GlueVersion)},
		{"Worker Type", string(job.WorkerType)},
		{"# Workers", fmt.Sprintf("%d", aws.ToInt32(job.NumberOfWorkers))},
		{"Max Retries", fmt.Sprintf("%d", job.MaxRetries)},
		{"Timeout", fmt.Sprintf("%dm", aws.ToInt32(job.Timeout))},
		{"Script", aws.ToString(job.Command.ScriptLocation)},
		{"Created On", formatTime(job.CreatedOn)},
		{"Last Modified", formatTime(job.LastModifiedOn)},
	}
	return rows, nil
}

func (c *GlueClient) GetJobRuns(jobName string) ([]table.Row, error) {
	input := glue.GetJobRunsInput{
		JobName: aws.String(jobName),
	}
	output, err := c.glue.GetJobRuns(c.ctx, &input)
	if err != nil {
		log.Fatalf("unable to get Glue job runs: %v", err)
	}

	var rows []table.Row
	for _, r := range output.JobRuns {
		rows = append(rows, table.Row{
			formatTime(r.StartedOn),
			string(r.JobRunState),
			fmt.Sprintf("%d", aws.ToInt32(&r.Attempt)),
			formatSeconds(int(r.ExecutionTime)),
		})
	}

	return rows, nil
}

func (c *GlueClient) GetJobScript(jobName string) (string, string, error) {
	jobInput := glue.GetJobInput{
		JobName: aws.String(jobName),
	}
	jobOutput, err := c.glue.GetJob(c.ctx, &jobInput)
	if err != nil {
		log.Fatalf("unable to get Glue job details: %v", err)
	}

	scriptLocation := aws.ToString(jobOutput.Job.Command.ScriptLocation)
	u, _ := url.Parse(scriptLocation)

	// log.Fatalf("scriptLocation: %s, bucket: %s, key: %s", scriptLocation, u.Host, u.Path)
	objInput := s3.GetObjectInput{
		Bucket: aws.String(u.Host),
		Key:    aws.String(u.Path[1:]),
	}
	objOutput, err := c.s3.GetObject(c.ctx, &objInput)
	if err != nil {
		return err.Error(), "", nil
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(objOutput.Body)
	if err != nil {
		log.Fatalf("unable to get Glue job details: %v", err)
	}

	return buf.String(), scriptLocation, nil
}
