CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(255),
    user_status INT(2) DEFAULT 1
);

CREATE TABLE IF NOT EXISTS categories
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS tags
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
)

CREATE TABLE IF NOT EXISTS tag_pets
(
    id SERIAL PRIMARY KEY,
    tag_id FOREIGN KEY REFERENCES tags(id),
    pet_id FOREIGN KEY REFERENCES pets(id)
);

CREATE TABLE IF NOT EXISTS pets
(
    id SERIAL PRIMARY KEY,
    category_id FOREIGN KEY REFERENCES categories(id),
    name VARCHAR(255) NOT NULL,
    status VARCHAR(255) DEFAULT 'available'
);

CREATE TABLE IF NOT EXISTS orders
(
    id SERIAL PRIMARY KEY,
    pet_id FOREIGN KEY REFERENCES pets(id),
    quantity INT NOT NULL,
    ship_date DATE NOT NULL,
    status VARCHAR(255) DEFAULT 'placed',
    complete BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS pet_photos
(
    id SERIAL PRIMARY KEY,
    pet_id FOREIGN KEY REFERENCES pets(id),
    photo_url VARCHAR(255) NOT NULL
);

--INSERT INTO users(name, email) VALUES
--('test','test'),
--('test2','test2');