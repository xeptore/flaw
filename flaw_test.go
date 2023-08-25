package flaw_test

import (
	"os"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/xeptore/flaw/v2"
)

type ExpectedRecord struct {
	Key     string         `json:"key"`
	Payload map[string]any `json:"payload"`
}

func parseRecordPayload(t *testing.T, p []byte) (v map[string]any) {
	assert.True(t, gjson.ValidBytes(p), "expected record payload to be a valid serialized json string, got: %s", string(p))
	result := gjson.ParseBytes(p)
	assert.True(t, result.IsObject(), "expected record payload to be a json object, got type: %q", result.Type.String())
	assert.NoError(t, json.Unmarshal(p, &v), "failed to unmarshal record payload from json to map: %s", string(p))
	return
}

func requireJSONEq(t *testing.T, expected map[string]any, recordPayload []byte) {
	actual := parseRecordPayload(t, recordPayload)
	expectedBytes, err := json.Marshal(expected)
	require.NoError(t, err, "failed to marshal expected json value: %#+v", expected)
	actualBytes, err := json.Marshal(actual)
	require.NoError(t, err, "failed to marshal record payload json value: %#+v", actual)
	require.Exactlyf(t, string(expectedBytes), string(actualBytes), "expected two json objects to have equal values")
}

func requireErrEq(t *testing.T, expectedRecords []ExpectedRecord, f error) {
	require.NotNil(t, f)
	expectedBytes, err := json.Marshal(expectedRecords)
	require.NoError(t, err, "failed to marshal expected records value to json: %#+v", expectedRecords)
	var parsedErrMessage []map[string]any
	errMsg := f.Error()
	require.NoError(t, json.Unmarshal([]byte(errMsg), &parsedErrMessage), "failed to unmarshal error message string from json to map: %s", errMsg)
	errMessageBytes, err := json.Marshal(parsedErrMessage)
	require.NoError(t, err, "failed to marshal error message json: %#+v", parsedErrMessage)
	require.Exactlyf(t, string(expectedBytes), string(errMessageBytes), "expected two json objects to have equal values")
	flawErr := new(flaw.Flaw)
	require.ErrorAs(t, f, &flawErr)
	for i, r := range flawErr.Records {
		require.Exactly(t, expectedRecords[i].Key, r.Key)
		requireJSONEq(t, expectedRecords[i].Payload, r.Payload)
	}
}

func TestNew(t *testing.T) {
	t.Parallel()
	f := flaw.New(
		"failed to connect to database",
		flaw.Dict("db").
			Str("host", "localhost").
			Int("port", 5643).
			Str("username", "root").
			Str("password", "root"),
	)
	expectedRecords := []ExpectedRecord{
		{
			Key: "db",
			Payload: map[string]any{
				"error":    "failed to connect to database",
				"host":     "localhost",
				"password": "root",
				"port":     5643,
				"username": "root",
			},
		},
	}
	requireErrEq(t, expectedRecords, f)
}

func TestFrom(t *testing.T) {
	t.Parallel()
	t.Run("Existing", testFromExisting)
	t.Run("NonExisting", testFromNonExisting)
}

func testFromNonExisting(t *testing.T) {
	t.Parallel()
	err := flaw.From(
		os.ErrClosed,
		"failed to connect to database",
		flaw.Dict("db").
			Str("host", "localhost").
			Int("port", 5643).
			Str("username", "root").
			Str("password", "root"),
	)
	expectedRecords := []ExpectedRecord{
		{
			Key: "db",
			Payload: map[string]any{
				"error":    "failed to connect to database: file already closed",
				"host":     "localhost",
				"password": "root",
				"port":     5643,
				"username": "root",
			},
		},
	}
	requireErrEq(t, expectedRecords, err)
}

func testFromExisting(t *testing.T) {
	t.Parallel()
	err := flaw.New(
		"failed to connect to database",
		flaw.Dict("db").
			Str("host", "localhost").
			Int("port", 5643).
			Str("username", "root").
			Str("password", "root"),
	)
	expectedRecords := []ExpectedRecord{
		{
			Key: "db",
			Payload: map[string]any{
				"error":    "failed to connect to database",
				"host":     "localhost",
				"password": "root",
				"port":     5643,
				"username": "root",
			},
		},
	}
	requireErrEq(t, expectedRecords, err)

	err = flaw.From(
		err,
		"failed to create user: permission denied",
		flaw.Dict("api").
			Str("request_id", "8fbbb51f-6f3a-4c9d-885a-92eb8e09cc31").
			Str("time", "2023-08-25T04:44:41.059Z").
			Str("client_ip", "127.0.0.1").
			Int("client_port", 58763),
	)
	expectedRecords = append(
		expectedRecords,
		ExpectedRecord{
			Key: "api",
			Payload: map[string]any{
				"error":       "failed to create user: permission denied",
				"time":        "2023-08-25T04:44:41.059Z",
				"request_id":  "8fbbb51f-6f3a-4c9d-885a-92eb8e09cc31",
				"client_port": 58763,
				"client_ip":   "127.0.0.1",
			},
		},
	)
	requireErrEq(t, expectedRecords, err)
}
