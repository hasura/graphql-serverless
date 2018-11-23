create table author (id serial primary key, name text);

create table article (id serial primary key, title text, content text, author_id integer not null references authors(id) );
