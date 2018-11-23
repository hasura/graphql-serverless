const { ApolloServer, gql } = require('apollo-server');
const ApolloServerLambda = require('apollo-server-lambda').ApolloServer;
const Sequelize = require("sequelize");

const POSTGRES_CONNECTION_STRING = process.env.POSTGRES_CONNECTION_STRING || "postgres://postgres:password@localhost:6432/postgres";

const sequelize = new Sequelize(
    POSTGRES_CONNECTION_STRING, {}
);

const User = sequelize.define('user', {
    id: { type: Sequelize.INTEGER, autoIncrement: true, primaryKey: true },
    name: Sequelize.TEXT,
    balance: Sequelize.INTEGER
},
{
    timestamps: false
});

const typeDefs = gql`
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
`;

// Resolvers define the technique for fetching the types in the
// schema.  We'll retrieve books from the "books" array above.
const resolvers = {
    Query: {
        users: () => User.findAll()
    },
    Mutation: {
        addUser: async (_, { name, balance }) => {
            try {
                const user = await User.create({
                    name: name,
                    balance: balance
                });

                return user;
            } catch (e) {
                console.log(e);
                throw new Error(e);
            }
        },
        transfer: async (_, { userIdFrom, userIdTo, amount }) => {
            return await sequelize.transaction(async (t) => {
                if (userIdFrom == userIdTo) {
                    throw new Error("can't transfer on same account");
                }

                var userFrom = await User.findOne({
                    where: {id: userIdFrom}
                }, {transaction: t});

                if(userFrom === null) {
                    throw new Error("can't find user: " + userIdFrom);
                }

                var userTo = await User.findOne({
                    where: {id: userIdTo}
                }, {transaction: t});

                if(userTo === null) {
                    throw new Error("can't find user: " + userIdTo);
                }

                if((userFrom.balance - amount) < 0){
                    throw new Error("balance is too low");
                }

                userFrom.balance = userFrom.balance - amount;
                const resultFrom = await userFrom.save({transaction: t});

                userTo.balance = userTo.balance + amount;
                const resultTo = await userFrom.save({transaction: t});

                return resultFrom;

            });
        }
    }
};

// In the most basic sense, the ApolloServer can be started
// by passing type definitions (typeDefs) and the resolvers
// responsible for fetching the data for those types.
const server = new ApolloServerLambda({
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

// For local development
if( process.env.LAMBDA_EXECUTION_ENVIRONMENT == "local") {
    const serverLocal = new ApolloServer({ typeDefs, resolvers });

    serverLocal.listen().then(({ url }) => {
        console.log(`Server ready at ${url}`);
    });
}
