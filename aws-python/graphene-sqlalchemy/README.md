# Python + AWS Lambda + Graphene + SQLAlchemy

This is a GraphQL backend boilerplate in Python that can be deployed on AWS Lambda.

## Stack

Python 2.7

AWS RDS Postgres

AWS Lambda

#### Frameworks/Libraries

Graphene

SQL Alchemy (Postgres ORM)

## Schema

We consider an Author/Article schema where an author can have many articles.

```
type User {
  id:       Int
  name:     String
  balance:  Int
}

type Query {
  users:  [User]
  hello: String
}

type Mutation {
  addUser(name: String, balance: Int): User
  transfer(userIdFrom: Int, userIdTo: Int, amount: Int): User
}
```

## Development

1. Clone the repo

    ```bash
    git clone git@github.com:hasura/graphql-serverless
    cd graphql-serverless/aws-python/graphene-sqlalchemy
    ```

2. Set up your development environment
    ```bash
    pip install virtualenv
    virtualenv env
    source env/bin/activate
    pip install -r requirements.txt
    ```

3. Set required environment variables

    ```bash
    # your postgres connection string
    export POSTGRES_CONNECTION_STRING='postgres://username:password@rds-database-endpoint.us-east-1.rds.amazonaws.com:5432/mydb' 
    ```


4. Next, lets create the tables required for our schema. The SQL commands are in `migrations.sql` file.

    ```bash
    $ psql $POSTGRES_CONNECTION_STRING < ../../schema/migrations.sql
    ```

5. Run the server with an environment variable `LAMBDA_LOCAL_DEVELOPMENT=1`

    ```bash
    LAMBDA_LOCAL_DEVELOPMENT=1 ./main.py 
    ```

6. Try out graphql queries at `http://localhost:5000/graphql`

7. (Optional) Open GraphiQL using [graphQurl](https://github.com/hasura/graphqurl). `graphQurl` gives a local graphiQL environment for any graphql endpoint.

## Deployment

Now that you have run the graphql service locally and made any required changes, it's time to deploy your service to AWS Lambda and get an endpoint. The easiest way to do this is through the AWS console.

1) Create a Lambda function by clicking on Create Function on your [Lambda console](https://console.aws.amazon.com/lambda/home). Choose the `Python 2.7` runtime.

2) In the next page (or Lambda instance page), select API Gateway as the trigger.

   ![create-api-gateway](../../_assets/create-api-gateway.png)

3) Configure the API Gateway as you wish. The simplest configuration is shown below.

   ![configure-api-gateway](../../_assets/configure-api-gateway.png)

Save your changes. You will receive a HTTPS endpoint for your lambda.

   ![output-api-gateway](../../_assets/output-api-gateway.png)

If you go to the endpoint, you will receive a "Hello from Lambda!" message. This is because we haven't uploaded any code yet!

4) Zip and upload code (follow the instructions [here](https://docs.aws.amazon.com/lambda/latest/dg/lambda-python-how-to-create-deployment-package.html#python-package-venv)). You can also run the shell script to genereate a zipped package (`_package.zip`)

   ```bash
   $ ./package.sh
   ```

5. Make sure to add the `POSTGRES_CONNECTION_STRING` environment variable.

6. Finally upload the code while setting the handler as `main.lambda_handler`.

## Connection Pooling

As discussed in the main [readme](../../README.md), without connection pooling our GraphQL backend will not scale at the same rate as serverless invocations. With Postgres, we can add a standalone connection pooler like [pgBouncer](https://pgbouncer.github.io/) to accomplish this. 

Deploying pgBouncer requires an EC2 instance. We can use the CloudFormation template present in this folder: [cloudformation.json](../../cloudformation/cloudformation.json) to deploy a pgBouncer EC2 instance in few clicks.

#### Deploy pgBouncer

1. Goto CloudFormation in AWS Console and select Create Stack.

2. Upload the file [cloudformation.json](../../cloudformation/cloudformation.json) as the template.

3. In the next step, fill in your Postgres connection details:

![cloudformation-params](../../_assets/cloudformation-params.png)

4. You do not need any other configuration, so just continue by clicking NEXT and finally click CREATE.

5. After the creation is complete, you will see your new `POSTGRES_CONNECTION_STRING` in the output:

![cloudformation-output](../../_assets/cloudformation-output.png)

Now, change your `POSTGRES_CONNECTION_STRING` in your `zappa_settings.json` to the new value. Update the environment variable by running `zappa update` and, everything should just work!

#### Results

Using pgBouncer, here are the results for corresponding rate of lambda invocations. The tests were conducted with the `addAuthor` mutation using [jmeter](https://jmeter.apache.org/).

|  Error Rate -> | Without pgBouncer | With pgBouncer|
| -------------- | ----------------- | ------------- |
| 100 req/s      | 86%               | 0%            |
| 1000 req/s     | 92%               | 4%            |
| 10000 req/s    | NA                | 3%            |

