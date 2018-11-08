const { ApolloServer, gql } = require('apollo-server-lambda');
const Sequelize = require("sequelize");

const POSTGRES_CONNECTION_STRING = process.env.POSTGRES_CONNECTION_STRING || "postgres://postgres@localhost:5432/postgres"

const sequelize = new Sequelize(
    POSTGRES_CONNECTION_STRING, {}
);

const Author = sequelize.define('author', {
    id: { type: Sequelize.INTEGER, autoIncrement: true, primaryKey: true },
    name: Sequelize.TEXT
},
{
    timestamps: false
});

const Article = sequelize.define('article', {
    id: { type: Sequelize.INTEGER, autoIncrement: true, primaryKey: true },
    title: Sequelize.TEXT,
    content: Sequelize.TEXT,
    author_id: { type: Sequelize.INTEGER, references: { model: Author, key: 'id'} }
},
{
    timestamps: false
});

Author.hasMany(Article, { foreignKey: 'author_id' });
// Type definitions define the "shape" of your data and specify
// which ways the data can be fetched from the GraphQL server.
const typeDefs = gql`
  type Author {
    id: Int
    name: String
    articles: [Article]
  }

  type Article {
    id: Int
    title: String
    content: String
    author_id: Int
  }

  type Query {
    authors: [Author]
    articles: [Article]
  }

  type Mutation {
    addAuthor(name: String): Author
    addArticle(title: String, content: String, author_id: Int): Article
  }
`;

// Resolvers define the technique for fetching the types in the
// schema.  We'll retrieve books from the "books" array above.
const resolvers = {
    Query: {
        authors: () => Author.findAll({ include: [Article] }),
        articles: () => Article.findAll()
    },
    Mutation: {
        addAuthor: async (_, { name }) => {
            try {
                const author = await Author.create({
                    name: name
                });

                return author;
            } catch (e) {
                throw new Error(e);
            }
        },
        addArticle: async (_, { title, content, author_id}) => {
            try {
                const article = await Article.create({
                    title: title,
                    content: content,
                    author_id: author_id
                });
                return article;
            } catch (e) {
                throw new Error(e);
            }
        }
    }
};

// In the most basic sense, the ApolloServer can be started
// by passing type definitions (typeDefs) and the resolvers
// responsible for fetching the data for those types.
const server = new ApolloServer({
    typeDefs,
    resolvers,
    context: ({ event, context }) => ({
        headers: event.headers,
        functionName: context.functionName,
        event,
        context,
    }),
});

exports.handler = server.createHandler({
    cors: {
        origin: '*',
        credentials: true,
        allowedHeaders: 'Content-Type, Authorization'
    },
});
