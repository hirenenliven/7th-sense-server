-- Migration: "log in as customer" sessions
--
-- The superadmin opens a customer's account in a new browser tab (POST /api/v1/auth/login-as).
-- That starts an ordinary session of the customer; impersonated_by records which superadmin
-- started it, for auditing (NULL = the user logged in with their own password). Run it once:
--   psql -U postgres -d 7th_sense -f migrations/031_add_auth_tokens_impersonated_by.sql

ALTER TABLE auth_tokens ADD COLUMN IF NOT EXISTS impersonated_by integer REFERENCES users(id);
