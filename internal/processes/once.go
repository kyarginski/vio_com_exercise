package processes

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"vio/internal/database"
)

func RunOnce(path string, connectString string, parallelFlag string) ([]byte, error) {
	locations, loadStatistics, err := loadData(path)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stderr, "finished reading: %s\n", loadStatistics.LoadTime)
	var loadTimeProcessStr string

	db, err := database.GetDB(connectString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	ctx := context.Background()

	if parallelFlag == "-1" {
		fmt.Fprintf(os.Stderr, "no parallel working\n")
		loadTimeProcessStr, err = process(ctx, db, &locations)
	} else {
		// Getting the number of processors in the system.
		maxWorkers, _ := strconv.Atoi(parallelFlag)
		if maxWorkers == 0 {
			maxWorkers = runtime.NumCPU()
		}
		fmt.Fprintf(os.Stderr, "start with %d parallel working\n", maxWorkers)

		loadTimeProcessStr, err = processParallelWithMaxProcs(ctx, db, &locations, maxWorkers)
	}

	if err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stderr, "finished db inserting: %s\n", loadTimeProcessStr)
	loadStatistics.LoadTime, err = addTimeDurations(loadStatistics.LoadTime, loadTimeProcessStr)
	if err != nil {
		return nil, err
	}
	jsonStatistics, err := json.Marshal(loadStatistics)
	if err != nil {
		return nil, err
	}

	return jsonStatistics, nil
}
