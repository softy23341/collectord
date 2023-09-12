-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- +goose StatementBegin

CREATE TABLE "task" (
       id bigserial PRIMARY KEY,
       title character varying NOT NULL,
       description text NOT NULL DEFAULT '',
       creator_user_id bigint REFERENCES "user" ON DELETE SET NULL,
       assigned_user_id bigint REFERENCES "user" ON DELETE SET NULL,
       deadline timestamp WITH time zone NULL,

       status character varying NOT NULL,
       archive boolean NOT NULL DEFAULT false,

       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       user_uniq_id bigint NOT NULL,
       UNIQUE(creator_user_id, user_uniq_id)
);

-- +goose StatementEnd

CREATE INDEX task_assigned_user_id_idx ON task (assigned_user_id);


-- +goose StatementBegin

CREATE TABLE task_media_ref (
       task_id bigint REFERENCES "task" ON DELETE CASCADE NOT NULL,
       media_id bigint REFERENCES media ON DELETE CASCADE NOT NULL,
       "position" integer DEFAULT 0 NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       PRIMARY KEY (task_id, media_id)
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE task_group_ref (
       task_id bigint REFERENCES "task" ON DELETE CASCADE NOT NULL,
       group_id bigint REFERENCES "group" ON DELETE CASCADE NOT NULL,
       "position" integer DEFAULT 0 NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       PRIMARY KEY (task_id, group_id)
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE task_collection_ref (
       task_id bigint REFERENCES "task" ON DELETE CASCADE NOT NULL,
       collection_id bigint REFERENCES "collection" ON DELETE CASCADE NOT NULL,
       "position" integer DEFAULT 0 NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       PRIMARY KEY (task_id, collection_id)
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE task_object_ref (
       task_id bigint REFERENCES "task" ON DELETE CASCADE NOT NULL,
       object_id bigint REFERENCES "object" ON DELETE CASCADE NOT NULL,
       "position" integer DEFAULT 0 NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       PRIMARY KEY (task_id, object_id)
);

-- +goose StatementEnd

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE IF EXISTS task_object_ref;
DROP TABLE IF EXISTS task_collection_ref;
DROP TABLE IF EXISTS task_group_ref;
DROP TABLE IF EXISTS task_media_ref;
DROP TABLE IF EXISTS task;
