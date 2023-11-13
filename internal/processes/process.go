package processes

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"vio/internal/models"
)

var SQLInsert = `INSERT INTO location (ip_address, country_code, country, city, latitude, longitude, mystery_value)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (ip_address) DO UPDATE
SET
	country_code = EXCLUDED.country_code,
	country = EXCLUDED.country,
	city = EXCLUDED.city,
	latitude = EXCLUDED.latitude,
	longitude = EXCLUDED.longitude,
	mystery_value = EXCLUDED.mystery_value
`

var SQLSelect = `SELECT * FROM location WHERE ip_address = $1`

func process(ctx context.Context, db *sql.DB, locations *[]models.Location) (string, error) {
	startTime := time.Now()
	stmt, err := db.Prepare(SQLInsert)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var errs []error

	for _, loc := range *locations {
		_, err := stmt.ExecContext(ctx, loc.IPAddress, loc.CountryCode, loc.Country, loc.City, loc.Latitude, loc.Longitude, loc.MysteryValue)
		if err != nil {
			errs = append(errs, err)
		}
	}
	finishTime := time.Since(startTime).String()

	return finishTime, errors.Join(errs...)
}

func GetLocation(db *sql.DB, ipAddress string) (*models.Location, error) {
	rows, err := db.Query(SQLSelect, ipAddress)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loc models.Location

	for rows.Next() {
		err := rows.Scan(
			&loc.IPAddress,
			&loc.CountryCode,
			&loc.Country,
			&loc.City,
			&loc.Latitude,
			&loc.Longitude,
			&loc.MysteryValue,
		)
		if err != nil {
			return nil, err
		}
	}

	return &loc, nil
}
