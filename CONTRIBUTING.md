# Contributing

We strongly encourage contributions of boilerplates in other languages, frameworks/ORM or serverless platforms. The guidelines below will help you get started. Please feel free to raise an issue if you'd like to ask any specific question about contributing.

## Structure

The directory structure of the repo is as follows:

```
├── aws-go
│   ├── graphqlgo-gorm
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── handler.go
│   │   └── README.md
│   └── README.md
├── aws-nodejs
```

If you are contributing a new language or platform, then you should create a top level directory as follows:

```
├── aws-go
│   ├── graphqlgo-gorm
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── handler.go
│   │   └── README.md
│   └── README.md
├── <MY-CLOUD>-java
│   ├── <SOME-FRAMEWORK>
│   │   ├── main.java
│   │   └── README.md
│   └── README.md
├── aws-nodejs
```

If you are contributing a new framework/ORM for an existing language + platform combination, then you should create a sub-directory within the <platform>-<language> directory:

```
├── aws-go
│   ├── graphqlgo-gorm
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── handler.go
│   ├── <MY-GO-FRAMEWORK>-<MY-GO-ORM>
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── handler.go
│   └── README.md
├── aws-nodejs
```

## Development Guidelines

1. Use the GraphQL schema as described in the README.
2. Use local development/deployment workflows with no or minimum dependencies. Refer to existing boilerplates for development/deployment methods without any external tools.
3. Use `graphQurl` to launch GraphiQL locally on your schema and test thoroughly.
4. Remember, you are helping the community so follow any best practices for your language/framework.
