package processes

import (
	"path/filepath"
	"runtime"
	"testing"

	"vio/internal/models"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_loadData(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		wantResult      []models.Location
		wantStatistics  *models.LoadStatistics
		wantError       bool
		wantErrorString string
	}{
		{
			name:       "Load with errors",
			path:       "process_data_error",
			wantResult: []models.Location{},
			wantStatistics: &models.LoadStatistics{
				LoadTime:   "395µs",
				FilesCount: 1,
				Accepted:   0,
				Discarded:  0,
				Total:      0,
			},
			wantError:       true,
			wantErrorString: "record on line 2: wrong number of fields",
		},
		{
			name: "Load with success",
			path: "process_data_good",
			wantResult: []models.Location{
				{
					IPAddress:    "70.95.73.73",
					CountryCode:  "TL",
					Country:      "Saudi Arabia",
					City:         "Gradymouth",
					Latitude:     -49.16675918861615,
					Longitude:    -86.05920084416894,
					MysteryValue: 2559997162,
				},
				{
					IPAddress:    "125.159.20.54",
					CountryCode:  "LI",
					Country:      "Guyana",
					City:         "Port Karson",
					Latitude:     -78.2274228596799,
					Longitude:    -163.26218895343357,
					MysteryValue: 1337885276,
				},
			},
			wantStatistics: &models.LoadStatistics{
				LoadTime:   "395µs",
				FilesCount: 1,
				Accepted:   2,
				Discarded:  2,
				Total:      4,
			},
			wantError:       false,
			wantErrorString: "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := readFixture(t, tt.path)
			got, got1, err := loadData(path)
			if !tt.wantError {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErrorString)
			}

			diff := cmp.Diff(tt.wantResult, got)
			if diff != "" {
				t.Fatal("Result mismatch\n", diff)
			}
			diff = cmp.Diff(tt.wantStatistics, got1, cmpopts.IgnoreFields(models.LoadStatistics{}, "LoadTime"))
			if diff != "" {
				t.Fatal("LoadStatistics mismatch\n", diff)
			}
		})
	}
}

func TestIsValidIPAddress(t *testing.T) {
	tests := []struct {
		name      string
		ipAddress string
		want      bool
	}{
		{
			name:      "Correct IPv4 address",
			ipAddress: "70.95.73.73",
			want:      true,
		},
		{
			name:      "Correct IPv6 address 1",
			ipAddress: "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			want:      true,
		},
		{
			name:      "Correct IPv6 address 2",
			ipAddress: "2001:db8:85a3::8a2e:370:7334",
			want:      true,
		},
		{
			name:      "Empty IP address",
			ipAddress: "",
			want:      false,
		},
		{
			name:      "Error IPv4 address",
			ipAddress: "70.95.73.7399",
			want:      false,
		},
		{
			name:      "Error IPv6 address",
			ipAddress: "2001:ZZ8:85a3::8a2e:370:7334",
			want:      false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := IsValidIPAddress(tt.ipAddress); got != tt.want {
				t.Errorf("IsValidIPAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addTimeDurations(t *testing.T) {
	tests := []struct {
		name    string
		time1   string
		time2   string
		want    string
		wantErr string
	}{
		{
			name:    "Check errors",
			time1:   "aaaaa",
			time2:   "bbbb",
			want:    "cccc",
			wantErr: "time: invalid duration",
		},
		{
			name:    "Check valid duration",
			time1:   "3.225433834s",
			time2:   "4m17.929842458s",
			want:    "4m21.155276292s",
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := addTimeDurations(tt.time1, tt.time2)
			if len(tt.wantErr) > 0 {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				assert.Equalf(t, tt.want, got, "addTimeDurations(%v, %v)", tt.time1, tt.time2)
			}
		})
	}
}

func TestExtractIPAddress(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
	}{
		{
			name:   "Valid IP address",
			source: "api/geolocation/70.95.73.73",
			want:   "70.95.73.73",
		},
		{
			name:   "Empty IP address",
			source: "api/geolocation/",
			want:   "",
		},
		{
			name:   "Error IP address",
			source: "api/geolocation/70.95.73.739",
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ExtractIPAddress(tt.source), "ExtractIPAddress(%v)", tt.source)
		})
	}
}

func readFixture(t *testing.T, name string) string {
	t.Helper()

	_, curFile, _, ok := runtime.Caller(0)
	require.True(t, ok)

	path := filepath.Join(filepath.Dir(curFile), "testdata", name)

	return path
}
