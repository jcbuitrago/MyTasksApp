-- Simple schema to test DB setup
CREATE TABLE test_table (
    id SERIAL PRIMARY KEY,
    message TEXT NOT NULL
);
INSERT INTO test_table (message) VALUES ('Hello from PostgreSQL!');
