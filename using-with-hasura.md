# Using GraphQL on Serverless with Hasura GraphQL Engine

Hasura GraphQL Engine gives a wide variety of CRUD GraphQL APIs instantly. If you have deployed a GraphQL backend on serverless platforms, you can merge those schemas with Hasura to augment your GraphQL schema.

Your serverless GraphQL backends are referred as "Remote schemas" in Hasura. This is what Hasura running with Remote schemas looks like:

![remote-schemas](_assets/remote-schemas-arch.png)

### Steps

There are broadly 2 steps:

1) [Deploy Hasura GraphQL Engine](https://docs.hasura.io/1.0/graphql/manual/deployment/index.html)
2) [Add Remote Schema](https://docs.hasura.io/1.0/graphql/manual/remote-schemas/index.html)
