package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	ctx := context.TODO()
	endpoint := os.Getenv("DYNAMO_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("ap-northeast-1"),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint, HostnameImmutable: true}, nil
			},
		)),
	)
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	client := ddb.NewFromConfig(cfg)

	createTableIfNotExists(ctx, client, &ddb.CreateTableInput{
		TableName: aws.String("pytori_commits"),
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeN},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})

	createTableIfNotExists(ctx, client, &ddb.CreateTableInput{
		TableName: aws.String("pytori_repos"),
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeN},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})

	putTestData(ctx, client)
}

func createTableIfNotExists(ctx context.Context, client *ddb.Client, input *ddb.CreateTableInput) {
	tables, err := client.ListTables(ctx, &ddb.ListTablesInput{})
	if err != nil {
		log.Printf("❌ ListTables failed: %v", err)
		return
	}
	for _, name := range tables.TableNames {
		if name == *input.TableName {
			fmt.Printf("✅ %s は既に存在します\n", *input.TableName)
			return
		}
	}
	_, err = client.CreateTable(ctx, input)
	if err != nil {
		log.Printf("❌ %s の作成に失敗しました: %v", *input.TableName, err)
	} else {
		fmt.Printf("✅ %s を作成しました\n", *input.TableName)
	}
}

func putTestData(ctx context.Context, client *ddb.Client) {
	repos := []map[string]types.AttributeValue{
		{"id": &types.AttributeValueMemberN{Value: "101"}, "name": &types.AttributeValueMemberS{Value: "team-a"}, "status": &types.AttributeValueMemberN{Value: "1"}},
		{"id": &types.AttributeValueMemberN{Value: "102"}, "name": &types.AttributeValueMemberS{Value: "team-b"}, "status": &types.AttributeValueMemberN{Value: "1"}},
	}

	commits := []map[string]types.AttributeValue{
		{
			"id": &types.AttributeValueMemberN{Value: "1"}, "repository_id": &types.AttributeValueMemberN{Value: "101"},
			"review_comment": &types.AttributeValueMemberS{Value: "ちょーすごい"}, "current_word": &types.AttributeValueMemberS{Value: "def"},
			"theme": &types.AttributeValueMemberS{Value: "春"}, "is_merged": &types.AttributeValueMemberN{Value: "1"},
			"merged_on": &types.AttributeValueMemberS{Value: "2025-07-10T15:20:00Z"},
		},
		{
			"id": &types.AttributeValueMemberN{Value: "2"}, "repository_id": &types.AttributeValueMemberN{Value: "102"},
			"review_comment": &types.AttributeValueMemberS{Value: "ナイスコミット！"}, "current_word": &types.AttributeValueMemberS{Value: "eval"},
			"theme": &types.AttributeValueMemberS{Value: "おやつ"}, "is_merged": &types.AttributeValueMemberN{Value: "1"},
			"merged_on": &types.AttributeValueMemberS{Value: "2025-07-11T11:45:00Z"},
		},
		{
			"id": &types.AttributeValueMemberN{Value: "3"}, "repository_id": &types.AttributeValueMemberN{Value: "102"},
			"review_comment": &types.AttributeValueMemberS{Value: "ナイス!!"}, "current_word": &types.AttributeValueMemberS{Value: "list"},
			"theme": &types.AttributeValueMemberS{Value: "冬"}, "is_merged": &types.AttributeValueMemberN{Value: "1"},
			"merged_on": &types.AttributeValueMemberS{Value: "2025-07-12T11:45:00Z"},
		},
	}

	for _, item := range repos {
		_, err := client.PutItem(ctx, &ddb.PutItemInput{
			TableName: aws.String("pytori_repos"),
			Item:      item,
		})
		if err != nil {
			log.Printf("❌ Insert repo failed: %v", err)
		}
	}

	for _, item := range commits {
		_, err := client.PutItem(ctx, &ddb.PutItemInput{
			TableName: aws.String("pytori_commits"),
			Item:      item,
		})
		if err != nil {
			log.Printf("❌ Insert commit failed: %v", err)
		}
	}
	fmt.Println("✅ テストデータを投入しました")
}
