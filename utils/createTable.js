// createTable.js
import {
    DynamoDBClient,
    CreateTableCommand,
    ListTablesCommand,
  } from "@aws-sdk/client-dynamodb";

  const client = new DynamoDBClient({
    region: "ap-northeast-1",
    endpoint: "http://localhost:8000",
  });

  async function createTableIfNotExists(params) {
    const { TableName } = params;

    const existingTables = await client.send(new ListTablesCommand({}));
    if (existingTables.TableNames.includes(TableName)) {
      console.log(`✅ ${TableName} は既に存在します`);
      return;
    }

    await client.send(new CreateTableCommand(params));
    console.log(`✅ ${TableName} を作成しました`);
  }

  await createTableIfNotExists({
    TableName: "game_commits",
    KeySchema: [{ AttributeName: "id", KeyType: "HASH" }],
    AttributeDefinitions: [{ AttributeName: "id", AttributeType: "N" }],
    ProvisionedThroughput: { ReadCapacityUnits: 5, WriteCapacityUnits: 5 },
  });

  await createTableIfNotExists({
    TableName: "game_repos",
    KeySchema: [{ AttributeName: "id", KeyType: "HASH" }],
    AttributeDefinitions: [{ AttributeName: "id", AttributeType: "N" }],
    ProvisionedThroughput: { ReadCapacityUnits: 5, WriteCapacityUnits: 5 },
  });
