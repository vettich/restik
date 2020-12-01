# restik [![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/vettich/restik) [![license](https://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/vettich/restik/master/MIT-LICENSE)

<p align="center"><img src="docs/images/logo.png" width="300"></p>

`restik` a Go library that provides an easy and convenient way to implement a REST service.

# Import

```go
import "github.com/vettich/restik"
```

# Examples

## Echo example

Simple GET request resolver

```go
func echo(req *restik.Request) string {
	return req.Vars.String("msg")
}

func main() {
  r := restik.NewRouter()
  r.Get("/echo/{msg}", echo)
  http.Handle("/", r.Handler())
  http.ListenAndServe("0.0.0.0:8000", nil)
}
```

Request:
```bash
$ curl "localhost:8000/echo/example"
```

Response:
```json
{"response":"example"}
```

## Hello world example

POST request

```go
type helloArg struct {
  Value string `json:"value"`
}

func hello(arg helloArg) (string, error) {
  if arg.Value == "" {
    return "", errors.New("value is empty")
  }
  return fmt.Sprint("Hello, " arg.Value)
}

func main() {
  r := restik.NewRouter()
  r.Post("/hello", hello)
  http.Handle("/", r.Handler())
  http.ListenAndServe("0.0.0.0:8000", nil)
}
```

Request:
```bash
$ curl -XPOST -d '{"value":"world!"}' "localhost:8000/hello"
```

Response:
```json
{"response":"Hello, world!"}
```
