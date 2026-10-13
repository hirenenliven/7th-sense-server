-- Migration: Make code uniqueness ignore soft-deleted rows
-- Version: 006
-- Description: The code columns of categories, brands, products and controllers were created
-- with a plain UNIQUE constraint, so a soft-deleted row kept its code reserved forever.
-- Each constraint is replaced by a partial unique index with the same name, so repository
-- error mapping (<table>_code_key) keeps working.

BEGIN;

ALTER TABLE categories DROP CONSTRAINT IF EXISTS categories_code_key;
CREATE UNIQUE INDEX IF NOT EXISTS categories_code_key ON categories (code) WHERE is_deleted = false;

ALTER TABLE brands DROP CONSTRAINT IF EXISTS brands_code_key;
CREATE UNIQUE INDEX IF NOT EXISTS brands_code_key ON brands (code) WHERE is_deleted = false;

ALTER TABLE products DROP CONSTRAINT IF EXISTS products_code_key;
CREATE UNIQUE INDEX IF NOT EXISTS products_code_key ON products (code) WHERE is_deleted = false;

ALTER TABLE controllers DROP CONSTRAINT IF EXISTS controllers_code_key;
CREATE UNIQUE INDEX IF NOT EXISTS controllers_code_key ON controllers (code) WHERE is_deleted = false;

COMMIT;
