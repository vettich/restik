package restik

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Vars is endpoint path variables
type Vars map[string]interface{}

// NewVars return new Vars instance
func NewVars(r *http.Request) Vars {
	vars := Vars{}
	for k, v := range mux.Vars(r) {
		vars[k] = v
	}
	return vars
}

// Set set value in key
func (vars Vars) Set(key string, value interface{}) {
	vars[key] = value
}

// Int return int value from vars by key
func (vars Vars) Int(key string) (ret int, ok bool) {
	var v interface{}
	v, ok = vars[key]
	if !ok {
		return
	}
	switch i := v.(type) {
	case int:
		ret = i
	case int32:
		ret = int(i)
	case int64:
		ret = int(i)
	case uint:
		ret = int(i)
	case uint32:
		ret = int(i)
	case uint64:
		ret = int(i)
	case float32:
		ret = int(i)
	case float64:
		ret = int(i)
	case string:
		ret, _ = strconv.Atoi(i)
	default:
		ok = false
	}
	return
}

// Int32 return int32 value from vars by key
func (vars Vars) Int32(key string) (ret int32, ok bool) {
	var v interface{}
	v, ok = vars[key]
	if !ok {
		return
	}
	switch i := v.(type) {
	case int:
		ret = int32(i)
	case int32:
		ret = i
	case int64:
		ret = int32(i)
	case uint:
		ret = int32(i)
	case uint32:
		ret = int32(i)
	case uint64:
		ret = int32(i)
	case float32:
		ret = int32(i)
	case float64:
		ret = int32(i)
	case string:
		s, _ := strconv.Atoi(i)
		ret = int32(s)
	default:
		ok = false
	}
	return
}

// Int64 return int64 value from vars by key
func (vars Vars) Int64(key string) (ret int64, ok bool) {
	var v interface{}
	v, ok = vars[key]
	if !ok {
		return
	}
	switch i := v.(type) {
	case int:
		ret = int64(i)
	case int32:
		ret = int64(i)
	case int64:
		ret = i
	case uint:
		ret = int64(i)
	case uint32:
		ret = int64(i)
	case uint64:
		ret = int64(i)
	case float32:
		ret = int64(i)
	case float64:
		ret = int64(i)
	case string:
		s, _ := strconv.Atoi(i)
		ret = int64(s)
	default:
		ok = false
	}
	return
}

// Uint return uint value from vars by key
func (vars Vars) Uint(key string) (ret uint, ok bool) {
	var v interface{}
	v, ok = vars[key]
	if !ok {
		return
	}
	switch i := v.(type) {
	case int:
		ret = uint(i)
	case int32:
		ret = uint(i)
	case int64:
		ret = uint(i)
	case uint:
		ret = uint(i)
	case uint32:
		ret = uint(i)
	case uint64:
		ret = uint(i)
	case float32:
		ret = uint(i)
	case float64:
		ret = uint(i)
	case string:
		s, _ := strconv.Atoi(i)
		ret = uint(s)
	default:
		ok = false
	}
	return
}

// Uint32 return uint32 value from vars by key
func (vars Vars) Uint32(key string) (ret uint32, ok bool) {
	var v interface{}
	v, ok = vars[key]
	if !ok {
		return
	}
	switch i := v.(type) {
	case int:
		ret = uint32(i)
	case int32:
		ret = uint32(i)
	case int64:
		ret = uint32(i)
	case uint:
		ret = uint32(i)
	case uint32:
		ret = uint32(i)
	case uint64:
		ret = uint32(i)
	case float32:
		ret = uint32(i)
	case float64:
		ret = uint32(i)
	case string:
		s, _ := strconv.Atoi(i)
		ret = uint32(s)
	default:
		ok = false
	}
	return
}

// Uint64 return uint64 value from vars by key
func (vars Vars) Uint64(key string) (ret uint64, ok bool) {
	var v interface{}
	v, ok = vars[key]
	if !ok {
		return
	}
	switch i := v.(type) {
	case int:
		ret = uint64(i)
	case int32:
		ret = uint64(i)
	case int64:
		ret = uint64(i)
	case uint:
		ret = uint64(i)
	case uint32:
		ret = uint64(i)
	case uint64:
		ret = uint64(i)
	case float32:
		ret = uint64(i)
	case float64:
		ret = uint64(i)
	case string:
		s, _ := strconv.Atoi(i)
		ret = uint64(s)
	default:
		ok = false
	}
	return
}

// String return string value from vars by key
func (v Vars) String(key string) (string, bool) {
	val, ok := v[key]
	if !ok {
		return "", false
	}
	return fmt.Sprint(val), true
}
