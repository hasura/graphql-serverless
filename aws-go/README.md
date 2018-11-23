# Go + AWS Lambda GraphQL Boilerplates

These are GraphQL backend boilerplates in Go that can be deployed on AWS Lambda.

## Frameworks/Libraries

You can choose a boilerplate for your favourite GraphQL framework and ORM from below:

[graphqlgo-gorm](graphqlgo-gorm)

## Schema

We consider a bank account schema where a user can transfer money to another user. This will involve writing a `transfer` resolver which does complex business logic in a transaction.

```
type User {
  id:       Int
  name:     String
  balance:  Int
}

type Query {
  users:  [User]
}

type Mutation {
  addUser(name: String, balance: Int): User
  transfer(userIdFrom: Int, userIdTo: Int, amount: Int): User
}
```

## Development

These boilerplates do not require any additional tooling for local development. Just set the environment variable `LAMBDA_EXECUTION_ENVIRONMENT=local` and you are good to go.

Detailed steps are available in individual readmes.

## Deployment

These boilerplates do not require any additional tooling for deployment. We will use the AWS console to upload and expose the service.

Detailed steps are available in individual readmes.

## Connection Pooling

As discussed in the main [readme](../README.md), without connection pooling our GraphQL backend will not scale at the same rate as serverless invocations. With Postgres, we can add a standalone connection pooler like [pgBouncer](https://pgbouncer.github.io/) to accomplish this.

Detailed steps are available in individual readmes.
