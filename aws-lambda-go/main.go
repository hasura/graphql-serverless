package main

import (
	"net/http"
	"fmt"

	gh "github.com/graphql-go/handler"
	g "github.com/graphql-go/graphql"
)

type Author struct {
	Id int `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Articles []Article `json:"articles,omitempty"`
}

type Article struct {
	Id int `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	AuthorId int `json:"author_id,omitempty"`
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
				return nil, nil

			},
		},
		"articles": &g.Field{
			Type: g.NewList(articleType),
			Resolve: func(p g.ResolveParams) (interface{}, error) {
				return nil, nil
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
				return nil, nil
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
				return nil, nil
			},
		},

	},
})

var schema, _ = g.NewSchema(g.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func main() {
	h := gh.New(&gh.Config{
		Schema: &schema,
		Pretty: true,
		GraphiQL: true ,
	})
	http.Handle("/graphql", h)
	fmt.Printf("graphql server running on port 8080\n")
	http.ListenAndServe(":8080", nil)
}
