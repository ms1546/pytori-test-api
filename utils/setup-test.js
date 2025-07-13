import {
  DynamoDBClient,
  PutItemCommand
} from "@aws-sdk/client-dynamodb";

const client = new DynamoDBClient({
  region: "ap-northeast-1",
  endpoint: "http://localhost:8000"
});

const setup = async () => {
  await client.send(new PutItemCommand({
    TableName: "game_repos",
    Item: {
      id:     { N: "101" },
      name:   { S: "team-a" },
      status: { N: "1" }
    }
  }));

  await client.send(new PutItemCommand({
    TableName: "game_commits",
    Item: {
      id:            { N: "1" },
      repository_id: { N: "101" },
      review_comment:{ S: "ちょーすごい" },
      current_word:  { S: "ぬいぐるみ" },
      theme:         { S: "春" },
      is_merged:     { N: "1" },
      merged_on:     { S: "2025-07-10T15:20:00Z" }
    }
  }));

  await client.send(new PutItemCommand({
    TableName: "game_repos",
    Item: {
      id:     { N: "102" },
      name:   { S: "team-b" },
      status: { N: "1" }
    }
  }));

  await client.send(new PutItemCommand({
    TableName: "game_commits",
    Item: {
      id:            { N: "2" },
      repository_id: { N: "102" },
      review_comment:{ S: "ナイスコミット！" },
      current_word:  { S: "みかん" },
      theme:         { S: "テストカバレッジ" },
      is_merged:     { N: "1" },
      merged_on:     { S: "2025-07-11T11:45:00Z" }
    }
  }));

  console.log("✅ テストデータを 2 件投入しました");
};

setup();
