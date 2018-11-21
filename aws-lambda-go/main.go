package main

import (
	"net/http"
	"fmt"
	"os"

	gh "github.com/graphql-go/handler"
	g "github.com/graphql-go/graphql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	POSTGRES_CONNECTION_STRING = func () string {
		connection_string := "postgres://postgres:password@localhost:6432/postgres?sslmode=disable"
		if cs := os.Getenv("POSTGRES_CONNECTION_STRING"); cs != "" {
			connection_string = cs
		}
		return connection_string
	}()
	DB, DBErr = gorm.Open("postgres", POSTGRES_CONNECTION_STRING)
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

func main() {
	if DBErr != nil {
		fmt.Printf("%v", DBErr)
		return
	}
	defer DB.Close()
	h := gh.New(&gh.Config{
		Schema: &schema,
		Pretty: true,
		GraphiQL: true ,
	})
	http.Handle("/graphql", h)
	fmt.Printf("graphql server running on port 8080\n")
	http.ListenAndServe(":8080", nil)
}
