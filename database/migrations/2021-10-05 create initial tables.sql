create table if not exists users (
    id serial primary key,
	email varchar(255) unique not null,
    password varchar (255) not null,
	created_at timestamp not null,
    last_login timestamp not null
);

create table if not exists tasks (
   id bigserial primary key,
   user_id int not null,
   title varchar(255) not null,
   created_at timestamp not null,
   completed_at timestamp null
);