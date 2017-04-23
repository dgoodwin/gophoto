
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE media(
    id SERIAL PRIMARY KEY,
    created TIMESTAMP WITH TIME ZONE NOT NULL,
    uploaded TIMESTAMP WITH TIME ZONE NOT NULL,
    filename VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    res_x INTEGER NOT NULL,
    res_y INTEGER NOT NULL,
    size BIGINT NOT NULL
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE media;

