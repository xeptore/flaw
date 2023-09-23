package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/tidwall/pretty"
	"github.com/xeptore/flaw/v5"
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
			New("redis: attempt to insert a bad key into redis").
			With(flaw.NewDict().Str("key", key).Str("value", value))
	}

	if time.Now().Day()%2 == 0 {
		return errors.New("bad day error")
	}

	return nil
}

func createUser(userID string, age int, isAdmin bool) error {
	if age > 40 {
		if err := insertRedisKey("bad-key", userID); nil != err {
			return flaw.
				From(err, "user: failed to insert user into redis").
				With(flaw.NewDict().Str("id", userID).Int("age", age).Bool("is_admin", isAdmin))
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
					records := zerolog.Arr()
					for _, v := range flawErr.Records {
						record := zerolog.Dict().Str("message", v.Message)
						if v.Payload == nil {
							record.RawJSON("payload", []byte("null"))
						} else {
							record.RawJSON("payload", v.Payload)
						}
						records.Dict(record)
					}
					e.Array("records", records)
					e.Func(func(e *zerolog.Event) {
						arr := zerolog.Arr()
						lo.ForEach(flawErr.Traces, func(v flaw.StackTrace, _ int) {
							arr.Dict(zerolog.Dict().Str("location", fmt.Sprintf("%s:%d", v.File, v.Line)).Str("function", v.Function))
						})
						e.Array("stack_traces", arr)
					})
					return
				}
				e.Err(err)
			},
		).
		Send()
}

// Will print the following JSON object:
// {
//   "level": "error",
//   "records": [
//     {
//       "message": "redis: attempt to insert a bad key into redis",
//       "payload": {
//         "key": "bad-key",
//         "value": "a"
//       }
//     },
//     {
//       "message": "user: failed to insert user into redis",
//       "payload": {
//         "id": "a",
//         "age": 42,
//         "is_admin": true
//       }
//     },
//     {
//       "message": "http: failed to process request",
//       "payload": {
//         "request": {
//           "headers": {
//             "Accept": ["*/*"],
//             "User-Agent": ["curl/8.3.0"]
//           }
//         }
//       }
//     }
//   ],
//   "stack_traces": [
//     {
//       "location": "/path/flaw/example/main.go:37",
//       "function": "main.insertRedisKey"
//     },
//     {
//       "location": "/path/flaw/example/main.go:50",
//       "function": "main.createUser"
//     },
//     {
//       "location": "/path/flaw/example/main.go:147",
//       "function": "main.main.func1"
//     },
//     {
//       "location": "/path/env/go/go/src/net/http/server.go:2136",
//       "function": "net/http.HandlerFunc.ServeHTTP"
//     },
//     {
//       "location": "/path/env/go/go/src/net/http/server.go:2514",
//       "function": "net/http.(*ServeMux).ServeHTTP"
//     },
//     {
//       "location": "/path/env/go/go/src/net/http/server.go:2938",
//       "function": "net/http.serverHandler.ServeHTTP"
//     },
//     {
//       "location": "/path/env/go/go/src/net/http/server.go:2009",
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
				dict := flaw.NewDict().Dict("request", flaw.NewDict().Dict("headers", flaw.NewDict().Func(func(d *flaw.Dict) {
					for k, v := range r.Header {
						d.Strs(k, v)
					}
				})))
				logErr(flaw.From(err, "http: failed to process request").With(dict))
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
