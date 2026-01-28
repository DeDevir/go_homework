-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- +goose StatementBegin
DO
$$
    BEGIN
        IF NOT EXISTS (
            SELECT 1 FROM pg_type WHERE typname = 'order_status'
        ) THEN
            CREATE TYPE order_status AS ENUM (
                'PENDING_PAYMENT',
                'PAID',
                'CANCELLED'
                );
        END IF;

        IF NOT EXISTS (
            SELECT 1 FROM pg_type WHERE typname = 'payment_method'
        ) THEN
            CREATE TYPE payment_method AS ENUM (
                'UNKNOWN',
                'CARD',
                'SBP',
                'INVESTOR_MONEY',
                'CREDIT_CARD'
                );
        END IF;
    END
$$;
-- +goose StatementEnd


CREATE TABLE IF NOT EXISTS orders
(
    uuid             UUID PRIMARY KEY        DEFAULT gen_random_uuid(),
    user_uuid        UUID           NOT NULL,
    part_uuids       UUID[]         NOT NULL DEFAULT '{}',
    total_price      NUMERIC(14, 2) NOT NULL CHECK ( total_price > 0 ),
    transaction_uuid UUID           NULL,
    order_status     order_status   NOT NULL DEFAULT 'PENDING_PAYMENT',
    payment_method   payment_method NULL,

    created_at       TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ    NULL
);

CREATE INDEX IF NOT EXISTS idx_orders_uuid on orders (uuid);
CREATE INDEX IF NOT EXISTS idx_orders_user_uuid on orders (user_uuid);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_orders_updated_at()
    RETURNS trigger as
$$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
end;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_set_updated_at_orders
    BEFORE UPDATE
    ON orders
    FOR EACH ROW
EXECUTE FUNCTION set_orders_updated_at();

-- +goose Down

DROP TRIGGER IF EXISTS trg_set_updated_at_orders ON orders;

-- +goose StatementBegin
DO
$$
    BEGIN
        IF EXISTS(SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
            DROP TYPE order_status;
        END IF;

        IF EXISTS(SELECT 1 FROM pg_type WHERE typname = 'payment_method') THEN
            DROP TYPE payment_method;
        END IF;
    end
$$;

-- +goose StatementEnd

-- +goose StatementBegin
DROP FUNCTION IF EXISTS set_orders_updated_at();
-- +goose StatementEnd
DROP TABLE IF EXISTS orders;