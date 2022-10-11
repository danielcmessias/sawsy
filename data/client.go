package data

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	gluetypes "github.com/aws/aws-sdk-go-v2/service/glue/types"
	"github.com/aws/aws-sdk-go-v2/service/lakeformation"
	lftypes "github.com/aws/aws-sdk-go-v2/service/lakeformation/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/danielcmessias/lfq/ui/components/table"
	"github.com/danielcmessias/lfq/utils"
)

type Client struct {
    ctx  context.Context
    glue *glue.Client
    lf   *lakeformation.Client
    sts  *sts.Client

    S3   *S3Client
}

func NewClient() Client {
    ctx := context.TODO()
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }
    return Client {
        ctx: ctx,
        glue: glue.NewFromConfig(cfg),
        lf: lakeformation.NewFromConfig(cfg),
        sts: sts.NewFromConfig(cfg),

        S3: NewS3Client(cfg, ctx),
    }
}

func (c *Client) FetchCurrentAWSAccountId() string {
    input := sts.GetCallerIdentityInput{}
    output, err := c.sts.GetCallerIdentity(c.ctx, &input)
    if err != nil {
        log.Fatalf("unable to get caller identity: %v", err)
    }
    return *output.Account
}

func (c *Client) FetchDatabaseRows(nextToken *string) ([]table.Row, *string, error) {
    input := glue.GetDatabasesInput{
        ResourceShareType: gluetypes.ResourceShareTypeAll,
        MaxResults: aws.Int32(MAX_RESULTS),
    }
    if nextToken != nil {
        input.NextToken = nextToken
    }

    output, err := c.glue.GetDatabases(c.ctx, &input)
    if err != nil {
        log.Fatalf("unable to get databases: %v", err)
    }

    var rows []table.Row
    for _, d := range output.DatabaseList {
        row := table.Row{
            *d.Name,
            *d.CatalogId,
            aws.ToString(d.Description),
            aws.ToString(d.LocationUri),
        }
        if (d.TargetDatabase != nil) {
            row = append(row, aws.ToString(d.TargetDatabase.DatabaseName), aws.ToString(d.TargetDatabase.CatalogId))
        } else {
            row = append(row, "", "")
        }
    
        rows = append(rows, row)
    }
    return rows, output.NextToken, nil
}

func (c *Client) FetchTableRows(nextToken *string) ([]table.Row, *string, error) {
    input := glue.SearchTablesInput{
        MaxResults: aws.Int32(MAX_RESULTS),
        ResourceShareType: gluetypes.ResourceShareTypeAll,
        SortCriteria: []gluetypes.SortCriterion{
            {
                FieldName: aws.String("UpdateTime"),
                Sort: gluetypes.SortDescending,
            },
        },
    }
    if nextToken != nil {
        input.NextToken = nextToken
    }

    output, err := c.glue.SearchTables(c.ctx, &input)
    if err != nil {
        log.Fatalf("unable to get tables: %v", err)
    }

    var rows []table.Row
    for _, t := range output.TableList {
        row := table.Row{
            *t.Name,
            *t.DatabaseName,
            *t.CatalogId,
            aws.ToString(t.Description),
        }

        if (t.StorageDescriptor != nil) {
            row = append(row, aws.ToString(t.StorageDescriptor.Location))
        } else {
            row = append(row, "")
        }

        if (t.TargetTable != nil) {
            row = append(row, aws.ToString(t.TargetTable.DatabaseName), aws.ToString(t.TargetTable.CatalogId))
        } else {
            row = append(row, "", "")
        }
    
        rows = append(rows, row)
    }
    return rows, output.NextToken, nil
}

func (c *Client) FetchLFTagRows(nextToken *string) ([]table.Row, *string, error) {
    input := lakeformation.ListLFTagsInput{
        ResourceShareType: lftypes.ResourceShareTypeAll,
        MaxResults: aws.Int32(MAX_RESULTS),
    }
    if nextToken != nil {
        input.NextToken = nextToken
    }
    
    output, err := c.lf.ListLFTags(c.ctx, &input)
    if err != nil {
        log.Fatalf("unable to get tables: %v", err)
    }
    var rows []table.Row
    for _, t := range output.LFTags {
        row := table.Row{
            *t.TagKey,
            strings.Join(t.TagValues, ","),
            *t.CatalogId,
        }	
        rows = append(rows, row)
    }
    return rows, output.NextToken, nil
}

