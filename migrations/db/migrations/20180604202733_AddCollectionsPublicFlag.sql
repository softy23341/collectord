
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE "collection" ADD COLUMN "public" boolean NOT NULL DEFAULT FALSE;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE "collection" DROP COLUMN IF EXISTS "public";

