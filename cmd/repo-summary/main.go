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
	RepositoryID   int    `json:"repository_id"`
	RepositoryName string `json:"repository_name"`
	Status         int    `json:"status"`
	ShiritoriCount int    `json:"shiritori_count"`
	CurrentWord    string `json:"current_word"`
	ReviewComment  string `json:"review_comment"`
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

func getRepoSummary(ctx context.Context, repoId int, repoItem map[string]types.AttributeValue) (*RepoSummary, error) {
	repoName := getString(repoItem["name"])
	status := getInt(repoItem["status"])

	out, err := client.Scan(ctx, &ddb.ScanInput{
		TableName:        aws.String("pytori_commits"),
		FilterExpression: aws.String("repository_id = :repo_id AND is_merged = :merged"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":repo_id": &types.AttributeValueMemberN{Value: strconv.Itoa(repoId)},
			":merged":  &types.AttributeValueMemberN{Value: "1"},
		},
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(out.Items, func(i, j int) bool {
		return getString(out.Items[i]["merged_on"]) > getString(out.Items[j]["merged_on"])
	})
	latest := map[string]types.AttributeValue{}
	if len(out.Items) > 0 {
		latest = out.Items[0]
	}

	return &RepoSummary{
		RepositoryID:   repoId,
		RepositoryName: repoName,
		Status:         status,
		ShiritoriCount: len(out.Items),
		CurrentWord:    getString(latest["current_word"]),
		ReviewComment:  getString(latest["review_comment"]),
		MergedOn:       getString(latest["merged_on"]),
	}, nil
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	repoIdParam := req.QueryStringParameters["repository_id"]

	if repoIdParam != "" {
		repoId, _ := strconv.Atoi(repoIdParam)
		repo, err := client.GetItem(ctx, &ddb.GetItemInput{
			TableName: aws.String("pytori_repos"),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberN{Value: strconv.Itoa(repoId)},
			},
		})
		if err != nil || repo.Item == nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 404,
				Body:       `{"error":"Repository not found"}`,
			}, nil
		}

		summary, err := getRepoSummary(ctx, repoId, repo.Item)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"error":"Failed to get summary"}`,
			}, nil
		}
		body, _ := json.Marshal(summary)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(body),
		}, nil
	}

	allRepos, err := client.Scan(ctx, &ddb.ScanInput{
		TableName: aws.String("pytori_repos"),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error":"Failed to get repositories"}`,
		}, nil
	}

	var summaries []RepoSummary
	for _, item := range allRepos.Items {
		repoId := getInt(item["id"])
		summary, err := getRepoSummary(ctx, repoId, item)
		if err == nil {
			summaries = append(summaries, *summary)
		}
	}
	body, _ := json.Marshal(summaries)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}

func main() {
	endpoint := os.Getenv("DYNAMO_ENDPOINT")
	if endpoint == "" {
		log.Fatal("DYNAMO_ENDPOINT is not set")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-northeast-1"),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:               endpoint,
					HostnameImmutable: true,
				}, nil
			}),
		),
	)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	client = ddb.NewFromConfig(cfg)

	lambda.Start(handler)
}
