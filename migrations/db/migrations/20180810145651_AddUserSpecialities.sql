-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE "user" ADD COLUMN speciality text;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE "user" DROP COLUMN IF EXISTS speciality;
ALTER TABLE "user" DROP COLUMN IF EXISTS specialities;
