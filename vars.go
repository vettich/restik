package restik

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Vars is endpoint path variables
type Vars map[string]string

// NewVars return new Vars instance
func NewVars(r *http.Request) Vars {
	return mux.Vars(r)
}

// Int return int value from vars by key
func (v Vars) Int(key string) int {
	val, ok := v[key]
	if !ok {
		return 0
	}
	i, _ := strconv.Atoi(val)
	return i
}

// String return string value from vars by key
func (v Vars) String(key string) string {
	val, ok := v[key]
	if !ok {
		return ""
	}
	return val
}
