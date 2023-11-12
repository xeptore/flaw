package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/tidwall/pretty"

	"github.com/xeptore/flaw/v8"
)

var (
	log = zerolog.New(NewPretty(os.Stderr))
)

func NewPretty(out io.Writer) Pretty {
	return Pretty{out: out}
}

type Pretty struct {
	out io.Writer
}

func (p Pretty) Write(line []byte) (int, error) {
	return os.Stderr.Write(pretty.Color(pretty.Pretty(line), nil))
}

func insertRedisKey(key string, value string) error {
	if key == "bad-key" {
		return flaw.
			From(errors.New("redis: attempt to insert a bad key into redis")).
			Append(map[string]any{"key": key, "value": value}, nil, map[string]any{"x": 2})
	}

	if time.Now().Day()%2 == 0 {
		return errors.New("bad day error")
	}

	return nil
}

func createUser(userID string, age int, isAdmin bool) error {
	if age > 40 {
		if err := insertRedisKey("bad-key", userID); nil != err {
			if flawErr := new(flaw.Flaw); errors.As(err, &flawErr) {
				return flawErr.Append(flaw.P{"id": userID, "age": age, "is_admin": isAdmin})
			}
			return flaw.From(fmt.Errorf("user: failed to insert user into redis: %v", err))
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
					e.Str("error", flawErr.Inner)
					records := zerolog.Arr()
					for _, v := range flawErr.Records {
						payload, err := json.Marshal(v.Payload)
						if nil != err {
							panic(fmt.Errorf("failed to marshal record paylod: %v", err))
						}
						records.Dict(zerolog.Dict().Str("function", v.Function).RawJSON("payload", payload))
					}
					e.Array("records", records)
					stackTrace := zerolog.Arr()
					for _, v := range flawErr.StackTrace {
						stackTrace.Dict(zerolog.Dict().Str("location", fmt.Sprintf("%s:%d", v.File, v.Line)).Str("function", v.Function))
					}
					e.Array("stack_traces", stackTrace)
					return
				}
				e.Err(err)
			},
		).
		Send()
}

// Will print the following JSON object:
//
// {
//   "level": "error",
//   "error": "redis: attempt to insert a bad key into redis",
//   "records": [
//     {
//       "function": "main.insertRedisKey",
//       "payload": {
//         "key": "bad-key",
//         "value": "a",
//         "x": 2
//       }
//     },
//     {
//       "function": "main.createUser",
//       "payload": {
//         "age": 42,
//         "id": "a",
//         "is_admin": true
//       }
//     }
//   ],
//   "stack_traces": [
//     {
//       "location": "cwd/example/main.go:38",
//       "function": "main.insertRedisKey"
//     },
//     {
//       "location": "cwd/example/main.go:51",
//       "function": "main.createUser"
//     },
//     {
//       "location": "cwd/example/main.go:149",
//       "function": "main.main.func1"
//     },
//     {
//       "location": "goroot/src/net/http/server.go:2136",
//       "function": "net/http.HandlerFunc.ServeHTTP"
//     },
//     {
//       "location": "goroot/src/net/http/server.go:2514",
//       "function": "net/http.(*ServeMux).ServeHTTP"
//     },
//     {
//       "location": "goroot/src/net/http/server.go:2938",
//       "function": "net/http.serverHandler.ServeHTTP"
//     },
//     {
//       "location": "goroot/src/net/http/server.go:2009",
//       "function": "net/http.(*conn).serve"
//     }
//   ]
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
	if err := http.ListenAndServe("127.0.0.1:8080", mux); nil != err {
		log.Fatal().Err(err).Msg("http server listener stopped")
	}
}
