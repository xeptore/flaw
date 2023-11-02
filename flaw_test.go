package flaw_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xeptore/flaw/v6"
)

func TestNew(t *testing.T) {
	t.Parallel()
	err := flaw.
		From(fmt.Errorf("db: failed to connect to database: %v", os.ErrPermission)).
		Append(map[string]any{
			"host":     "localhost",
			"port":     5643,
			"username": "root",
			"password": "root",
		})
	callerErr := func() error {
		if flawErr := new(flaw.Flaw); errors.As(err, &flawErr) {
			return flawErr.Append(map[string]any{
				"artist": "Ramin Djawadi",
				"year":   2012,
			})
		} else {
			require.FailNow(t, "expected flaw error to pass errors.As, but failed")
			return nil
		}
	}()
	require.Exactly(t, "db: failed to connect to database: permission denied", err.Error())
	require.Exactly(t, "db: failed to connect to database: permission denied", err.Inner)
	require.Truef(t, len(err.StackTrace) > 0, "expected flaw stack trace not to be empty")
	require.Len(t, err.Records, 2)
	require.Exactly(t, err.Records[0].Function, "command-line-arguments_test.TestNew")
	require.Exactly(t, err.Records[0].Payload, map[string]any{"host": "localhost", "port": 5643, "username": "root", "password": "root"})
	require.NotNil(t, callerErr)
	require.Exactly(t, err.Records[1].Function, "command-line-arguments_test.TestNew.func1")
	require.Exactly(t, err.Records[1].Payload, map[string]any{"artist": "Ramin Djawadi", "year": 2012})
}
