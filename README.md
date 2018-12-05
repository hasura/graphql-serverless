# GraphQL Serverless

GraphQL backend boilerplates that can be deployed on serverless platforms.

In theory, GraphQL and Serverless are supposed to work well together. Serverless gives a scalable and no-ops platform to deploy a GraphQL backend instantly. Although in practice, there are few reasons why this may not work: 

1) **Serverless is new**:

Development, testing and deployment mechanisms on serverless are not mature. Although there are few tools out there which ease out some of the processes, these tools are rapidly changing.

2) **Managing state in serverless**:

Serverless backends cannot hold state between different requests. This means the state must be created and destroyed in each serverless request which becomes a bottleneck. In the case of GraphQL backend, a database connection must be created in each request which quickly exhausts the database limits.

The only way to solve these problems are:

1) Setup a simple development/deployment workflow on serverless which is reliable.
2) Setup a connection pool for managing database connections in serverless backends.

## Getting Started

The boilerplates in this repo provide sample source code for building GraphQL backends while solving the above problems:

Get started with the following languages and serverless platforms:

[NodeJS + AWS Lambda](aws-nodejs/apollo-sequelize)

[Python + AWS Lambda](aws-python/graphene-sqlalchemy)

[Go + AWS Lambda](aws-go/graphqlgo-gorm)

### Using with Hasura GraphQL Engine (Optional)

While you can use these boilerplates to create any kind of GraphQL schema, you may wish to merge your schema with [Hasura GraphQL Engine](https://hasura.io) to augment your schema with a wide range of CRUD APIs.

Follow this guide to merge your schema with Hasura: [using-with-hasura.md](using-with-hasura.md)
