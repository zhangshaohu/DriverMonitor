package influx

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/influxdata/influxdb/models"
)

// Record describes the schema for an observation from
// our sensors.
type Record struct {
	Timestamp     time.Time
	Token         string
	HeartRate     float64
	Lat           float64
	Lng           float64
	Speed         float64
	FuelRemaining float64
	RPM           float64
}

// nolint: gocyclo
func populateRecord(row models.Row, valNdx int) (Record, error) {
	record := Record{
		HeartRate:     -999,
		Lat:           -999,
		Lng:           -999,
		Speed:         -999,
		FuelRemaining: -999,
		RPM:           -999,
	}
	var err error
	for i, col := range row.Columns {
		val := row.Values[valNdx][i]
		var convertedVal json.Number
		var ok bool
		switch col {
		case "time":
			var parsedNs int64
			parsedNs, err = (val.(json.Number).Int64())
			parsedTimestamp := time.Unix(0, parsedNs).UTC()
			record.Timestamp = parsedTimestamp
		case "RPM":
			var rpm float64
			convertedVal, ok = val.(json.Number)
			if ok {
				rpm, err = convertedVal.Float64()
				record.RPM = rpm
			}
		case "FuelRemaining":
			var fuelRemaining float64
			convertedVal, ok = val.(json.Number)
			if ok {
				fuelRemaining, err = convertedVal.Float64()
				record.FuelRemaining = fuelRemaining
			}
		case "HeartRate":
			var heartRate float64
			convertedVal, ok = val.(json.Number)
			if ok {
				heartRate, err = convertedVal.Float64()
				record.HeartRate = heartRate
			}
		case "Lat":
			var lat float64
			convertedVal, ok = val.(json.Number)
			if ok {
				lat, err = convertedVal.Float64()
				if lat != 0 {
					record.Lat = lat
				}
			}
		case "Lng":
			var lng float64
			convertedVal, ok = val.(json.Number)
			if ok {
				lng, err = convertedVal.Float64()
				if lng != 0 {
					record.Lng = lng
				}
			}
		case "Speed":
			var speed float64
			convertedVal, ok = val.(json.Number)
			if ok {
				speed, err = convertedVal.Float64()
				record.Speed = speed
			}
		case "Token":
			record.Token = val.(string)
		}
		if err != nil {
			return Record{}, err
		}
	}
	return record, nil
}

// AbstractClient defines the behavior an influx client should have; useful for
// allowing mock clients for testing
type AbstractClient interface {
	GetRecords(startTime time.Time, endTime time.Time, token string) ([]Record, error)
}

// Client is our wrapper for the influxDB client, which
// we'll use to send queries.
type Client struct {
	influxClient client.Client
	db           string
}

// NewClient creates a Client object connected to the influxDB
// instance at the given host and port.
func NewClient(dbHost string, dbPort int, db string) (*Client, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: fmt.Sprintf("http://%s:%d", dbHost, dbPort),
	})

	if err != nil {
		return nil, err
	}

	return &Client{
		influxClient: c,
		db:           db,
	}, nil
}

func (c *Client) influxQuery(query string, params map[string]interface{}) ([]client.Result, error) {
	queryObj := client.NewQueryWithParameters(query, c.db, "n", params)

	var res []client.Result
	if response, err := c.influxClient.Query(queryObj); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

// GetRecords returns all the records for the given time interval.
func (c *Client) GetRecords(startTime time.Time, endTime time.Time, token string) ([]Record, error) {
	results, err := c.influxQuery("SELECT * FROM records WHERE time >= $startTime AND time <= $endTime AND Token = $token;",
		map[string]interface{}{"startTime": startTime, "endTime": endTime, "token": token})
	if err != nil {
		return nil, err
	}

	records := []Record{}

	for _, row := range results[0].Series {
		for i := 0; i < len(row.Values); i++ {
			record, err := populateRecord(row, i)
			if err != nil {
				return nil, err
			}
			records = append(records, record)
		}
	}

	return records, nil
}
