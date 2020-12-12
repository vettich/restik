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
func (vars Vars) Int(key string, defaultValue ...int) int {
	i, ok := vars.Int64Ok(key)
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return int(i)
}

// IntOk return int value from vars by key
func (vars Vars) IntOk(key string) (ret int, ok bool) {
	i, ok := vars.Int64Ok(key)
	return int(i), ok
}

// Int32 return int32 value from vars by key
func (vars Vars) Int32(key string, defaultValue ...int32) int32 {
	i, ok := vars.Int64Ok(key)
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return int32(i)
}

// Int32Ok return int32 value from vars by key
func (vars Vars) Int32Ok(key string) (ret int32, ok bool) {
	i, ok := vars.Int64Ok(key)
	return int32(i), ok
}

// Int64 return int64 value from vars by key
func (vars Vars) Int64(key string, defaultValue ...int64) int64 {
	i, ok := vars.Int64Ok(key)
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return i
}

// Int64Ok return int64 value from vars by key
func (vars Vars) Int64Ok(key string) (ret int64, ok bool) {
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
func (vars Vars) Uint(key string, defaultValue ...uint) uint {
	i, ok := vars.Uint64Ok(key)
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return uint(i)
}

// UintOk return uint value from vars by key
func (vars Vars) UintOk(key string) (ret uint, ok bool) {
	i, ok := vars.Uint64Ok(key)
	return uint(i), ok
}

// Uint32 return uint32 value from vars by key
func (vars Vars) Uint32(key string, defaultValue ...uint32) uint32 {
	i, ok := vars.Uint64Ok(key)
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return uint32(i)
}

// Uint32Ok return uint32 value from vars by key
func (vars Vars) Uint32Ok(key string) (ret uint32, ok bool) {
	i, ok := vars.Uint64Ok(key)
	return uint32(i), ok
}

// Uint64 return uint64 value from vars by key
func (vars Vars) Uint64(key string, defaultValue ...uint64) uint64 {
	i, ok := vars.Uint64Ok(key)
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return i
}

// Uint64Ok return uint64 value from vars by key
func (vars Vars) Uint64Ok(key string) (ret uint64, ok bool) {
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
func (v Vars) String(key string, defaultValue ...string) string {
	s, ok := v.StringOk(key)
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return s
}

// StringOk return string value from vars by key
func (v Vars) StringOk(key string) (string, bool) {
	val, ok := v[key]
	if !ok {
		return "", false
	}
	return fmt.Sprint(val), true
}
