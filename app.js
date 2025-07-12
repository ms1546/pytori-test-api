import {
    DynamoDBClient,
    GetItemCommand,
    ScanCommand,
  } from "@aws-sdk/client-dynamodb";

  const client = new DynamoDBClient({
    region: "ap-northeast-1",
    endpoint: process.env.DYNAMO_ENDPOINT,
  });

  const getRepoSummary = async (repoId, repoItem) => {
    const repository_name = repoItem.name.S;
    const status = parseInt(repoItem.status.N);

    const mergedCommits = await client.send(new ScanCommand({
      TableName: "game_commits",
      FilterExpression: "repository_id = :repo_id AND is_merged = :merged",
      ExpressionAttributeValues: {
        ":repo_id": { N: repoId.toString() },
        ":merged": { N: "1" },
      },
    }));

    const sorted = mergedCommits.Items.sort((a, b) =>
      new Date(b.merged_on.S) - new Date(a.merged_on.S)
    );
    const latest = sorted[0];

    return {
      repository_id: repoId,
      repository_name,
      status,
      shiritori_count: mergedCommits.Items.length,
      current_word: latest?.current_word?.S || null,
      review_comment: latest?.review_comment?.S || null,
      merged_on: latest?.merged_on?.S || null,
    };
  };

  export const handler = async (event) => {
    const repoIdParam = event.queryStringParameters?.repository_id;

    if (repoIdParam) {
      const repoId = parseInt(repoIdParam);
      const repoRes = await client.send(new GetItemCommand({
        TableName: "game_repos",
        Key: { id: { N: repoId.toString() } },
      }));

      if (!repoRes.Item) {
        return { statusCode: 404, body: JSON.stringify({ error: "Repository not found" }) };
      }

      const summary = await getRepoSummary(repoId, repoRes.Item);
      return { statusCode: 200, body: JSON.stringify(summary) };
    }

    // クエリなし：全リポジトリ取得
    const allRepos = await client.send(new ScanCommand({
      TableName: "game_repos"
    }));

    const summaries = await Promise.all(
      allRepos.Items.map(async (item) => {
        const repoId = parseInt(item.id.N);
        return await getRepoSummary(repoId, item);
      })
    );

    return {
      statusCode: 200,
      body: JSON.stringify(summaries),
    };
  };
