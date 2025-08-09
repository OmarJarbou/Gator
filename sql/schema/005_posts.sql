-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    description TEXT,
    published_at TIMESTAMP,
    feed_id UUID NOT NULL, 
    CONSTRAINT fk_posts_feeds FOREIGN KEY (feed_id) REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE posts;