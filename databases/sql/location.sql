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
