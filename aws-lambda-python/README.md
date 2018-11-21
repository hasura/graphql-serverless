# Python GraphQL Server

This is a GraphQL server with a simple article-author schema. It uses the following:

1. Flask
2. Graphene
3. SQL Alcheme
4. Postgres

## Local development

1. Clone the repo

    ```bash
    git clone git@github.com:hasura/graphql-serverless
    cd graphql-serverless
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
    export POSTGRES_CONNECTION_STRING

    # should the database be initialized (0 or 1)
    export DATABASE_INIT=1
    ```

3. Run the server

    ```bash
    ./main.py 
    ```

4. Try out graphql queries at `http://localhost:5000/graphql`

## Deployment

Lets deploy this function to a lambda using [Zappa](www.zappa.io)

1. Configure your amazon credentials. [Install amazon CLI](https://docs.aws.amazon.com/cli/latest/userguide/installing.html) and run this command

    ```
    aws configure
    ```



2. Now, from your python virtual environment, initialize your zappa configuration. Run `zappa init`

    ```
    zappa init
    ```

3. Set `main.app` as the modular path to your app's function when prompted for it

4. Finally deploy the function by running

    ```
    zappa deploy dev
    ```

