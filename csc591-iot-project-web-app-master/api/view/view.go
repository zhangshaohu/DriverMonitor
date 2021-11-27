package view

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.ncsu.edu/jmnance/csc591-iot-project-web-app/api/influx"
)

// Router returns the router used to handle all our requests.
func Router(viewContext *Context) *mux.Router {
	router := mux.NewRouter()
	router.Path("/api/records/{token}").
		Queries("startTime", "{startTime}", "endTime", "{endTime}").
		HandlerFunc(viewContext.GetRecords)

	return router
}

// Context wraps app-level context for the views, providing
// access to a persistent database connection.
type Context struct {
	db influx.AbstractClient
}

// NewContext creates a new Context object using the given parameters.
func NewContext(db influx.AbstractClient) *Context {
	return &Context{db: db}
}

// GetRecords returns a JSON array of records for the given startTime-endTime
// range and the given token.
func (c *Context) GetRecords(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var err error

	startTime, err := time.Parse(time.RFC3339, params["startTime"])

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse startTime: %s", err), 400)
		return
	}

	endTime, err := time.Parse(time.RFC3339, params["endTime"])

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse endTime: %s", err), 400)
		return
	}

	token := params["token"]

	records, err := c.db.GetRecords(startTime, endTime, token)

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to query records: %s", err), 400)
		return
	}

	err = json.NewEncoder(w).Encode(records)

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write records to JSON: %s", err), 500)
		return
	}
}
