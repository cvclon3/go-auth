-- migrations/000001_create_users_table.down.sql
-- Add down migration script here
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS user_profile;