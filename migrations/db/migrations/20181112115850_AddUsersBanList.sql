
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE user_ban (
    id bigserial PRIMARY KEY,
    creator_user_id bigint REFERENCES "user" ON DELETE CASCADE NOT NULL,
    user_id bigint REFERENCES "user" ON DELETE CASCADE NOT NULL,
    creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
    UNIQUE(creator_user_id, user_id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS user_ban;
