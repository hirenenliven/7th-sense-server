-- Migration: clienteles, the clients a customer makes parts for
--
-- A clientele belongs to a customer (user_id). Client ID (client_code) is the customer's
-- own number for the client and is unique per customer (case-insensitive). Country and
-- state are ISO 3166 codes ("IN", "GJ"); the contact number is stored as +<code><number>.
--
-- parts.clientele_id (migration 016) now references clienteles. A clientele that parts
-- use cannot be deleted (service.ClienteleService).
--
-- Rollback:
--   ALTER TABLE parts DROP CONSTRAINT IF EXISTS parts_clientele_id_fkey;
--   DROP TABLE IF EXISTS clienteles;

BEGIN;

CREATE TABLE IF NOT EXISTS clienteles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    client_code VARCHAR(50) NOT NULL,
    name VARCHAR(250) NOT NULL,
    email VARCHAR(255),
    contact_number VARCHAR(20) NOT NULL,
    address VARCHAR(250) NOT NULL,
    country_code VARCHAR(2) NOT NULL,
    state_code VARCHAR(10) NOT NULL,
    city VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- A customer cannot use the same Client ID twice
CREATE UNIQUE INDEX IF NOT EXISTS clienteles_user_client_code_key
    ON clienteles (user_id, LOWER(client_code)) WHERE is_deleted = false;

CREATE INDEX IF NOT EXISTS idx_clienteles_user_id ON clienteles(user_id);
CREATE INDEX IF NOT EXISTS idx_clienteles_is_deleted ON clienteles(is_deleted);

CREATE TRIGGER update_clienteles_updated_at BEFORE UPDATE ON clienteles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

ALTER TABLE parts
    ADD CONSTRAINT parts_clientele_id_fkey FOREIGN KEY (clientele_id) REFERENCES clienteles(id);

CREATE INDEX IF NOT EXISTS idx_parts_clientele_id ON parts(clientele_id);

COMMIT;
