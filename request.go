package restik

import "net/http"

// Request represent arg for handle func
type Request struct {
	Vars    Vars
	Headers http.Header
	Request *http.Request
}

// NewRequest create new Request instance
func NewRequest(r *http.Request) *Request {
	return &Request{
		NewVars(r),
		r.Header,
		r,
	}
}
