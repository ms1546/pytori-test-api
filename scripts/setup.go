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

	deleteTableIfExists(ctx, client, "pytori_shiritori")
	deleteTableIfExists(ctx, client, "pytori_repos")

	createTableIfNotExists(ctx, client, &ddb.CreateTableInput{
		TableName: aws.String("pytori_shiritori"),
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("repository_name"), KeyType: types.KeyTypeHash},
			{AttributeName: aws.String("merged_on"), KeyType: types.KeyTypeRange},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("repository_name"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("merged_on"), AttributeType: types.ScalarAttributeTypeS},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})

	createTableIfNotExists(ctx, client, &ddb.CreateTableInput{
		TableName: aws.String("pytori_repos"),
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("name"), KeyType: types.KeyTypeHash},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("name"), AttributeType: types.ScalarAttributeTypeS},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})

	putTestData(ctx, client)
}

func deleteTableIfExists(ctx context.Context, client *ddb.Client, name string) {
	_, err := client.DeleteTable(ctx, &ddb.DeleteTableInput{
		TableName: aws.String(name),
	})
	if err == nil {
		fmt.Printf("🗑️ %s を削除しました\n", name)
	}
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
		{"name": &types.AttributeValueMemberS{Value: "team-a"}, "status": &types.AttributeValueMemberN{Value: "1"}},
		{"name": &types.AttributeValueMemberS{Value: "team-b"}, "status": &types.AttributeValueMemberN{Value: "1"}},
	}

	shiritori := []map[string]types.AttributeValue{
		{"repository_name": &types.AttributeValueMemberS{Value: "team-a"}, "current_word": &types.AttributeValueMemberS{Value: "def"}, "merged_on": &types.AttributeValueMemberS{Value: "2025-07-10T15:20:00Z"}},
		{"repository_name": &types.AttributeValueMemberS{Value: "team-b"}, "current_word": &types.AttributeValueMemberS{Value: "eval"}, "merged_on": &types.AttributeValueMemberS{Value: "2025-07-11T11:45:00Z"}},
		{"repository_name": &types.AttributeValueMemberS{Value: "team-b"}, "current_word": &types.AttributeValueMemberS{Value: "list"}, "merged_on": &types.AttributeValueMemberS{Value: "2025-07-12T11:45:00Z"}},
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

	for _, item := range shiritori {
		_, err := client.PutItem(ctx, &ddb.PutItemInput{
			TableName: aws.String("pytori_shiritori"),
			Item:      item,
		})
		if err != nil {
			log.Printf("❌ Insert shiritori failed: %v", err)
		}
	}
	fmt.Println("✅ テストデータを投入しました")
}
