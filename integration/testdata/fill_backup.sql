SELECT current_database();

CREATE TABLE IF NOT EXISTS users (
    id       SERIAL PRIMARY KEY,
    name     VARCHAR(100) NOT NULL
);

INSERT INTO users (name)
VALUES
    ('Alice'),
    ('Bob'),
    ('Charlie'),
    ('Diana'),
    ('Ethan'),
    ('Fiona'),
    ('George'),
    ('Hannah'),
    ('Ian'),
    ('Julia');

SELECT * FROM users;

CREATE TABLE IF NOT EXISTS numbers (
    id       SERIAL PRIMARY KEY,
    name     VARCHAR(100) NOT NULL
);

INSERT INTO numbers (name)
VALUES
    ('One'),
    ('Two'),
    ('Three'),
    ('Four'),
    ('Five'),
    ('Six'),
    ('Seven'),
    ('Eight'),
    ('Nine'),
    ('Ten');

SELECT * FROM numbers;