package flaw_test

import (
	"os"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/require"

	"github.com/xeptore/flaw/v4"
)

type ExpectedRecord struct {
	Message string         `json:"message"`
	Key     string         `json:"key"`
	Payload map[string]any `json:"payload"`
}

func requireErrEq(t *testing.T, expectedRecords []ExpectedRecord, f error) {
	require.NotNil(t, f)
	expectedBytes, err := json.Marshal(expectedRecords)
	require.NoError(t, err, "failed to marshal expected records value to json: %#+v", expectedRecords)
	require.JSONEq(t, string(expectedBytes), f.Error())
}

func TestNew(t *testing.T) {
	t.Parallel()
	f := flaw.
		New("failed to connect to database").
		With(
			flaw.NewDict("db").
				Str("host", "localhost").
				Int("port", 5643).
				Str("username", "root").
				Str("password", "root"),
		)
	expectedRecords := []ExpectedRecord{
		{
			Message: "failed to connect to database",
			Key:     "db",
			Payload: map[string]any{
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
	t.Run("NoRecord", testFromNoRecord)
}

func testFromNonExisting(t *testing.T) {
	t.Parallel()
	err := flaw.From(
		os.ErrClosed,
		"failed to connect to database",
	).
		With(
			flaw.NewDict("db").
				Str("host", "localhost").
				Int("port", 5643).
				Str("username", "root").
				Str("password", "root"),
		)
	expectedRecords := []ExpectedRecord{
		{
			Message: "failed to connect to database: file already closed",
			Key:     "db",
			Payload: map[string]any{
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
	err := flaw.
		New("failed to connect to database").
		With(
			flaw.NewDict("db").
				Str("host", "localhost").
				Int("port", 5643).
				Str("username", "root").
				Str("password", "root"),
		)
	expectedRecords := []ExpectedRecord{
		{
			Message: "failed to connect to database",
			Key:     "db",
			Payload: map[string]any{
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
	).
		With(
			flaw.NewDict("api").
				Str("request_id", "8fbbb51f-6f3a-4c9d-885a-92eb8e09cc31").
				Str("time", "2023-08-25T04:44:41.059Z").
				Str("client_ip", "127.0.0.1").
				Int("client_port", 58763),
		)
	expectedRecords = append(
		expectedRecords,
		ExpectedRecord{
			Message: "failed to create user: permission denied",
			Key:     "api",
			Payload: map[string]any{
				"time":        "2023-08-25T04:44:41.059Z",
				"request_id":  "8fbbb51f-6f3a-4c9d-885a-92eb8e09cc31",
				"client_port": 58763,
				"client_ip":   "127.0.0.1",
			},
		},
	)
	requireErrEq(t, expectedRecords, err)
}

func testFromNoRecord(t *testing.T) {
	t.Parallel()
	err := flaw.From(
		os.ErrClosed,
		"failed to connect to database",
	)
	expectedRecords := []ExpectedRecord{
		{
			Message: "failed to connect to database: file already closed",
			Key:     "",
			Payload: nil,
		},
	}
	requireErrEq(t, expectedRecords, err)
}
