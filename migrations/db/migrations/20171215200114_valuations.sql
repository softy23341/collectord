
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- +goose StatementBegin

CREATE TABLE "object_valuations" (
       id bigserial PRIMARY KEY,
       object_id bigint REFERENCES object ON DELETE CASCADE NOT NULL,
       name character varying NOT NULL DEFAULT '',
       comment text NOT NULL DEFAULT '',
       date date,
       price int8range NOT NULL,
       price_currency_id bigint REFERENCES currency ON DELETE RESTRICT NOT NULL,
       CONSTRAINT price_from CHECK (
         (price IS NOT NULL AND price_currency_id IS NOT NULL
          AND lower(price) > 0 AND upper(price) > 0
         )
       )
);

-- +goose StatementEnd

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
