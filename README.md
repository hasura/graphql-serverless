# GraphQL Serverless

GraphQL backend boilerplates that can be deployed on serverless platforms.

In theory, GraphQL and Serverless are supposed to work well together. Serverless gives a scalable platform to deploy a GraphQL backend instantly. Although in practice, there are few reasons why this is not convenient:

1) **Serverless is new**: Developing, testing and deploying on serverless is not mature.
2) **Database connections with serverless**: GraphQL backends need database connections that cannot scale at the same rate as serverless invocations.

The only way to solve these problems are:

1) Setup a simple development/deployment workflow on serverless.
2) Setup a connection pool for the GraphQL backend.

The boilerplates in this repo provide sample source code to solve the above problems.

Get started with the following languages and serverless platforms:


[Go + AWS Lambda](https://github.com/hasura/graphql-serverless/tree/master/aws-lambda-go)

[NodeJS + AWS Lambda](https://github.com/hasura/graphql-serverless/tree/master/aws-lambda-nodejs)

[Python + AWS Lambda](https://github.com/hasura/graphql-serverless/tree/master/aws-lambda-python)

