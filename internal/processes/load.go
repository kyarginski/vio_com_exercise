package processes

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"vio/internal/models"
)

// Regular expression to check IPv4 and IPv6.
var ipRegex = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}$|^([0-9a-fA-F]{1,4}:){1,7}:([0-9a-fA-F]{1,4}:){1,7}[0-9a-fA-F]{1,4}$|^([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}$|^([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}$|^([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}$|^([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}$|^([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}$|^([0-9a-fA-F]{1,4}:){1,1}(:[0-9a-fA-F]{1,4}){1,6}$|^:(:[0-9a-fA-F]{1,4}){1,7}$|^::$`)

func loadData(path string) ([]models.Location, *models.LoadStatistics, error) {
	startTime := time.Now()
	fileInfos, err := os.ReadDir(path)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading directory: %v", err)
	}

	var errs []error
	uniqLocations := make(map[string]models.Location)
	locations := make([]models.Location, 0)
	var loadStatistics models.LoadStatistics
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() && filepath.Ext(fileInfo.Name()) == ".csv" {
			filePath := filepath.Join(path, fileInfo.Name())
			file, err := os.Open(filePath)
			if err != nil {
				errs = append(errs, fmt.Errorf("error opening file %s: %v", fileInfo.Name(), err))

				continue
			}
			loadStatistics.FilesCount++

			reader := csv.NewReader(file)

			// Skip the header of the CSV file.
			_, err = reader.Read()
			if err != nil {
				errs = append(errs, err)

				break
			}

			for {
				record, err := reader.Read()
				if err != nil {
					if !errors.Is(err, io.EOF) {
						errs = append(errs, err)
					}

					break
				}

				latitude, _ := strconv.ParseFloat(record[4], 64)
				longitude, _ := strconv.ParseFloat(record[5], 64)
				mysteryValue, _ := strconv.ParseInt(record[6], 10, 64)

				location := models.Location{
					IPAddress:    record[0],
					CountryCode:  record[1],
					Country:      record[2],
					City:         record[3],
					Latitude:     latitude,
					Longitude:    longitude,
					MysteryValue: mysteryValue,
				}

				if !IsValidIPAddress(location.IPAddress) {
					loadStatistics.Discarded++
					continue
				}

				uniqLocations[location.IPAddress] = location
				loadStatistics.Accepted++
			}

			errClose := file.Close()
			if errClose != nil {
				errs = append(errs, fmt.Errorf("error closing file %s: %v", fileInfo.Name(), errClose))
			}
		}
	}

	for _, location := range uniqLocations {
		locations = append(locations, location)
	}

	err = errors.Join(errs...)
	loadStatistics.Total = loadStatistics.Accepted + loadStatistics.Discarded
	loadStatistics.LoadTime = time.Since(startTime).String()

	return locations, &loadStatistics, err
}

func IsValidIPAddress(ipAddress string) bool {
	return ipRegex.MatchString(ipAddress)
}

func ExtractIPAddress(s string) string {
	re := regexp.MustCompile(`\b(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`)
	match := re.FindString(s)
	if match == "" {
		return ""
	}
	return match
}

func addTimeDurations(time1, time2 string) (string, error) {
	duration1, err := time.ParseDuration(time1)
	if err != nil {
		return "", err
	}
	duration2, err := time.ParseDuration(time2)
	if err != nil {
		return "", err
	}
	// Addition of time intervals.
	result := duration1 + duration2

	resultString := result.String()

	return resultString, nil
}
