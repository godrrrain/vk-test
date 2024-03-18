-- file: 10-create-user.sql
CREATE ROLE program WITH PASSWORD 'test';
ALTER ROLE program WITH LOGIN;

CREATE DATABASE movies;
GRANT ALL PRIVILEGES ON DATABASE movies TO program;

\c movies;

CREATE TABLE movie
(
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(150)  NOT NULL,
    description  VARCHAR(1000) NOT NULL,
    release_date VARCHAR(30)   NOT NULL,
    rating       INT           NOT NULL
        CHECK (rating BETWEEN 0 AND 10)
);

CREATE TABLE actor
(
    id        SERIAL PRIMARY KEY,
    name      VARCHAR(80) NOT NULL,
    gender    VARCHAR(30),
    birthday  VARCHAR(30)
);

CREATE TABLE movie_actor
(
    movie_id      INT REFERENCES movie (id),
    actor_id      INT REFERENCES actor (id)
);

GRANT ALL ON ALL TABLES IN SCHEMA public TO program;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO program;

-- INSERT INTO movie VALUES (1, 'Бойцовский клуб', 'Страховой работник разрушает рутину своей благополучной жизни. Культовая драма по книге Чака Паланика', '1999-09-10', 10);
-- INSERT INTO actor VALUES (1, 'Брэд Питт', 'Мужчина', '1963-12-18');
-- INSERT INTO movie_actor VALUES (1, 1);


-- INSERT INTO movie VALUES (2, 'Троя', '1193 год до нашей эры. Парис украл прекрасную Елену, жену царя Спарты Менелая. За честь Менелая вступается его брат – царь Агамемнон', '2004-05-14', 9);
-- INSERT INTO movie_actor VALUES (2, 1);