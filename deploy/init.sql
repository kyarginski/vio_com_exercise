--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

CREATE TABLE IF NOT EXISTS location (
    ip_address VARCHAR(15) not null,
    country_code VARCHAR(2),
    country VARCHAR(250),
    city VARCHAR(250),
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    mystery_value BIGINT
    );

ALTER TABLE location
    ADD CONSTRAINT location_ip_address_key PRIMARY KEY (ip_address);

COMMENT ON TABLE location IS 'A table for store locations. Author: Victor Kyarginskiy ';
COMMENT ON COLUMN location.ip_address IS 'IP Address of the location';
COMMENT ON COLUMN location.country_code IS 'Country Code of the location';
COMMENT ON COLUMN location.country IS 'Country of the location';
COMMENT ON COLUMN location.city IS 'City of the location';
COMMENT ON COLUMN location.latitude IS 'Latitude of the location';
COMMENT ON COLUMN location.longitude IS 'Longitude of the location';
COMMENT ON COLUMN location.mystery_value IS 'A mysterious value associated with the location';

INSERT INTO location (ip_address, country_code, country, city, latitude, longitude, mystery_value) VALUES ('70.95.73.73', 'TL', 'Saudi Arabia', 'Gradymouth', -49.16675918861615, -86.05920084416894, 2559997162)
    ON CONFLICT (ip_address) DO NOTHING
;
