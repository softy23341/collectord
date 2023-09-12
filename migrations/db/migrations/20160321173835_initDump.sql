-- +goose Up

SET client_min_messages = warning;
SET client_encoding = 'UTF8';
-- CREATE EXTENSION pg_trgm;
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- ALTER extension pg_trgm SET SCHEMA pg_catalog;

-- +goose StatementBegin

CREATE TABLE root (
       id bigserial PRIMARY KEY,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE media (
       id bigserial PRIMARY KEY,
       root_id bigint NULL REFERENCES root,
       user_id bigint NULL,
       user_uniq_id bigint NOT NULL,
       type smallint NOT NULL,
       extra jsonb NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       UNIQUE(user_id, user_uniq_id)
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE "user" (
       id bigserial PRIMARY KEY,
       first_name character varying NOT NULL DEFAULT '',
       last_name character varying NOT NULL DEFAULT '',
       email character varying NOT NULL,
       encrypted_password character varying NOT NULL,
       avatar_media_id bigint REFERENCES media ON DELETE SET NULL,
       description text NOT NULL DEFAULT '',
       system_user boolean NOT NULL DEFAULT FALSE,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       last_event_seq_no bigint DEFAULT 0 NOT NULL,
       n_unread_messages integer DEFAULT 0 NOT NULL,
       n_unread_notifications integer DEFAULT 0 NOT NULL,
       UNIQUE(email)
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE invite (
       id bigserial PRIMARY KEY,
       creator_user_id bigint REFERENCES "user" ON DELETE CASCADE NOT NULL,
       root_id bigint REFERENCES root ON DELETE CASCADE NOT NULL,
       to_user_id bigint REFERENCES "user" ON DELETE CASCADE NOT NULL,
       token character varying NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       status smallint NOT NULL,
       UNIQUE(token)
);

-- +goose StatementEnd

CREATE UNIQUE INDEX one_active_invite_per_user_root ON invite(root_id, creator_user_id, to_user_id) WHERE status = 0;

-- +goose StatementBegin

CREATE TABLE message (
       id bigserial PRIMARY KEY,
       user_id bigint REFERENCES "user" ON DELETE CASCADE NOT NULL,
       user_uniq_id bigint NOT NULL,
       peer_id bigint NOT NULL,
       peer_type smallint NOT NULL,
       type smallint NOT NULL,
       extra jsonb NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       UNIQUE(user_id, user_uniq_id)
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE dialog (
       id bigserial PRIMARY KEY,
       least_user_id bigint REFERENCES "user" ON DELETE CASCADE NOT NULL,
       greatest_user_id bigint REFERENCES "user" ON DELETE CASCADE NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       UNIQUE(least_user_id, greatest_user_id)
);

-- +goose StatementEnd


-- +goose StatementBegin

CREATE TABLE conversation_user (
       id bigserial PRIMARY KEY,
       user_id bigint REFERENCES "user" ON DELETE CASCADE NOT NULL,
       peer_id bigint NOT NULL,
       peer_type smallint NOT NULL,
       joined_at timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       last_read_message_id bigint REFERENCES message ON DELETE RESTRICT NULL,
       n_unread_messages integer DEFAULT 0 NOT NULL,
       inviter_user_id bigint REFERENCES "user" ON DELETE RESTRICT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       UNIQUE(user_id, peer_type, peer_id)
);

-- +goose StatementEnd


-- +goose StatementBegin

CREATE TABLE unread_message (
       id bigserial PRIMARY KEY,
       message_id bigint REFERENCES message ON DELETE RESTRICT NOT NULL,
       user_id bigint REFERENCES "user" ON DELETE RESTRICT NOT NULL,
       peer_id bigint NOT NULL,
       peer_type smallint NOT NULL
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE chat (
    id bigserial PRIMARY KEY,
    creator_user_id bigint REFERENCES "user" ON DELETE CASCADE NOT NULL,
    user_uniq_id bigint NOT NULL,
    admin_user_id bigint REFERENCES "user" ON DELETE CASCADE NOT NULL,
    name text NOT NULL,
    avatar_media_id bigint REFERENCES media ON DELETE SET NULL,
    last_read_message_id bigint DEFAULT 0 NOT NULL,
    creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
    UNIQUE(creator_user_id, user_uniq_id)
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE event (
       id bigserial PRIMARY KEY,
       user_id bigint REFERENCES "user" ON DELETE CASCADE NOT NULL,
       seq_no bigint NOT NULL,
       type smallint NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       status smallint NOT NULL,
       extra jsonb NOT NULL,
       UNIQUE(user_id, seq_no)
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE FUNCTION user_search_name("user")
    RETURNS text
    STABLE
    LANGUAGE SQL
    COST 5
AS $$
    SELECT lower('  ' || coalesce($1.first_name, '') || '  ' || coalesce($1.last_name, '') || '  ' || coalesce($1.email, ''))
$$;

-- +goose StatementEnd

-- +goose StatementBegin

CREATE INDEX user_trgm_name ON "user" USING gist (user_search_name("user") gist_trgm_ops);

-- +goose StatementEnd

--+goose StatementBegin

CREATE TABLE user_session (
       id bigserial PRIMARY KEY,
       user_id bigint REFERENCES "user" ON DELETE CASCADE NOT NULL,
       auth_token character varying NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       UNIQUE(auth_token)
);

-- +goose StatementEnd

--+goose StatementBegin

CREATE TABLE user_session_device (
       id bigserial PRIMARY KEY,
       session_id bigint REFERENCES user_session ON DELETE CASCADE NOT NULL,
       typo smallint NOT NULL,
       token character varying NOT NULL,
       push_notification_sandbox boolean,
       UNIQUE(session_id, typo, token, push_notification_sandbox)
);

-- +goose StatementEnd

--+goose StatementBegin

CREATE TABLE user_root_ref (
       user_id bigint REFERENCES "user" ON DELETE CASCADE NOT NULL,
       root_id bigint REFERENCES root ON DELETE CASCADE NOT NULL,
       typo smallint NOT NULL DEFAULT 0,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       PRIMARY KEY (user_id, root_id)
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE material (
       id bigserial PRIMARY KEY,
       root_id bigint NULL REFERENCES root,
       name character varying NOT NULL,
       normal_name character varying NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       UNIQUE(root_id, normal_name)
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE actor (
       id bigserial PRIMARY KEY,
       root_id bigint NULL REFERENCES root,
       name character varying NOT NULL,
       normal_name character varying NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       UNIQUE(root_id, normal_name)
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE origin_location (
       id bigserial PRIMARY KEY,
       root_id bigint NULL REFERENCES root,
       name character varying NOT NULL,
       normal_name character varying NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       UNIQUE(root_id, normal_name)
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE "group" (
       id bigserial PRIMARY KEY,
       user_id bigint REFERENCES "user" ON DELETE SET NULL,
       user_uniq_id bigint NOT NULL,
       root_id bigint NOT NULL REFERENCES root,
       "name" character varying NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       UNIQUE(user_id, user_uniq_id),
       UNIQUE(root_id, "name")
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE collection (
       id bigserial PRIMARY KEY,
       user_id bigint REFERENCES "user" ON DELETE SET NULL,
       user_uniq_id bigint,
       root_id bigint NOT NULL REFERENCES root,
       name character varying NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       typo smallint NOT NULL DEFAULT 0,
       description text NOT NULL DEFAULT '',
       image_media_id bigint REFERENCES media ON DELETE SET NULL,
       UNIQUE(user_id, user_uniq_id)
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE collection_group_ref (
       collection_id bigint REFERENCES collection ON DELETE CASCADE NOT NULL,
       group_id bigint REFERENCES "group" ON DELETE CASCADE NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       PRIMARY KEY (collection_id, group_id)
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE currency (
       id bigserial PRIMARY KEY,
       symbol character varying NOT NULL,
       code character varying NOT NULL,
       num character varying NOT NULL,
       e smallint NOT NULL,
       currency character varying NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE object_status (
       id bigserial PRIMARY KEY,
       name character varying NOT NULL,
       description text NOT NULL DEFAULT '',
       image_media_id bigint REFERENCES media,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE named_date_interval (
       id bigserial PRIMARY KEY,
       root_id bigint NULL REFERENCES root,
       production_date_interval_from bigint,
       production_date_interval_to bigint,
       name character varying NOT NULL,
       normal_name character varying NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       UNIQUE(root_id, normal_name)
)

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE object (
       id bigserial PRIMARY KEY,
       system_id uuid NOT NULL DEFAULT uuid_generate_v1(),
       collection_id bigint NOT NULL REFERENCES collection,
       user_id bigint REFERENCES "user" ON DELETE SET NULL,
       user_uniq_id bigint NOT NULL,
       name character varying NOT NULL,
       production_date_interval_id bigint REFERENCES named_date_interval ON DELETE SET NULL,
       production_date_interval_from bigint,
       production_date_interval_to bigint,
       description text NOT NULL DEFAULT '',
       provenance text NOT NULL DEFAULT '',
       purchase_date date,
       purchase_price bigint,
       purchase_price_currency_id bigint REFERENCES currency ON DELETE RESTRICT,
       root_id_number character varying NOT NULL DEFAULT '',
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       update_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       CONSTRAINT interval_direction CHECK (
         (
           production_date_interval_id IS NOT NULL
           AND (production_date_interval_to IS NULL AND production_date_interval_from IS NULL)
         )
         OR
         (
           production_date_interval_id IS NULL
           AND production_date_interval_to IS NOT NULL AND production_date_interval_from IS NOT NULL
           AND (production_date_interval_to >= production_date_interval_from)
         )
         OR
         (
           production_date_interval_id IS NULL
           AND (production_date_interval_to IS NULL AND production_date_interval_from IS NULL)
         )
       ),
       CONSTRAINT purchase_price CHECK (
         (purchase_price IS NULL AND purchase_price_currency_id IS NULL)
         OR
         (purchase_price IS NOT NULL AND purchase_price_currency_id IS NOT NULL
          AND purchase_price > 0
         )
       ),
       UNIQUE(user_id, user_uniq_id),
       UNIQUE(system_id)
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE badge (
       id bigserial PRIMARY KEY,
       name character varying NOT NULL,
       normal_name character varying NOT NULL,
       color character varying NOT NULL,
       root_id bigint NULL REFERENCES root,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       UNIQUE(root_id, normal_name),
       UNIQUE(root_id, color)
);
-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE object_badge_ref (
       object_id bigint REFERENCES object ON DELETE CASCADE NOT NULL,
       badge_id bigint REFERENCES badge ON DELETE CASCADE NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       PRIMARY KEY (object_id, badge_id)
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE object_status_ref (
       id bigserial PRIMARY KEY,
       object_id bigint NOT NULL REFERENCES object ON DELETE CASCADE,
       object_status_id bigint NOT NULL REFERENCES object_status ON DELETE RESTRICT,
       start_date timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       notification_date timestamp WITH time zone,
       description text NOT NULL DEFAULT '',
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE object_media_ref (
       object_id bigint REFERENCES object ON DELETE CASCADE NOT NULL,
       media_id bigint REFERENCES media ON DELETE CASCADE NOT NULL,
       media_position integer DEFAULT 0 NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       PRIMARY KEY (object_id, media_id)
);

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TABLE object_material_ref (
       object_id bigint REFERENCES object ON DELETE CASCADE NOT NULL,
       material_id bigint REFERENCES material ON DELETE CASCADE NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       PRIMARY KEY(object_id, material_id)
);

-- +goose StatementEnd


-- +goose StatementBegin

CREATE TABLE object_actor_ref (
       id bigserial PRIMARY KEY,
       object_id bigint REFERENCES object ON DELETE CASCADE NOT NULL,
       actor_id bigint REFERENCES actor ON DELETE CASCADE NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       UNIQUE(object_id, actor_id)
);

-- +goose StatementEnd


-- +goose StatementBegin

CREATE TABLE object_origin_location_ref (
       object_id bigint REFERENCES object ON DELETE CASCADE NOT NULL,
       origin_location_id bigint REFERENCES origin_location ON DELETE CASCADE NOT NULL,
       creation_time timestamp WITH time zone NOT NULL DEFAULT current_timestamp,
       PRIMARY KEY (object_id, origin_location_id)
);

-- +goose StatementEnd

-- +goose StatementBegin

ALTER TABLE media ADD CONSTRAINT media_user_ref FOREIGN KEY (user_id) REFERENCES "user" ON DELETE SET NULL;

-- +goose StatementEnd

-- indexes
CREATE INDEX collection_root_id_idx ON collection (root_id);
CREATE INDEX object_collection_id_idx ON object (collection_id);
CREATE INDEX object_status_ref_object_id_idx ON object_status_ref (object_id);


-- SQL in section 'Up' is executed when this migration is applied

-------------------------------------------------------------------------------------

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
