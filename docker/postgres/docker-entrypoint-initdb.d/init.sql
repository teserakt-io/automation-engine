
-- Create an overall role "E4" to own the E4 database.
CREATE ROLE c2ae WITH NOSUPERUSER NOCREATEROLE NOINHERIT NOLOGIN NOREPLICATION NOBYPASSRLS;

-- Create a database with en_US.UTF-8 locale; set owner to e4.
CREATE DATABASE c2ae WITH OWNER=c2ae LC_COLLATE="en_US.UTF-8" LC_CTYPE="en_US.UTF-8" ENCODING=UTF8 TEMPLATE=template0;
\connect c2ae;

-- Create a specific login role:
CREATE ROLE c2ae_test WITH NOSUPERUSER NOCREATEROLE NOINHERIT LOGIN NOREPLICATION NOBYPASSRLS;
ALTER ROLE c2ae_test WITH ENCRYPTED PASSWORD 'teserakte4';

-- Create a specific schema for that role to operate in.
CREATE SCHEMA IF NOT EXISTS c2ae_test AUTHORIZATION c2ae_test;
-- Configure specific role to login using specified schema by default:
ALTER ROLE c2ae_test SET search_path = c2ae_test;

-- Give overall role E4 access to the schema.
GRANT ALL ON SCHEMA c2ae_test TO c2ae;



-- this removes access
-- for all users to the public schema.
-- Each user will need a schema they can access configured
-- as their default via search_path once you have done this.
REVOKE ALL PRIVILEGES ON SCHEMA public FROM PUBLIC
