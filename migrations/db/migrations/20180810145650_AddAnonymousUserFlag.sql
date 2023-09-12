
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE "user" ADD COLUMN "is_anonymous" boolean NOT NULL DEFAULT FALSE;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE "user" DROP COLUMN IF EXISTS "is_anonymous";
