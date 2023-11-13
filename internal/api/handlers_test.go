package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"

	"vio/internal/models"
	"vio/internal/testhelpers"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestGetGeoLocation(t *testing.T) {
	wantResult := models.Location{
		IPAddress:    "70.95.73.73",
		CountryCode:  "TL",
		Country:      "Saudi Arabia",
		City:         "Gradymouth",
		Latitude:     -49.16675918861615,
		Longitude:    -86.05920084416894,
		MysteryValue: 2559997162,
	}
	testDB, err := testhelpers.NewTestDatabase(t)
	if err != nil {
		fmt.Println("Error connecting to the database: ", err)
		t.Skip("This test is excluded from unit tests.")
	}

	// Creating a test HTTP server.
	server := httptest.NewServer(GetGeoLocation(testDB.DB()))
	defer server.Close()

	// Use the server.URL and append the IP address to it
	ipAddress := "70.95.73.73"
	url := fmt.Sprintf("%s/api/geolocation/%s", server.URL, ipAddress)

	// Create a request with the modified URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Error creating HTTP request: %v", err)
	}

	// Make the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error making HTTP request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	got, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var gotResult models.Location

	err = json.Unmarshal(got, &gotResult)
	assert.NoError(t, err)

	diff := cmp.Diff(wantResult, gotResult)
	if diff != "" {
		t.Fatal("result mismatch\n", diff)
	}
}
