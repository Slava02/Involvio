BEGIN;

CREATE TABLE space (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    description VARCHAR(255)
);

CREATE TABLE category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE category_space (
    category_id INTEGER REFERENCES category(id) ON DELETE CASCADE,
    space_id INTEGER REFERENCES space(id) ON DELETE CASCADE,
    PRIMARY KEY (category_id, space_id)
);

CREATE TABLE "user" (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255)
);

CREATE TABLE user_space (
    user_id INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
    space_id INTEGER REFERENCES space(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, space_id)
);

CREATE TABLE user_event (
    user_id INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
    event_id INTEGER REFERENCES event(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, event_id)
);

CREATE TABLE event (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255),
    description VARCHAR(255),
    begin_date TIMESTAMP,
    end_date TIMESTAMP
);

CREATE TABLE category_event (
    event_id INTEGER REFERENCES event(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES category(id) ON DELETE CASCADE,
    parameter BOOLEAN,
    PRIMARY KEY (event_id, category_id)
);

COMMIT;