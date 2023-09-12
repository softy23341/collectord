
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin

CREATE OR REPLACE FUNCTION sum_int8range(input int8range[]) RETURNS bigint AS $$
DECLARE
  result bigint := 0;
  irange int8range;
BEGIN
  FOREACH irange IN ARRAY input LOOP
    result = result + ((upper(irange) - 1 ) - (((upper(irange) - 1) - lower(irange))/2));
  END LOOP;
  RETURN result;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
