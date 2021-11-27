package influx

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDB = "__test_csc591"

func getTestClient() (*Client, error) {
	return NewClient("influxdb", 8086, testDB)
}

func init() {
	client, err := getTestClient()
	if err != nil {
		log.Fatal(err)
	}
	_, _ = client.influxQuery(
		fmt.Sprintf("DROP DATABASE %s;", testDB),
		map[string]interface{}{})

	_, err = client.influxQuery(
		fmt.Sprintf("CREATE DATABASE %s;", testDB),
		map[string]interface{}{})

	if err != nil {
		log.Fatal(err)
	}
}

func insertRecord(c client.Client, record Record) {
	batch, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: testDB,
	})

	if err != nil {
		log.Fatal(err)
	}

	tags := map[string]string{"Token": record.Token}
	fields := map[string]interface{}{
		"HeartRate":     record.HeartRate,
		"Lat":           record.Lat,
		"Lng":           record.Lng,
		"Speed":         record.Speed,
		"FuelRemaining": record.FuelRemaining,
		"RPM":           record.RPM,
	}

	pt, err := client.NewPoint("records", tags, fields, record.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	batch.AddPoint(pt)

	if err := c.Write(batch); err != nil {
		log.Fatal(err)
	}
}

func TestClientConnects(t *testing.T) {
	_, err := NewClient("influxdb", 8086, testDB)
	require.Nil(t, err)
}

func TestGetRecords(t *testing.T) {
	client, err := getTestClient()
	require.Nil(t, err)

	minTime := time.Now().Add(-time.Minute * 5)
	maxTime := time.Now().Add(time.Minute * 5)

	testToken := "abcd"
	otherToken := "efgh"

	// No records in the DB
	records, err := client.GetRecords(minTime, maxTime, testToken)
	require.Nil(t, err)
	assert.Equal(t, 0, len(records))

	record1 := Record{
		Timestamp:     time.Now().UTC(),
		Token:         testToken,
		HeartRate:     120,
		Lat:           36.005110,
		Lng:           -78.927760,
		Speed:         45,
		FuelRemaining: 50,
		RPM:           5,
	}

	record2 := Record{
		Timestamp:     time.Now().Add(time.Minute).UTC(),
		Token:         testToken,
		HeartRate:     130,
		Lat:           35.995738,
		Lng:           -78.904028,
		Speed:         35,
		FuelRemaining: 25,
		RPM:           1,
	}

	record3 := Record{
		Timestamp: time.Now().UTC(),
		Token:     otherToken,
	}

	insertRecord(client.influxClient, record1)

	// One record in the DB (correct token)
	records, err = client.GetRecords(minTime, maxTime, testToken)
	require.Nil(t, err)
	require.Equal(t, 1, len(records))
	assert.Equal(t, record1, records[0])

	insertRecord(client.influxClient, record2)

	// Two records in the DB (correct tokens)
	records, err = client.GetRecords(minTime, maxTime, testToken)
	require.Nil(t, err)
	require.Equal(t, 2, len(records))
	assert.Equal(t, record1, records[0])
	assert.Equal(t, record2, records[1])

	insertRecord(client.influxClient, record3)

	// 3 records in the DB (one has other token)
	records, err = client.GetRecords(minTime, maxTime, testToken)
	require.Nil(t, err)
	require.Equal(t, 2, len(records))
	assert.Equal(t, record1, records[0])
	assert.Equal(t, record2, records[1])

	// 1 record for the other token
	records, err = client.GetRecords(minTime, maxTime, otherToken)
	require.Nil(t, err)
	require.Equal(t, 1, len(records))
	assert.Equal(t, record3, records[0])
}
