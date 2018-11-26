package main

import (
	"net/http"
	"os"
	"log"
	"encoding/json"
	"io/ioutil"
	"errors"

	g "github.com/graphql-go/graphql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	POSTGRES_CONNECTION_STRING = func () string {
		connection_string := "postgres://postgres:password@localhost:6432/postgres"
		if cs := os.Getenv("POSTGRES_CONNECTION_STRING"); cs == "" {
			connection_string = cs
		}
		return connection_string
	}()
	DB, DBErr = gorm.Open("postgres", POSTGRES_CONNECTION_STRING + "?sslmode=disable")
    )


type User struct {
	Id int `json:"id,omitempty" gorm:"primary_key"`
	Name string `json:"name,omitempty"`
	Balance int `json:"balance,omitempty"`
}

var userType = g.NewObject(g.ObjectConfig{
	Name: "User",
	Fields: g.Fields{
		"id": &g.Field{
			Type: g.Int,
		},
		"name": &g.Field{
			Type: g.String,
		},
		"balance": &g.Field{
			Type: g.Int,
		},
	},
})

var rootQuery = g.NewObject(g.ObjectConfig{
	Name: "Query",
	Fields: g.Fields{
		"hello": &g.Field{
			Type: g.String,
			Resolve: func(p g.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
		"users": &g.Field{
			Type: g.NewList(userType),
			Resolve: func (p g.ResolveParams) (interface{}, error) {
				users:= make([]User, 0)
				result := DB.Find(&users)
				if result.Error != nil{
					return nil, result.Error
				}
				return users, nil
			},
		},
	},

})

var rootMutation = g.NewObject(g.ObjectConfig{
	Name: "Mutation",
	Fields: g.Fields{
		"addUser": &g.Field{
			Type: userType,
			Args: g.FieldConfigArgument{
				"name": &g.ArgumentConfig{
					Type: g.NewNonNull(g.String),
				},
				"balance": &g.ArgumentConfig{
					Type: g.NewNonNull(g.Int),
				},
			},
			Resolve: func(p g.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				balance, _ := p.Args["balance"].(int)
				user := User{
					Name: name,
					Balance: balance,
				}
				result := DB.Create(&user)
				if result.Error != nil{
					return nil, result.Error
				}
				return result.Value, nil
			},
		},
		"transfer": &g.Field{
			Type: userType,
			Args: g.FieldConfigArgument{
				"userIdFrom": &g.ArgumentConfig{
					Type: g.NewNonNull(g.Int),
				},
				"userIdTo": &g.ArgumentConfig{
					Type: g.NewNonNull(g.Int),
				},
				"amount": &g.ArgumentConfig{
					Type: g.NewNonNull(g.Int),
				},
			},
			Resolve: func(p g.ResolveParams) (interface{}, error) {
				userIdFrom, _ := p.Args["userIdFrom"].(int)
				userIdTo, _ := p.Args["userIdTo"].(int)
				amount, _ := p.Args["amount"].(int)

				if userIdFrom == userIdTo {
					return nil, errors.New("can't transfer on same account")
				}

				tx := DB.Begin()
				var userFrom User
				result := tx.Where("id= ?", userIdFrom).First(&userFrom)
				if result.Error != nil{
					return nil, result.Error
				}
				var userTo User
				result = tx.Where("id= ?", userIdTo).First(&userTo)
				if result.Error != nil{
					return nil, result.Error
				}

				if (userFrom.Balance - amount) < 0 {
					return nil, errors.New("balance is too low")

				}

				userFrom.Balance = userFrom.Balance - amount
				resultFrom := tx.Save(&userFrom)

				userTo.Balance = userTo.Balance + amount
				tx.Save(&userTo)

				tx.Commit()

				return resultFrom.Value, nil
			},
		},

	},
})

var schema, _ = g.NewSchema(g.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

type MyEvent struct {
        Name string `json:"name"`
}

type GraphQLRequest struct {
	Query string `json:"query"`
	OperationName string `json:"operationName"`
	Variables map[string]interface{} `json:"variables"`
}

func executeQuery(graphqlRequest GraphQLRequest, schema g.Schema) *g.Result {
	result := g.Do(g.Params{
		Schema:        schema,
		RequestString: graphqlRequest.Query,
		OperationName: graphqlRequest.OperationName,
		VariableValues: graphqlRequest.Variables,
	})
	if len(result.Errors) > 0 {
		log.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	//parse the graphql request from lambda context
	graphQLRequest := func (r events.APIGatewayProxyRequest) GraphQLRequest {
		query := r.QueryStringParameters["query"]
		operationName := r.QueryStringParameters["operationName"]
		variablesString := r.QueryStringParameters["variables"]
		variables := make(map[string]interface{})
		_ = json.Unmarshal([]byte(variablesString), &variables)
		graphqlRequest := GraphQLRequest{
			Query: query,
			OperationName: operationName,
			Variables: variables,
		}
		if query == "" {
			var graphqlRequest GraphQLRequest
			_ = json.Unmarshal([]byte(r.Body), &graphqlRequest)
			return graphqlRequest
		}
		return graphqlRequest
	}(request)

	log.Printf("%v", graphQLRequest)
	result := executeQuery(graphQLRequest, schema)
	resultJson, _ :=  json.Marshal(result)

	return events.APIGatewayProxyResponse{
		Body:       string(resultJson),
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE",
			"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		},
	}, nil
}

func main() {
	if DBErr != nil {
		log.Printf("%v", DBErr)
	}
	defer DB.Close()
	// run a graphql server in  local development mode
	if isLocal := os.Getenv("LAMBDA_LOCAL_DEVELOPMENT"); isLocal == "1" {
		http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
			graphqlRequest := func (r *http.Request) GraphQLRequest {
				query := r.URL.Query().Get("query")
				operationName := r.URL.Query().Get("operationName")
				variablesString := r.URL.Query().Get("variables")
				variables := make(map[string]interface{})
				_ = json.Unmarshal([]byte(variablesString), &variables)
				graphqlRequest := GraphQLRequest{
					Query: query,
					OperationName: operationName,
					Variables: variables,
				}
				if query == "" {
					b, _ := ioutil.ReadAll(r.Body)
					var graphqlRequest GraphQLRequest
					_ = json.Unmarshal(b, &graphqlRequest)
					return graphqlRequest
				}
				return graphqlRequest
			}(r)
			log.Printf("%v", graphqlRequest)
			result := executeQuery(graphqlRequest, schema)
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			json.NewEncoder(w).Encode(result)
		})
		log.Printf("graphql server running on port 8080\n")
		http.ListenAndServe(":8080", nil)
	} else {
		lambda.Start(handler)
	}
}
