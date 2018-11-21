WIP: The serverless part is a work in progress.

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
    ./app.py 
    ```

4. Try out graphql queries at `http://localhost:5000/graphql`
