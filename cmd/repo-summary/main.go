package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var client *ddb.Client

type RepoSummary struct {
	RepositoryName string `json:"repository_name"`
	Status         int    `json:"status"`
	CurrentWord    string `json:"current_word"`
	MergedOn       string `json:"merged_on"`
}

func getString(attr types.AttributeValue) string {
	if v, ok := attr.(*types.AttributeValueMemberS); ok {
		return v.Value
	}
	return ""
}

func getInt(attr types.AttributeValue) int {
	if v, ok := attr.(*types.AttributeValueMemberN); ok {
		i, _ := strconv.Atoi(v.Value)
		return i
	}
	return 0
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	repoName := req.QueryStringParameters["repository_name"]
	var summaries []RepoSummary

	if repoName != "" {
		repo, err := client.GetItem(ctx, &ddb.GetItemInput{
			TableName: aws.String("pytori_repos"),
			Key: map[string]types.AttributeValue{
				"name": &types.AttributeValueMemberS{Value: repoName},
			},
		})
		if err != nil || repo.Item == nil {
			return events.APIGatewayProxyResponse{StatusCode: 404, Body: `{"error":"Repository not found"}`}, nil
		}
		status := getInt(repo.Item["status"])

		out, err := client.Query(ctx, &ddb.QueryInput{
			TableName:              aws.String("ShiritoriMergedWords"),
			IndexName:              aws.String("repository_name-index"),
			KeyConditionExpression: aws.String("repository_name = :name"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":name": &types.AttributeValueMemberS{Value: repoName},
			},
		})
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"error":"Failed to query"}`}, nil
		}

		sort.Slice(out.Items, func(i, j int) bool {
			return getString(out.Items[i]["merged_on"]) > getString(out.Items[j]["merged_on"])
		})

		for _, item := range out.Items {
			summaries = append(summaries, RepoSummary{
				RepositoryName: repoName,
				Status:         status,
				CurrentWord:    getString(item["current_word"]),
				MergedOn:       getString(item["merged_on"]),
			})
		}
	} else {
		allRepos, err := client.Scan(ctx, &ddb.ScanInput{TableName: aws.String("pytori_repos")})
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"error":"Failed to scan repos"}`}, nil
		}

		for _, repo := range allRepos.Items {
			repoName := getString(repo["name"])
			status := getInt(repo["status"])

			out, err := client.Query(ctx, &ddb.QueryInput{
				TableName:              aws.String("ShiritoriMergedWords"),
				IndexName:              aws.String("repository_name-index"),
				KeyConditionExpression: aws.String("repository_name = :name"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":name": &types.AttributeValueMemberS{Value: repoName},
				},
			})
			if err != nil {
				continue
			}

			sort.Slice(out.Items, func(i, j int) bool {
				return getString(out.Items[i]["merged_on"]) > getString(out.Items[j]["merged_on"])
			})

			for _, item := range out.Items {
				summaries = append(summaries, RepoSummary{
					RepositoryName: repoName,
					Status:         status,
					CurrentWord:    getString(item["current_word"]),
					MergedOn:       getString(item["merged_on"]),
				})
			}
		}
	}

	body, _ := json.Marshal(summaries)
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(body)}, nil
}

func main() {
	ctx := context.TODO()
	var cfg aws.Config
	var err error

	endpoint := os.Getenv("DYNAMO_ENDPOINT")
	if endpoint != "" {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion("ap-northeast-1"),
			config.WithEndpointResolver(aws.EndpointResolverFunc(
				func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{URL: endpoint, HostnameImmutable: true}, nil
				})),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
	}

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client = ddb.NewFromConfig(cfg)
	lambda.Start(handler)
}
