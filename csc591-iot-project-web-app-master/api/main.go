package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.ncsu.edu/jmnance/csc591-iot-project-web-app/api/influx"
	"github.ncsu.edu/jmnance/csc591-iot-project-web-app/api/view"
)

func main() {
	strDBPort := os.Getenv("INFLUXDB_PORT")
	dbPort, err := strconv.Atoi(strDBPort)

	if err != nil {
		log.Fatalf("failed to parse db port: %s", strDBPort)
	}
	dbHost := os.Getenv("INFLUXDB_HOST")
	dbDB := os.Getenv("INFLUXDB_DB")
	db, err := influx.NewClient(dbHost, dbPort, dbDB)
	if err != nil {
		log.Fatalf("failed to connect to InfluxDB at host '%s', port '%d', db '%s'",
			dbHost, dbPort, dbDB)
	}

	viewContext := view.NewContext(db)

	log.Fatal(http.ListenAndServe(":8000", view.Router(viewContext)))
}
