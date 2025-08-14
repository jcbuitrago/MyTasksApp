-- User's table schema
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(25) NOT NULL,
    password VARCHAR(25) NOT NULL,
    picture VARCHAR(25)
);

-- Category's table schema

CREATE TABLE category (
    id SERIAL PRIMARY KEY,
    category_name VARCHAR(25) NOT NULL,
    description TEXT 
);

-- Task's table schema

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    task_description TEXT NOT NULL,
    creation_date DATE,
    tentative_due_date DATE,
    current_status ENUM('Backlog', 'In_Progress', 'Done') DEFAULT 'Backlog',
    parent_id INT,
    category_id INT,
    CONSTRAIN FOREIGN KEY (parent_id) REFERENCES category(id), ON DELETE CASCADE
);

INSERT INTO test_table (message) VALUES 
('Hello from PostgreSQL!'),
('This is a second test message.');
