import graphene
from graphene import relay
from graphene_sqlalchemy import SQLAlchemyConnectionField, SQLAlchemyObjectType, utils
from models import Article as ArticleModel
from models import Author as AuthorModel


class Article(SQLAlchemyObjectType):
    class Meta:
        model = ArticleModel


class Author(SQLAlchemyObjectType):
    class Meta:
        model = AuthorModel

class Query(graphene.ObjectType):
    articles = graphene.List(Article)
    def resolve_articles(self, info):
        query = Article.get_query(info)
        return query.all()

    authors = graphene.List(Author)
    def resolve_authors(self, info):
        query = Author.get_query(info)
        return query.all()

schema = graphene.Schema(query=Query)
