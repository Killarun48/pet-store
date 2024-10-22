CREATE TABLE IF NOT EXISTS users
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    phone TEXT,
    user_status INTEGER DEFAULT 1
);

CREATE TABLE IF NOT EXISTS categories
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tags
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tag_pets
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tag_id INTEGER,
    pet_id INTEGER,
    FOREIGN KEY (tag_id) REFERENCES tags(id),
    FOREIGN KEY (pet_id) REFERENCES pets(id)
);

CREATE TABLE IF NOT EXISTS pets
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER,
    name TEXT NOT NULL,
    status TEXT DEFAULT 'available',
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE TABLE IF NOT EXISTS orders
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pet_id INTEGER,
    quantity INTEGER NOT NULL,
    ship_date DATE NOT NULL,
    status TEXT DEFAULT 'placed',
    complete INTEGER DEFAULT 0,
    FOREIGN KEY (pet_id) REFERENCES pets(id)
);

CREATE TABLE IF NOT EXISTS pet_photos
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pet_id INTEGER,
    photo_url TEXT NOT NULL,
    FOREIGN KEY (pet_id) REFERENCES pets(id)
);

INSERT INTO users (username, first_name, last_name, email, password, phone, user_status) VALUES
('admin','John','Wick','JohnWick@continental.fake','admin','+123456789',1);

INSERT INTO categories (id, name) VALUES
(1, 'dog'),
(2, 'cat');

INSERT INTO tags (id, name) VALUES
(1, 'giftFromWife'),
(2, 'kennel');

INSERT INTO pets (id, category_id, name, status) VALUES
(1, 1, 'Daisy', 'sold'),
(2, 1, 'Good Boy', 'available');

INSERT INTO orders (id, pet_id, quantity, ship_date, status, complete) VALUES
(1, 1, 2, '2020-01-01', 'delivered', 1),
(2, 2, 1, '2021-06-01', 'placed', 1),
(3, 2, 1, '2021-06-01', 'approved', 0),
(4, 2, 1, '2021-09-01', 'delivered', 0),
(5, 2, 1, '2022-06-01', 'placed', 0),
(6, 2, 1, '2021-06-01', 'approved', 0),
(7, 2, 1, '2022-09-01', 'placed', 0),
(8, 2, 1, '2021-06-01', 'delivered', 0),
(9, 2, 1, '2021-09-01', 'placed', 0),
(10, 2, 1, '2021-06-01', 'placed', 0),
(11, 2, 1, '2022-06-01', 'approved', 0),
(12, 2, 1, '2021-06-01', 'placed', 0);

INSERT INTO pet_photos (id, pet_id, photo_url) VALUES
(1, 1, 'https://static1.srcdn.com/wordpress/wp-content/uploads/2024/01/keanu-reeves-as-john-wick-looks-down-at-daisy-in-a-scene-from-john-wick.jpg'),
(2, 2, 'https://static1.colliderimages.com/wordpress/wp-content/uploads/2022/11/pitbull-John-Wick.jpg');

INSERT INTO tag_pets (id, tag_id, pet_id) VALUES
(1, 1, 1),
(2, 2, 2);