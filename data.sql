CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    city TEXT NOT NULL
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    amount INT NOT NULL
);

INSERT INTO users (name, city) VALUES
('Alice', 'Almaty'),
('Bob', 'Astana'),
('Charlie', 'Almaty'),
('Dana', 'Shymkent');

INSERT INTO orders (user_id, amount) VALUES
(1, 100),
(1, 200),
(1, 50),
(2, 500),
(3, 300),
(3, 150),
(3, 400);
