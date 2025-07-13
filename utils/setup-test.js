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
      id: { N: "101" },
      name: { S: "team-a" },
      status: { N: "1" }
    }
  }));

  await client.send(new PutItemCommand({
    TableName: "game_commits",
    Item: {
      id: { N: "1" },
      repository_id: { N: "101" },
      review_comment: { S: "ちょーすごい" },
      current_word: { S: "ぬいぐるみ" },
      theme: { S: "春" },
      is_merged: { N: "1" },
      merged_on: { S: "2025-07-10T15:20:00Z" }
    }
  }));

  console.log("✅ テストデータを投入しました");
};

setup();
