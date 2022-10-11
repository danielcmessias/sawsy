package data

import (
	"context"
	"log"
	"net/url"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/danielcmessias/sawsy/ui/components/table"
)

type IAMClient struct {
	ctx context.Context
	iam *iam.Client
}

func NewIAMClient(ctx context.Context, iam *iam.Client) *IAMClient {
	return &IAMClient{
		ctx: ctx,
		iam: iam,
	}
}

func (c *IAMClient) GetUsers(nextToken *string) ([]table.Row, *string, error) {
	input := iam.ListUsersInput{
		Marker: nextToken,
	}
	output, err := c.iam.ListUsers(c.ctx, &input)
	if err != nil {
		log.Fatalf("unable to list users: %v", err)
	}

	var rows []table.Row
	for _, u := range output.Users {
		lastUsed := "None"
		if u.PasswordLastUsed != nil {
			lastUsed = formatTime(u.PasswordLastUsed)
		}

		rows = append(rows, table.Row{
			aws.ToString(u.UserName),
			aws.ToString(u.Arn),
			lastUsed,
			formatTime(u.CreateDate),
		})
	}

	return rows, output.Marker, nil
}

func (c *IAMClient) GetRoles(nextToken *string) ([]table.Row, *string, error) {
	input := iam.ListRolesInput{
		Marker: nextToken,
	}
	output, err := c.iam.ListRoles(c.ctx, &input)
	if err != nil {
		log.Fatalf("unable to list roles: %v", err)
	}

	var rows []table.Row
	for _, r := range output.Roles {
		rows = append(rows, table.Row{
			aws.ToString(r.RoleName),
			aws.ToString(r.Arn),
			formatTime(r.CreateDate),
		})
	}

	return rows, output.Marker, nil
}

func (c *IAMClient) GetUserPolicies(userName string, nextToken *string) ([]table.Row, *string, error) {
	inputAttached := iam.ListAttachedUserPoliciesInput{
		UserName: aws.String(userName),
	}
	outputAttached, err := c.iam.ListAttachedUserPolicies(c.ctx, &inputAttached)
	if err != nil {
		log.Fatalf("unable to list attached user policies: %v", err)
	}

	inputInline := iam.ListUserPoliciesInput{
		UserName: aws.String(userName),
		Marker:   nextToken,
	}
	outputInline, err := c.iam.ListUserPolicies(c.ctx, &inputInline)
	if err != nil {
		log.Fatalf("unable to list user policies: %v", err)
	}

	var rows []table.Row
	for _, p := range outputAttached.AttachedPolicies {
		rows = append(rows, table.Row{
			aws.ToString(p.PolicyName),
			"Attached",
			aws.ToString(p.PolicyArn),
		})
	}
	for _, p := range outputInline.PolicyNames {
		rows = append(rows, table.Row{
			p,
			"Inline",
			"",
		})
	}

	return rows, outputInline.Marker, nil
}

func (c *IAMClient) GetUserTags(userName string, nextToken *string) ([]table.Row, *string, error) {
	input := iam.ListUserTagsInput{
		UserName: aws.String(userName),
		Marker:   nextToken,
	}
	output, err := c.iam.ListUserTags(c.ctx, &input)
	if err != nil {
		log.Fatalf("unable to list user tags: %v", err)
	}

	var rows []table.Row
	for _, t := range output.Tags {
		rows = append(rows, table.Row{
			aws.ToString(t.Key),
			aws.ToString(t.Value),
		})
	}

	return rows, output.Marker, nil
}

