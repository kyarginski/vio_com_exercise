package processes

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"vio/internal/models"
)

// Worker represents a goroutine that processes records.
type Worker struct {
	ID       int
	JobQueue chan models.Location
	Result   chan error
}

func (w *Worker) start(ctx context.Context, wg *sync.WaitGroup, stmt *sql.Stmt, mu *sync.Mutex, errs *[]error) {
	defer wg.Done()

	for loc := range w.JobQueue {
		_, err := stmt.ExecContext(ctx, loc.IPAddress, loc.CountryCode, loc.Country, loc.City, loc.Latitude, loc.Longitude, loc.MysteryValue)
		if err != nil {
			mu.Lock()
			*errs = append(*errs, err)
			mu.Unlock()
		}
	}

	w.Result <- nil
}

func processParallelWithMaxProcs(ctx context.Context, db *sql.DB, locations *[]models.Location, maxWorkers int) (string, error) {
	startTime := time.Now()
	stmt, err := db.Prepare(SQLInsert)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	// Create channels for tasks and results.
	jobQueue := make(chan models.Location, len(*locations))
	results := make(chan error, maxWorkers)

	for i := 1; i <= maxWorkers; i++ {
		worker := &Worker{
			ID:       i,
			JobQueue: jobQueue,
			Result:   results,
		}
		wg.Add(1)
		go worker.start(ctx, &wg, stmt, &mu, &errs)
	}

	// Filling the channel with tasks.
	go func() {
		for _, loc := range *locations {
			jobQueue <- loc
		}
		close(jobQueue)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result != nil {
			return "", result
		}
	}

	finishTime := time.Since(startTime).String()

	return finishTime, errors.Join(errs...)
}
