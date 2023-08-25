package main

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/xeptore/flaw/v2"
)

var (
	log = zerolog.New(os.Stderr)
)

func insertRedisKey(key string, value string) error {
	if key == "bad-key" {
		return flaw.New(
			"attempt to insert a bad key into redis",
			flaw.Dict("redis").Str("key", key).Str("value", value),
		)
	}

	if time.Now().Day()%2 == 0 {
		return errors.New("bad day error")
	}

	return nil
}

func createUser(userID string, age int, isAdmin bool) error {
	if age > 40 {
		if err := insertRedisKey("bad-key", userID); nil != err {
			return flaw.From(
				err,
				"failed to insert user into redis",
				flaw.Dict("user").Str("id", userID).Int("age", age).Bool("is_admin", isAdmin),
			)
		}
	}
	return nil
}

func logErr(err error) {
	log.
		Error().
		Func(
			func(e *zerolog.Event) {
				if flawErr := new(flaw.Flaw); errors.As(err, &flawErr) {
					dict := zerolog.Dict()
					for _, v := range flawErr.Records {
						dict.RawJSON(v.Key, v.Payload)
					}
					e.Dict("info", dict)
					return
				}
				e.Err(err)
			},
		).
		Send()
}

// Results in the compacted version of the following JSON object:
//
// {
//   "level": "error",
//   "info": {
//     "redis": {
//       "key": "bad-key",
//       "value": "a",
//       "error": "attempt to insert a bad key into redis"
//     },
//     "user": {
//       "id": "a",
//       "age": 42,
//       "is_admin": true,
//       "error": "failed to insert user into redis"
//     }
//   }
// }

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/create-user") {
			if err := createUser("a", 42, true); nil != err {
				w.WriteHeader(http.StatusInternalServerError)
				logErr(err)
				return
			}
			w.WriteHeader(http.StatusCreated)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
	http.ListenAndServe("127.0.0.1:8080", mux)
}
