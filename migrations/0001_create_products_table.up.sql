CREATE TABLE products (
                          id          BIGSERIAL PRIMARY KEY,
                          name        TEXT       NOT NULL,
                          type        TEXT       NOT NULL,
                          price       NUMERIC(10,2) NOT NULL CHECK (price >= 0),
                          description TEXT,
                          picture_url TEXT
);
