# restik

`restik` a Go library that provides an easy and convenient way to implement a REST service.

# Import

`import "github.com/vettich/restik`

# Usage

## Echo example

```go
func echo(msg string) string {
  return msg
}

func main() {
  r := restik.NewRouter()
  r.Get("/echo", echo)
  http.Handle("/", r.Handler())
  http.ListenAndServe("0.0.0.0:8000", nil)
}
```

## Hello world example

```go
type helloArg struct {
  Value string `json:"value"`
}

func hello(arg helloArg) (string, error) {
  if arg.Value == "" {
    return "", errors.New("value is empty")
  }
  return fmt.Sprint("Hello," arg.Value)
}

func main() {
  r := restik.NewRouter()
  r.Post("/hello", hello)
  http.Handle("/", r.Handler())
  http.ListenAndServe("0.0.0.0:8000", nil)
}
```
