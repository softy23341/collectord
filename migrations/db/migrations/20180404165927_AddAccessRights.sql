-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- +goose StatementBegin

CREATE TABLE user_entity_right (
       id bigserial PRIMARY KEY,
       user_id bigint REFERENCES "user" ON DELETE CASCADE NOT NULL,
       entity_type character varying NOT NULL,
       entity_id bigint NOT NULL,

       "level" character varying NOT NULL,

       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       root_id bigint REFERENCES root ON DELETE CASCADE NOT NULL,
       UNIQUE(user_id, entity_type, entity_id)
);

-- +goose StatementEnd

-- +goose StatementBegin

INSERT INTO user_entity_right (user_id, entity_type, entity_id, "level", root_id)
SELECT u.id AS user_id, 'collection', c.id AS collection_id, CASE WHEN urr.typo = 30 THEN 'admin' ELSE 'write' END, urr.root_id AS root_id
FROM "user" AS u
INNER JOIN user_root_ref AS urr
  ON urr.user_id = u.id
INNER JOIN collection AS c
  ON c.root_id = urr.root_id
;

-- +goose StatementEnd

CREATE INDEX collection_image_media_id_idx ON collection (image_media_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE user_entity_right;
