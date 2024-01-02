package flaw_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xeptore/flaw/v8"
)

func TestNew(t *testing.T) {
	t.Parallel()
	err := flaw.
		From(fmt.Errorf("db: failed to connect to database: %v", os.ErrPermission)).
		Append(flaw.P{
			"host":     "localhost",
			"port":     5643,
			"username": "root",
			"password": "root",
		})
	assert.Panics(t, func() { err.Append(nil) })
	callerErr := func() error {
		if flawErr := new(flaw.Flaw); errors.As(err, &flawErr) {
			return flawErr.Append(
				flaw.P{
					"artist": "Ramin Djawadi",
					"year":   2012,
				},
				flaw.P{
					"sql": flaw.P{
						"query": "select * from artists",
					},
				},
			)
		} else {
			assert.FailNow(t, "expected flaw error to pass errors.As, but failed")
			return nil
		}
	}()
	assert.Exactly(t, "db: failed to connect to database: permission denied", err.Error())
	assert.Exactly(t, "db: failed to connect to database: permission denied", err.Inner)
	assert.Truef(t, len(err.StackTrace) > 0, "expected flaw stack trace not to be empty")
	assert.Len(t, err.Records, 2)
	assert.Exactly(t, "github.com/xeptore/flaw/v8_test.TestNew", err.Records[0].Function)
	assert.Exactly(t, flaw.P{"host": "localhost", "port": 5643, "username": "root", "password": "root"}, err.Records[0].Payload)
	assert.NotNil(t, callerErr)
	assert.Exactly(t, "github.com/xeptore/flaw/v8_test.TestNew.func2", err.Records[1].Function)
	assert.Exactly(t, flaw.P{"artist": "Ramin Djawadi", "year": 2012, "sql": flaw.P{"query": "select * from artists"}}, err.Records[1].Payload)
}

func TestFrom(t *testing.T) {
	t.Parallel()
	t.Run("ShouldPanicOnNilArg", func(t *testing.T) {
		t.Parallel()
		assert.Panics(t, func() { flaw.From(nil) })
	})
}
