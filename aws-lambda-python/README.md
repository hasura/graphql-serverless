# Python + AWS Lambda GraphQL Boilerplate

This is a GraphQL backend boilerplate in Python that can be deployed on AWS Lambda.

## Stack

Python 2.7

AWS RDS Postgres

AWS Lambda

#### Frameworks/Libraries

Graphene

SQL Alchemy (Postgres ORM)

Zappa (flask-serverless bundling)

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

1. Clone the repo

    ```bash
    git clone git@github.com:hasura/graphql-serverless
    cd graphql-serverless/aws-lambda-python
    ```

1. Set up your development environment
    ```bash
    pip install virtualenv
    virtualenv env
    source env/bin/activate
    pip install -r requirements.txt
    ```

2. Set required environment variables

    ```bash
    # your postgres connection string
    export POSTGRES_CONNECTION_STRING='postgres://username:password@rds-database-endpoint.us-east-1.rds.amazonaws.com:5432/mydb' 
    ```


3. Next, lets create the tables required for our schema. The SQL commands are in `migrations.sql` file.

    ```bash
    $ psql $POSTGRES_CONNECTION_STRING < migrations.sql
    ```

4. Run the server

    ```bash
    ./main.py 
    ```

5. Try out graphql queries at `http://localhost:5000/graphql`

## Deployment

Now that you have run the graphql service locally and made any required changes, it's time to deploy your service to AWS Lambda and get an endpoint. 

Lets deploy this function to a lambda using [Zappa](www.zappa.io)

1. Configure your amazon credentials. [Install amazon CLI](https://docs.aws.amazon.com/cli/latest/userguide/installing.html) and run this command

    ```
    aws configure
    ```



2. Now, from your python virtual environment, initialize your zappa configuration. Run `zappa init`

    ```
    zappa init
    ```

3. Set `main.app` as the modular path to your app's function when prompted for it. Now a `zappa_settings.json` must have been generated.

4. Set your POSTGRES_CONNECTION_STRING as an environment variable in `zappa_settings.json` (you will have to run the `migrations.sql` if you haven't )

    ```json
    {
        "app_function": 'main.app',
        "s3_bucket": '<something>',
        ...
        "environment_variables": {
            "POSTGRES_CONNECTION_STRING": "postgres://username:password@hostname:port/database"
        }
    }
    ```

5. Finally deploy the function by running

    ```
    zappa deploy dev
    ```



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

Now, change your `POSTGRES_CONNECTION_STRING` in your `zappa_settings.json` to the new value. Update the environment variable by running `zappa update` and, everything should just work!

#### Results

Using pgBouncer, here are the results for corresponding rate of lambda invocations. The tests were conducted with the `addAuthor` mutation using [jmeter](https://jmeter.apache.org/).

|  Error Rate -> | Without pgBouncer | With pgBouncer|
| -------------- | ----------------- | ------------- |
| 100 req/s      | 86%               | 0%            |
| 1000 req/s     | 92%               | 4%            |
| 10000 req/s    | NA                | 3%            |

