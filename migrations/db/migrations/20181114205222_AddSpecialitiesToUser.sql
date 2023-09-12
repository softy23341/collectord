-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied


DROP INDEX user_trgm_name;
DROP FUNCTION user_search_name("user");

-- +goose StatementBegin
CREATE FUNCTION user_search_name("user")
    RETURNS text
    LANGUAGE SQL
    COST 5
AS $$
    SELECT lower('  ' || coalesce($1.first_name, '') || '  ' || coalesce($1.last_name, '') || '  ' || coalesce($1.email, '') || '   ' || array_to_string($1.tags, '   ', '') || '  ' || coalesce($1.speciality, ''))
$$
    IMMUTABLE
;
-- +goose StatementEnd

CREATE INDEX user_trgm_name ON "user" USING gist (user_search_name("user") gist_trgm_ops);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back


DROP INDEX user_trgm_name;
DROP FUNCTION user_search_name("user");
ALTER TABLE "user" DROP COLUMN IF EXISTS speciality;

-- +goose StatementBegin
CREATE FUNCTION user_search_name("user")
    RETURNS text
    STABLE
    LANGUAGE SQL
    IMMUTABLE
    COST 5
AS $$
    SELECT lower('  ' || coalesce($1.first_name, '') || '  ' || coalesce($1.last_name, '') || '  ' || coalesce($1.email, '') || '   ' || array_to_string($1.tags, '   ', ''))
$$;
-- +goose StatementEnd

CREATE INDEX user_trgm_name ON "user" USING gist (user_search_name("user") gist_trgm_ops);
