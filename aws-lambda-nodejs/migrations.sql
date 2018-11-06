create table authors (id serial primary key, name text);

create table articles (id serial primary key, title text, content text, author_id integer not null references authors(id) );
