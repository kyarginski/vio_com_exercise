package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"vio/internal/processes"
)

var (
	sourceFlag   string
	databaseFlag string
	parallelFlag string
	helpFlag     string
)

var errUsage = errors.New("usage")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s [option...] 

Options:
`,
			os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&helpFlag, "help", "", "show help about arguments")
	flag.StringVar(&sourceFlag, "source", "data_source", "directory of input data files")
	flag.StringVar(&databaseFlag, "database", "host=localhost port=25432 dbname=postgres user=postgres password=postgres client_encoding=UTF8 sslmode=disable", "connection string to database")
	flag.StringVar(&parallelFlag, "parallel", "0", "parallel processing -1 = off, N = count of goroutines, 0 = count of CPU cores")
	flag.Parse()
	if len(helpFlag) > 0 {
		flag.Usage()
		os.Exit(2)
	}

	info, err := run()
	if err != nil {
		if errors.Is(err, errUsage) {
			flag.Usage()
			os.Exit(2)
		}
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(info))
}

func run() ([]byte, error) {
	return processes.RunOnce(sourceFlag, databaseFlag, parallelFlag)
}
