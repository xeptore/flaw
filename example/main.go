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
	"github.com/xeptore/flaw/v2"
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
//   },
//   "stack_traces": [
//     {
//       "location": "./main.go:36",
//       "function": "main.insertRedisKey"
//     },
//     {
//       "location": "./main.go:51",
//       "function": "main.createUser"
//     },
//     {
//       "location": "./main.go:134",
//       "function": "main.main.func1"
//     },
//     {
//       "location": "net/http/server.go:2136",
//       "function": "net/http.HandlerFunc.ServeHTTP"
//     },
//     {
//       "location": "net/http/server.go:2514",
//       "function": "net/http.(*ServeMux).ServeHTTP"
//     },
//     {
//       "location": "net/http/server.go:2938",
//       "function": "net/http.serverHandler.ServeHTTP"
//     },
//     {
//       "location": "net/http/server.go:2009",
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