func (c *IAMClient) GetRolePolicies(roleName string, nextToken *string) ([]table.Row, *string, error) {
	inputAttached := iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	}
	outputAttached, err := c.iam.ListAttachedRolePolicies(c.ctx, &inputAttached)
	if err != nil {
		log.Fatalf("unable to list attached user policies: %v", err)
	}

	inputInline := iam.ListRolePoliciesInput{
		RoleName: aws.String(roleName),
		Marker:   nextToken,
	}
	outputInline, err := c.iam.ListRolePolicies(c.ctx, &inputInline)
	if err != nil {
		log.Fatalf("unable to list user policies: %v", err)
	}

	var rows []table.Row
	for _, p := range outputAttached.AttachedPolicies {
		rows = append(rows, table.Row{
			aws.ToString(p.PolicyName),
			"Attached",
			aws.ToString(p.PolicyArn),
		})
	}
	for _, p := range outputInline.PolicyNames {
		rows = append(rows, table.Row{
			p,
			"Inline",
			"",
		})
	}

	return rows, outputInline.Marker, nil
}

func (c *IAMClient) GetRoleTags(roleName string, nextToken *string) ([]table.Row, *string, error) {
	input := iam.ListRoleTagsInput{
		RoleName: aws.String(roleName),
		Marker:   nextToken,
	}
	output, err := c.iam.ListRoleTags(c.ctx, &input)
	if err != nil {
		log.Fatalf("unable to list user tags: %v", err)
	}

	var rows []table.Row
	for _, t := range output.Tags {
		rows = append(rows, table.Row{
			aws.ToString(t.Key),
			aws.ToString(t.Value),
		})
	}

	return rows, output.Marker, nil
}

func (c *IAMClient) GetManagedPolicy(policyArn string) (string, error) {
	input := iam.GetPolicyInput{
		PolicyArn: aws.String(policyArn),
	}
	output, err := c.iam.GetPolicy(c.ctx, &input)
	if err != nil {
		log.Fatalf("unable to get IAM policy: %v", err)
	}

	inputVersion := iam.GetPolicyVersionInput{
		PolicyArn: aws.String(policyArn),
		VersionId: output.Policy.DefaultVersionId,
	}
	outputVersion, err := c.iam.GetPolicyVersion(c.ctx, &inputVersion)
	if err != nil {
		log.Fatalf("unable to get IAM policy version: %v", err)
	}

	decodedDocument, err := url.QueryUnescape(aws.ToString(outputVersion.PolicyVersion.Document))
	if err != nil {
		log.Fatalf("unable to URL decode IAM policy document: %v", err)
	}

	return formatJson(decodedDocument), nil
}

func (c *IAMClient) GetInlineUserPolicy(userName string, policyName string) (string, error) {
	input := iam.GetUserPolicyInput{
		UserName:   aws.String(userName),
		PolicyName: aws.String(policyName),
	}
	output, err := c.iam.GetUserPolicy(c.ctx, &input)
	if err != nil {
		log.Fatalf("unable to get IAM user policy: %v", err)
	}

	decodedDocument, err := url.QueryUnescape(aws.ToString(output.PolicyDocument))
	if err != nil {
		log.Fatalf("unable to URL decode IAM policy document: %v", err)
	}

	return formatJson(decodedDocument), nil
}

func (c *IAMClient) GetInlineRolePolicy(roleName string, policyName string) (string, error) {
	input := iam.GetRolePolicyInput{
		RoleName:   aws.String(roleName),
		PolicyName: aws.String(policyName),
	}
	output, err := c.iam.GetRolePolicy(c.ctx, &input)
	if err != nil {
		log.Fatalf("unable to get IAM user policy: %v", err)
	}

	decodedDocument, err := url.QueryUnescape(aws.ToString(output.PolicyDocument))
	if err != nil {
		log.Fatalf("unable to URL decode IAM policy document: %v", err)
	}

	return formatJson(decodedDocument), nil
}

func (c *IAMClient) GetAssumeRolePolicy(roleName string) (string, error) {
	input := iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}
	output, err := c.iam.GetRole(c.ctx, &input)
	if err != nil {
		log.Fatalf("unable to get IAM assume role policy: %v", err)
	}

	decodedDocument, err := url.QueryUnescape(aws.ToString(output.Role.AssumeRolePolicyDocument))
	if err != nil {
		log.Fatalf("unable to URL decode IAM policy document: %v", err)
	}

	return formatJson(decodedDocument), nil
}
