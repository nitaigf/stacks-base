-- 002_user_management.up.sql backfills compatibility for an older schema shape.
-- The canonical 001_init baseline already owns these columns, tables and indexes,
-- so dropping them here would roll the database back past the intended base schema.
select 1;
