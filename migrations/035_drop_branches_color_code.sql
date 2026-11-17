-- Migration: branches no longer have a colour
--
-- The branch colour (branches.color_code, a #RRGGBB value) was entered on the branch form and shown
-- in the Branches list, but nothing else used it. The field is removed from the API and the client,
-- and the column is dropped; its values are lost. Deploy the server without color_code first: an
-- older server reads the column and fails once it is gone. Safe to run more than once:
--   psql -U postgres -d 7th_sense -f migrations/035_drop_branches_color_code.sql
--
-- Rollback (the values do not come back):
--   ALTER TABLE branches ADD COLUMN IF NOT EXISTS color_code character varying(10);

ALTER TABLE branches DROP COLUMN IF EXISTS color_code;
