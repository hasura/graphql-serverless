# nodeJS + AWS Lambda GraphQL Boilerplate

This is a GraphQL backend boilerplate in nodeJS that can be deployed on AWS Lambda.

## Stack

node 8.10

AWS RDS Postgres

AWS Lambda

#### Frameworks/Libraries

Apollo Server (GraphQL framework)

Sequelize (Postgres ORM)

netlify-lambda (Local dev)

## Schema

We consider an Author/Article schema where an author can have many articles.

```
type Author {
  id:       Int
  name:     String
  articles: [Article]
}

type Article {
  id:        Int
  title:     String
  content:   String
  author_id: Int
}

type Query {
  authors:  [Author]
  articles: [Article]
}

type Mutation {
  addAuthor(name: String): Author
  addArticle(title: String, content: String, author_id: Int): Article
}
```

## Development

The sample source code is present in `src/graphql.js`. You can start a development environment by using `netlify-lambda`. Make sure you are in the `aws-lambda-nodejs` folder.

```bash
$ pwd
/home/tselvan/graphql-serverless/aws-lambda-nodejs
```

1) First, let's set the environment variable for connecting to the postgres instance on RDS. You can find this endpoint on your RDS instances page on AWS console:

```bash
$ export POSTGRES_CONNECTION_STRING='postgres://username:password@rds-database-endpoint.us-east-1.rds.amazonaws.com:5432/mydb' 
```

2) Next, lets create the tables required for our schema. The SQL commands are in `migrations.sql` file.

```bash
$ psql $POSTGRES_CONNECTION_STRING < migrations.sql
```

3) Finally, in the `aws-lambda-nodejs` folder, run `npm install` to install `netlify-lambda`. Then run `netlify-lambda serve`.

```bash
$ npm install

$ npx netlify-lambda -c webpack.config.js serve src

Output:

netlify-lambda: Starting server
Lambda server is listening on 9000
```

This will start a local server on `localhost:9000`. You can hit the graphql service at `localhost:9000/graphql`. This opens a graphql playground where you can query your schema.

Edit the source code as you wish in the `src` folder and it will auto-reload changes.

## Deployment

Now that you have run the graphql service locally and made any required changes, it's time to deploy your service to AWS Lambda and get an endpoint. The easiest way to do this is through the AWS console.

1) Create a Lambda function by clicking on Create Function on your Lambda console. Choose the `NodeJS 8.10` runtime and `lambda_basic_execution` role.

![create-lambda](assets/create-lambda.png)

2) In the next page (or Lambda instance page), select API Gateway as the trigger.

![create-api-gateway](assets/create-api-gateway.png)

3) Configure the API Gateway as you wish. The simplest configuration is shown below.

![configure-api-gateway](assets/configure-api-gateway.png)

Save your changes. You will receive a HTTPS endpoint for your lambda.

![output-api-gateway](assets/output-api-gateway.png)

If you go to the endpoint, you will receive a "Hello from Lambda!" message. This is because we haven't uploaded any code yet!

4) Upload code: You can use `netlify-lambda` to create a buildpack for your lambda.

```bash
$ npx netlify-lambda -c webpack.config.js build src
```

This will create a `graphql.js` file in your `build/` folder. You can either copy the contents of this file directly or zip it and upload it on the Lambda console.

Also, make sure to add the `POSTGRES_CONNECTION_STRING` environment variable.

![upload-code](assets/upload-code.png)

And that's it. Hit save and visit the endpoint again. You will see the graphql playground again.

NOTE: You may have to edit the GraphQL URL in the playground to reflect the right endpoint ( same as the URL in the browser ).

## Connection Pooling

As discussed in the main [readme](../README.md), without connection pooling our GraphQL backend will not scale at the same rate as serverless invocations. With Postgres, we can add a standalone connection pooler like [pgBouncer](https://pgbouncer.github.io/) to accomplish this. 

Deploying pgBouncer requires an EC2 instance. We can use the CloudFormation template present in this folder: [cloudformation.json](cloudformation.json) to deploy a pgBouncer EC2 instance in few clicks.

#### Deploy pgBouncer

1. Goto CloudFormation in AWS Console and select Create Stack.

2. Upload the file [cloudformation.json](cloudformation.json) as the template.

3. In the next step, fill in your Postgres connection details:

![cloudformation-params](assets/cloudformation-params.png)

4. You do not need any other configuration, so just continue by clicking NEXT and finally click CREATE.

5. After the creation is complete, you will see your new `POSTGRES_CONNECTION_STRING` in the output:

![cloudformation-output](assets/cloudformation-output.png)

Now, change your `POSTGRES_CONNECTION_STRING` in your lambda function to the new value. And, everything should just work!

#### Results

Using pgBouncer, here are the results for corresponding rate of lambda invocations. The tests were conducted with the `addAuthor` mutation using [jmeter](https://jmeter.apache.org/).

|  Error Rate -> | Without pgBouncer | With pgBouncer|
| -------------- | ----------------- | ------------- |
| 100 req/s      | 86%               | 0%            |
| 1000 req/s     | 92%               | 4%            |
| 10000 req/s    | NA                | 3%            |

