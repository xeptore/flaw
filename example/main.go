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
	return os.Stderr.Write(pretty.Pretty(line))
}

func closeFile() error {
	return errors.New("wtf")
}

func insertRedisKey(key string, value string) (err error) {
	defer func() {
		if nil != err {
			if closeErr := closeFile(); nil != closeErr {
				if flawErr := new(flaw.Flaw); errors.As(err, &flawErr) {
					flawErr.Join(fmt.Errorf("failed to close file: %v", closeErr)).Append(flaw.P{"tty": "putty"}).Join(os.ErrPermission)
				}
			}
		}
	}()

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
					e.Dict(
						"error",
						zerolog.
							Dict().
							Str("message", flawErr.Inner).
							Str("typename", flawErr.InnerType).
							Str("syntax_representation", flawErr.InnerSyntaxRepr),
					)

					records := zerolog.Arr()
					for _, v := range flawErr.Records {
						payload, err := json.Marshal(v.Payload)
						if nil != err {
							panic(fmt.Errorf("failed to marshal record paylod: %v", err))
						}
						records.
							Dict(
								zerolog.
									Dict().
									Str("function", v.Function).
									RawJSON("payload", payload),
							)
					}
					e.Array("records", records)

					stackTrace := zerolog.Arr()
					for _, v := range flawErr.StackTrace {
						stackTrace.
							Dict(
								zerolog.
									Dict().
									Str("location", fmt.Sprintf("%s:%d", v.File, v.Line)).
									Str("function", v.Function),
							)
					}
					e.Array("stack_traces", stackTrace)

					joined := zerolog.Arr()
					for _, v := range flawErr.JoinedErrors {
						d := zerolog.Dict()
						if st := v.CallerStackTrace; nil != st {
							d.Dict(
								"caller_stack_trace",
								zerolog.
									Dict().
									Str("location", fmt.Sprintf("%s:%d", st.File, st.Line)).
									Str("function", st.Function),
							)
						}
						d.Dict(
							"error",
							zerolog.
								Dict().
								Str("message", v.Message).
								Str("type_name", v.TypeName).
								Str("syntax_representation", v.SyntaxRepr),
						)
						joined.Dict(d)
					}
					e.Array("joined_errors", joined)

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
//   "error": {
//     "message": "redis: attempt to insert a bad key into redis",
//     "typename": "*errors.errorString",
//     "syntax_representation": "&errors.errorString{s:\"redis: attempt to insert a bad key into redis\"}"
//   },
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
//       "function": "main.insertRedisKey.func1",
//       "payload": {
//         "tty": "putty"
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
//       "location": "cwd/flaw/example/main.go:52",
//       "function": "main.insertRedisKey"
//     },
//     {
//       "location": "cwd/flaw/example/main.go:65",
//       "function": "main.createUser"
//     },
//     {
//       "location": "cwd/flaw/example/main.go:228",
//       "function": "main.main.func1"
//     },
//     {
//       "location": "goroot/src/net/http/server.go:2220",
//       "function": "net/http.HandlerFunc.ServeHTTP"
//     },
//     {
//       "location": "goroot/src/net/http/server.go:2747",
//       "function": "net/http.(*ServeMux).ServeHTTP"
//     },
//     {
//       "location": "goroot/src/net/http/server.go:3210",
//       "function": "net/http.serverHandler.ServeHTTP"
//     },
//     {
//       "location": "goroot/src/net/http/server.go:2092",
//       "function": "net/http.(*conn).serve"
//     },
//     {
//       "location": "goroot/src/runtime/asm_arm64.s:1223",
//       "function": "runtime.goexit"
//     }
//   ],
//   "joined_errors": [
//     {
//       "caller_stack_trace": {
//         "location": "cwd/flaw/example/main.go:44",
//         "function": "main.insertRedisKey.func1"
//       },
//       "error": {
//         "message": "failed to close file: wtf",
//         "type_name": "*errors.errorString",
//         "syntax_representation": "&errors.errorString{s:\"failed to close file: wtf\"}"
//       }
//     },
//     {
//       "caller_stack_trace": {
//         "location": "cwd/flaw/example/main.go:44",
//         "function": "main.insertRedisKey.func1"
//       },
//       "error": {
//         "message": "permission denied",
//         "type_name": "*errors.errorString",
//         "syntax_representation": "&errors.errorString{s:\"permission denied\"}"
//       }
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
