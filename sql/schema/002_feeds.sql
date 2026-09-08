-- +goose up
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL,
    name VARCHAR(50) NOT NULL,
    url VARCHAR(100) UNIQUE NOT NULL,
    user_id UUID NOT NULL,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose down
DROP TABLE feeds;