
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE "invite" ADD COLUMN "to_user_email" character varying;
ALTER TABLE "invite" ALTER COLUMN "to_user_id" DROP NOT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE "invite" DROP COLUMN IF EXISTS "to_user_email";
ALTER TABLE "invite" ALTER COLUMN "to_user_id" SET NOT NULL;
