-- +goose Up
CREATE TABLE referrals (
                           id SERIAL PRIMARY KEY,
                           user_id INT UNIQUE NOT NULL,
                           code VARCHAR(255) UNIQUE NOT NULL,
                           expiry TIMESTAMP NOT NULL,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE referrals;