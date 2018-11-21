#!/usr/bin/env python

import os
from flask import Flask
from flask_graphql import GraphQLView
import graphene
from graphene_sqlalchemy import SQLAlchemyConnectionField, SQLAlchemyObjectType, utils
from sqlalchemy import Column, DateTime, ForeignKey, Integer, String, func, create_engine
from sqlalchemy.orm import backref, relationship, scoped_session, sessionmaker
from sqlalchemy.ext.declarative import declarative_base

app = Flask(__name__)

################################## DATABASE #################################

POSTGRES_CONNECTION_STRING = os.environ.get('POSTGRES_CONNECTION_STRING')

engine = create_engine(POSTGRES_CONNECTION_STRING, convert_unicode=True)
db_session = scoped_session(sessionmaker(autocommit=False,
                     autoflush=False,
                     bind=engine))
Base = declarative_base()
Base.query = db_session.query_property()


def init_db():
  # import all modules here that might define models so that
  # they will be registered properly on the metadata.  Otherwise
  # you will have to import them first before calling init_db()
  from models import Article, Author
  Base.metadata.drop_all(bind=engine)
  Base.metadata.create_all(bind=engine)
  db_session.commit()

#############################################################################

################################## MODELS ###################################

class ArticleModel(Base):
  __tablename__ = 'article'
  id = Column(Integer, primary_key=True)
  title = Column(String)
  content = Column(String)
  author_id = Column(Integer, ForeignKey('author.id'))

class AuthorModel(Base):
  __tablename__ = 'author'
  id = Column(Integer, primary_key=True)
  name = Column(String)
  age = Column(Integer)
  articles = relationship(
    ArticleModel,
    backref=backref(
      'author',
      cascade='delete,all'
    )
  )
#############################################################################


################################## SCHEMA ###################################
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
#############################################################################


################################ FLASK APP ##################################
app.add_url_rule('/', view_func=GraphQLView.as_view('graphql', schema=schema, graphiql=True))

@app.teardown_appcontext
def shutdown_session(exception=None):
  db_session.remove()

#############################################################################

################################## EXECUTION ##################################
if __name__ == '__main__':
  if os.environ.get('DATABASE_INIT') == '1':
    init_db()
  app.run()
#############################################################################
