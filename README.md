# FindHotel Coding Challenge

## Geolocation Service

### Overview

You are provided with a CSV file (`data_dump.csv`) that contains raw geolocation data. The goal is to develop a service that imports such data and expose it via an API.

Sample data:
```
ip_address,country_code,country,city,latitude,longitude,mystery_value
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115
70.95.73.73,TL,Saudi Arabia,Gradymouth,-49.16675918861615,-86.05920084416894,2559997162
,PY,Falkland Islands (Malvinas),,75.41685191518815,-144.6943217219469,0
125.159.20.54,LI,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276
```

### Requirements

1. Develop a library/component with two main features:
    * A service that parses the CSV file containing the raw data and persists it in a database;
    * An interface to provide access to the geolocation data (model layer);
1. Develop a REST API that uses the aforementioned library to expose the geolocation data.

In doing so:
* Define a data format suitable for the data contained in the CSV file;
* Sanitise the entries: the file comes from an unreliable source, this means that the entries can be duplicated, may miss some value, the value can not be in the correct format or completely bogus;
* At the end of the import process, return some statistics about the time elapsed, as well as the number of entries accepted/discarded;
* The library should be configurable by an external configuration (particularly with regards to the DB configuration);
* The API layer should implement a single HTTP endpoint that, given an IP address, returns information about the IP address' location (e.g. country, city).

### Expected outcome and shipping:

* A library/component that packages the import service and the interface for accessing the geolocation data;
* A REST API application that uses the aforementioned library

### Notes

* The file's contents are fake, you do not have to worry about data correctness.
* For running the application locally (development) you can decide to include a docker-compose, or setup a Makefile.
* You can structure the repository as you see it fit.

### Evaluation

The following are the criteria we will use to evaluate your work:
- Code quality;
- Well-tested solution;
- Best code practices in general;
- Architectural design skills;
- API design and data structure skills (i.e. correctness of the API responses, etc);
- Communication and writing skills (i.e. documentation on how the apps works, setup and tradeoffs).

---

# Solution

## Criteria for checking data from a file

- The file is not empty
- The `ip_address` field is not empty and is a valid IP address
- Any other fields are not checked and may have empty values
- Duplicate processing strategy: a newer entry replaces the previous one (subject to IP address verification)

## Run service as CLI application (loader)

```shell
go run ./cmd/loader -source=data_source -parallel=0
```

## Run service as server application (geolocation)

```shell
go run ./cmd/geolocation
```


## Run service as docker container

```shell
docker-compose up -d
```

## Stop service with docker container

```shell
docker-compose down
```

### Results of testing

Run with 1 goroutine.

```shell
go run ./cmd/loader -source=data_source

finished reading: 3.225433834s
finished db inserting: 4m17.929842458s
{"load_time":"4m21.155276292s","files_count":2,"accepted":966482,"discarded":33524,"total":1000006}
```

Run with 30 parallel goroutines.

```shell
go run ./cmd/loader -source=data_source -parallel=30

finished reading: 3.218792959s
start with 30 parallel working
finished db inserting: 47.593757125s
{"load_time":"50.812550084s","files_count":2,"accepted":966482,"discarded":33524,"total":1000006}
```

#### So the parallel processing is 5 times faster than sequential processing.

## Get result

Get result with Postman

```shell
GET http://localhost:8087/api/geolocation/70.95.73.73
```

Result 200 OK

body:
```json
{
    "ip_address": "70.95.73.73",
    "country_code": "TL",
    "country": "Saudi Arabia",
    "city": "Gradymouth",
    "latitude": -49.16675918861615,
    "longitude": -86.05920084416894,
    "mystery_value": 2559997162
}
```

## Unit tests

```shell
 make test
```

## Integration tests

```shell
make integration_test
```


## Test coverage

```shell
go test -cover ./...
```

#### Test coverage is 69.9% of statements.

## Documentation of API (swagger)

See swagger documentation file [here](doc/swagger.json)
