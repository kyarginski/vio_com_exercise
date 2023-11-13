package processes

import (
	"context"
	"regexp"
	"testing"

	"vio/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestProcessParallelWithMaxProcs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	expectedSQL := regexp.QuoteMeta(SQLInsert)
	ctx := context.Background()

	mock.ExpectPrepare(expectedSQL).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

	testFunction := func(locations *[]models.Location, maxWorkers int) (string, error) {
		return processParallelWithMaxProcs(ctx, db, locations, maxWorkers)
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

	resultTime, resultErr := testFunction(&locations, 10)
	assert.Nil(t, resultErr, "unexpected error")
	assert.True(t, len(resultTime) > 0, "unexpected finish time")

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
