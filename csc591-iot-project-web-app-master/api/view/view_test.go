package view

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.ncsu.edu/jmnance/csc591-iot-project-web-app/api/influx"
)

var server *httptest.Server

const token1 = "abcd"
const token2 = "efgh"
const badToken = "zyxw"

var token1Records = []influx.Record{{
	Token:         token1,
	Timestamp:     time.Now().UTC(),
	HeartRate:     1,
	Lat:           1,
	Lng:           1,
	Speed:         1,
	FuelRemaining: 1,
	RPM:           1,
}}

var token2Records = []influx.Record{{
	Token:         token2,
	Timestamp:     time.Now().UTC(),
	HeartRate:     2,
	Lat:           2,
	Lng:           2,
	Speed:         2,
	FuelRemaining: 2,
	RPM:           2,
}, {
	Token:         token2,
	Timestamp:     time.Now().UTC(),
	HeartRate:     3,
	Lat:           3,
	Lng:           3,
	Speed:         3,
	FuelRemaining: 3,
	RPM:           3,
}}

var badTokenRecords = []influx.Record{}

type APITestClient struct{}

func (a *APITestClient) GetRecords(startTime time.Time, endTime time.Time, token string) ([]influx.Record, error) {
	switch token {
	case token1:
		return token1Records, nil
	case token2:
		return token2Records, nil
	case badToken:
		return badTokenRecords, nil
	}
	return nil, fmt.Errorf("unrecognized token: %s", token)
}

func init() {
	client := &APITestClient{}

	testContext := NewContext(client)
	server = httptest.NewServer(Router(testContext))
}

func getTestURL(token string) string {
	startTime := time.Time{}.Format(time.RFC3339)
	endTime := time.Now().Add(500 * time.Hour).Format(time.RFC3339)
	return fmt.Sprintf("%s/api/records/%s?startTime=%s&endTime=%s",
		server.URL, token, startTime, endTime)
}

func getRecordsFromResBody(res *http.Response) ([]influx.Record, error) {
	var records []influx.Record

	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(&records)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func TestGetRecords(t *testing.T) {
	// No records
	res, err := http.Get(getTestURL(badToken))
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	records, err := getRecordsFromResBody(res)
	require.Nil(t, err)
	assert.Equal(t, 0, len(records))

	// One record
	res, err = http.Get(getTestURL(token1))
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	records, err = getRecordsFromResBody(res)
	require.Nil(t, err)
	assert.Equal(t, 1, len(records))
	assert.Equal(t, token1Records, records)

	// Two records
	res, err = http.Get(getTestURL(token2))
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	records, err = getRecordsFromResBody(res)
	require.Nil(t, err)
	assert.Equal(t, 2, len(records))
	assert.Equal(t, token2Records, records)
}
