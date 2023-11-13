package processes

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	"vio/internal/models"
)

func TestGetLocation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testFunction := func(ipAddress string) (*models.Location, error) {
		return GetLocation(db, ipAddress)
	}

	expectedLocation := &models.Location{
		IPAddress:    "127.0.0.1",
		CountryCode:  "US",
		Country:      "United States",
		City:         "New York",
		Latitude:     40.7128,
		Longitude:    -74.0060,
		MysteryValue: 1234567,
	}

	rows := sqlmock.NewRows([]string{"IPAddress", "CountryCode", "Country", "City", "Latitude", "Longitude", "MysteryValue"}).
		AddRow(expectedLocation.IPAddress, expectedLocation.CountryCode, expectedLocation.Country, expectedLocation.City, expectedLocation.Latitude, expectedLocation.Longitude, expectedLocation.MysteryValue)

	mock.ExpectQuery("SELECT (.+) FROM location WHERE ip_address = (.+)").WithArgs("127.0.0.1").WillReturnRows(rows)

	gotResult, err := testFunction("127.0.0.1")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	diff := cmp.Diff(expectedLocation, gotResult)
	if diff != "" {
		t.Fatal("result mismatch\n", diff)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProcess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()
	testFunction := func(locations *[]models.Location) (time.Duration, error) {
		startTime := time.Now()
		_, err := process(ctx, db, locations)
		return time.Since(startTime), err
	}

	locations := []models.Location{
		{
			IPAddress:    "127.0.0.1",
			CountryCode:  "US",
			Country:      "United States",
			City:         "New York",
			Latitude:     40.7128,
			Longitude:    -74.0060,
			MysteryValue: 12345678,
		},
	}

	expectedSQL := regexp.QuoteMeta(SQLInsert)

	mock.ExpectPrepare(expectedSQL).ExpectExec().WithArgs(
		locations[0].IPAddress, locations[0].CountryCode, locations[0].Country, locations[0].City,
		locations[0].Latitude, locations[0].Longitude, locations[0].MysteryValue,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	resultTime, resultErr := testFunction(&locations)

	assert.True(t, resultTime > 0, "unexpected negative or zero finish time")
	assert.Nil(t, resultErr, "unexpected error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
