package main

import (
	"net/http"
	"os"
	"log"
	"encoding/json"
	"io/ioutil"

	g "github.com/graphql-go/graphql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	POSTGRES_CONNECTION_STRING = func () string {
		connection_string := "postgres://postgres:password@localhost:6432/postgres"
		if cs := os.Getenv("POSTGRES_CONNECTION_STRING"); cs != "" {
			connection_string = cs
		}
		return connection_string
	}()
	DB, DBErr = gorm.Open("postgres", POSTGRES_CONNECTION_STRING + "?sslmode=disable")
    )


type Author struct {
	Id int `json:"id,omitempty" gorm:"primary_key"`
	Name string `json:"name,omitempty"`
	Articles []Article `json:"articles,omitempty" gorm:"foreignkey:AuthorId"`
}

type Article struct {
	Id int `json:"id,omitempty" gorm:"primary_key"`
	Title string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	AuthorId int `json:"author_id,omitempty"`
	Author Author `json:"authors,omitempty" gorm:"association_foreignkey:Id"`
}

var authorType = g.NewObject(g.ObjectConfig{
	Name: "Author",
	Fields: g.Fields{
		"id": &g.Field{
			Type: g.String,
		},
		"name": &g.Field{
			Type: g.String,
		},
		"articles": &g.Field{
			Type: g.NewList(articleType),
		},
	},
})

var articleType = g.NewObject(g.ObjectConfig{
	Name: "Article",
	Fields: g.Fields{
		"id": &g.Field{
			Type: g.String,
		},
		"title": &g.Field{
			Type: g.String,
		},
		"content": &g.Field{
			Type: g.String,
		},
		"author_id": &g.Field{
			Type: g.Int,
		},
	},
})

var rootQuery = g.NewObject(g.ObjectConfig{
	Name: "Query",
	Fields: g.Fields{
		"authors": &g.Field{
			Type: g.NewList(authorType),
			Resolve: func (p g.ResolveParams) (interface{}, error) {
				authors := make([]Author, 0)
				DB.Preload("Articles").Find(&authors)
				return authors, nil
			},
		},
		"articles": &g.Field{
			Type: g.NewList(articleType),
			Resolve: func(p g.ResolveParams) (interface{}, error) {
				articles := make([]Article, 0)
				DB.Find(&articles)
				return articles, nil
			},
		},
	},

})

var rootMutation = g.NewObject(g.ObjectConfig{
	Name: "Mutation",
	Fields: g.Fields{
		"addAuthor": &g.Field{
			Type: authorType,
			Args: g.FieldConfigArgument{
				"name": &g.ArgumentConfig{
					Type: g.NewNonNull(g.String),
				},
			},
			Resolve: func(p g.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				author := Author{
					Name: name,
				}
				result := DB.Create(&author)
				if result.Error != nil{
					return nil, result.Error
				}
				return result.Value, nil
			},
		},
		"addArticle": &g.Field{
			Type: articleType,
			Args: g.FieldConfigArgument{
				"title": &g.ArgumentConfig{
					Type: g.NewNonNull(g.String),
				},
				"content": &g.ArgumentConfig{
					Type: g.NewNonNull(g.String),
				},
				"author_id": &g.ArgumentConfig{
					Type: g.NewNonNull(g.Int),
				},
			},
			Resolve: func(p g.ResolveParams) (interface{}, error) {
				title, _ := p.Args["title"].(string)
				content, _ := p.Args["content"].(string)
				author_id, _ := p.Args["author_id"].(int)
				article := Article{
					Title: title,
					Content: content,
					AuthorId: author_id,
				}
				result := DB.Create(&article)
				if result.Error != nil{
					return nil, result.Error
				}
				return result.Value, nil
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
		return
	}
	defer DB.Close()
	// run a graphql server in  local development mode
	if isLocal := os.Getenv("LAMBDA_EXECUTION_ENVIRONMENT"); isLocal == "local" {
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