func (c *Client) FetchLFTagPermissionRows(nextToken *string) ([]table.Row, *string, error) {
    input := lakeformation.ListPermissionsInput{
        ResourceType: lftypes.DataLakeResourceTypeLfTag,
        MaxResults: aws.Int32(MAX_RESULTS),
    }
    if nextToken != nil {
        input.NextToken = nextToken
    }
    
    output, err := c.lf.ListPermissions(c.ctx, &input)
    if err != nil {
        log.Fatalf("unable to get LF-Tag Permissions: %v", err)
    }
    var rows []table.Row
    for _, p := range output.PrincipalResourcePermissions {

        var perms strings.Builder
        for _, p := range p.Permissions {
            perms.WriteString(fmt.Sprintf("%s", p))
        }

        var grantable strings.Builder
        for _, p := range p.PermissionsWithGrantOption {
            grantable.WriteString(fmt.Sprintf("%s", p))
        }

        row := table.Row{
            *p.Principal.DataLakePrincipalIdentifier,
            *p.Resource.LFTag.TagKey,
            strings.Join(p.Resource.LFTag.TagValues, ","),
            perms.String(),
            grantable.String(),
        }
        rows = append(rows, row)
    }
    return rows, output.NextToken, nil
}

func (c *Client) FetchDataLakeLocationRows(nextToken *string) ([]table.Row, *string, error) {
    input := lakeformation.ListResourcesInput{
        MaxResults: aws.Int32(MAX_RESULTS),
    }
    if nextToken != nil {
        input.NextToken = nextToken
    }

    output, err := c.lf.ListResources(c.ctx, &input)
    if err != nil {
        log.Fatalf("unable to get LF resources: %v", err)
    }

    var rows []table.Row
    for _, r := range output.ResourceInfoList {
        rows = append(rows, []string{*r.ResourceArn, r.LastModified.String()})
    }
    return rows, output.NextToken, nil
}

func (c *Client) FetchDatabaseDetailRows(databaseName string) ([][]table.Row, error) {
    input := glue.GetDatabaseInput{
        Name: aws.String(databaseName),
    }
    output, err := c.glue.GetDatabase(c.ctx, &input)
    if err != nil {
        log.Fatalf("unable to get database: %v", err)
    }

    var detailRows []table.Row
    detailRows = append(detailRows, []string{"Name", utils.UseStr(output.Database.Name)})
    detailRows = append(detailRows, []string{"Location", utils.UseStr(output.Database.LocationUri)})
    detailRows = append(detailRows, []string{"Description", utils.UseStr(output.Database.Description)})

    
    lfTagInput := lakeformation.GetResourceLFTagsInput{
        Resource: &lftypes.Resource{
            Database: &lftypes.DatabaseResource{
                Name: aws.String(databaseName),
            },
        },
        ShowAssignedLFTags: utils.BoolPtr(true),
    }
    lfTagsOutput, err := c.lf.GetResourceLFTags(c.ctx, &lfTagInput)
    if err != nil {
        log.Fatalf("unable to get database: %v", err)
    }
    var lfTagRows []table.Row
    for _, t := range lfTagsOutput.LFTagOnDatabase {
        lfTagRows = append(lfTagRows, []string{utils.UseStr(t.TagKey), strings.Join(t.TagValues, ",")})
    }

    return [][]table.Row{detailRows, lfTagRows}, nil
}


func (c *Client) FetchTableDetailRows(tableName string, databaseName string) ([][]table.Row, error) {
    input := glue.GetTableInput{
        DatabaseName: aws.String(databaseName),
        Name:         aws.String(tableName),
    }
    output, err := c.glue.GetTable(c.ctx, &input)
    if err != nil {
        log.Fatalf("unable to get table: %v", err)
    }

    var detailRows []table.Row
    detailRows = append(detailRows, []string{"Table Name", utils.UseStr(output.Table.Name)})
    detailRows = append(detailRows, []string{"Database Name", utils.UseStr(output.Table.DatabaseName)})
    detailRows = append(detailRows, []string{"Location", utils.UseStr(output.Table.StorageDescriptor.Location)})
    detailRows = append(detailRows, []string{"Description", utils.UseStr(output.Table.Description)})
    detailRows = append(detailRows, []string{"Last Updated", output.Table.UpdateTime.String()})

    var schemaRows []table.Row
    for i, c := range output.Table.StorageDescriptor.Columns {
        schemaRows = append(schemaRows, []string{strconv.Itoa(i+1), utils.UseStr(c.Name), utils.UseStr(c.Type)})
    }

    lfTagInput := lakeformation.GetResourceLFTagsInput{
        Resource: &lftypes.Resource{
            Table: &lftypes.TableResource{
                DatabaseName: aws.String(databaseName),
                Name: aws.String(tableName),
            },
        },
        ShowAssignedLFTags: utils.BoolPtr(true),
    }
    lfTagsOutput, err := c.lf.GetResourceLFTags(c.ctx, &lfTagInput)
    if err != nil {
        log.Fatalf("unable to get database: %v", err)
    }
    var lfTagRows []table.Row
    for _, t := range lfTagsOutput.LFTagsOnTable {
        lfTagRows = append(lfTagRows, []string{utils.UseStr(t.TagKey), strings.Join(t.TagValues, ",")})
    }

    return [][]table.Row{detailRows, schemaRows, lfTagRows}, nil
}
