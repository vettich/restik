package restik

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

type routeHandlerType int

const (
	httpHandlerType routeHandlerType = iota
	restHandlerType
	customHandlerType
)

type (
	httpHandler func(http.ResponseWriter, *http.Request)
	restHandler func(ResponseWriter, *Request)
)

type routes map[string]*Route

func (rts routes) find(method, endpoint string) (*Route, bool) {
	key := fmt.Sprintf("%s:%s", method, endpoint)
	rt, ok := rts[key]
	return rt, ok
}

// Route is route of rest
type Route struct {
	method   string
	endpoint string

	handlerType routeHandlerType
	httpHandler httpHandler
	restHandler restHandler

	args        reflect.Type
	argsIsPtr   bool
	reply       reflect.Type
	inputCount  int
	outputCount int
	fn          reflect.Value
}

// NewRoute create new route by http method and endpoint.
// fn can be one of
//
//      func([*Request | [,] *struct]) [*struct | [,] error]
//      func(http.ResponseWriter, *http.Request)
//      func(rest.ResponseWriter, *rest.Request)
func NewRoute(method, endpoint string, fn interface{}) *Route {
	if httpHndl, ok := fn.(func(http.ResponseWriter, *http.Request)); ok {
		return &Route{
			method:      method,
			endpoint:    endpoint,
			handlerType: httpHandlerType,
			httpHandler: httpHndl,
		}
	}

	if restHndl, ok := fn.(func(ResponseWriter, *Request)); ok {
		return &Route{
			method:      method,
			endpoint:    endpoint,
			handlerType: restHandlerType,
			restHandler: restHndl,
		}
	}

	fnTypeOf := reflect.TypeOf(fn)
	args, argsIsPtr := parseInput(fnTypeOf)
	return &Route{
		method:      method,
		endpoint:    endpoint,
		inputCount:  fnTypeOf.NumIn(),
		outputCount: fnTypeOf.NumOut(),
		args:        args,
		argsIsPtr:   argsIsPtr,
		reply:       parseOutput(fnTypeOf),
		fn:          reflect.ValueOf(fn),
	}
}

func parseInput(fnType reflect.Type) (reflect.Type, bool) {
	cnt := fnType.NumIn()
	if cnt == 0 {
		return nil, false
	}
	if cnt > 2 {
		panic("Function arguments count must be 2 or less")
	}
	elem, isPtr := fnType.In(cnt-1), false
	if elem.Kind() == reflect.Ptr {
		isPtr = true
		elem = elem.Elem()
	}
	if elem.String() != "rest.Request" {
		return elem, isPtr
	}
	return nil, false
}

func parseOutput(fnType reflect.Type) reflect.Type {
	cnt := fnType.NumOut()
	if cnt == 0 {
		return nil
	}
	if cnt > 2 {
		panic("Function arguments count must be 2 or less")
	}
	elem := fnType.Out(0)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	if elem.String() != "rest.Request" {
		return elem
	}
	return nil
}

// Get return new route by GET method
func Get(endpoint string, fn interface{}) *Route {
	return NewRoute("GET", endpoint, fn)
}

// Post return new route by POST method
func Post(endpoint string, fn interface{}) *Route {
	return NewRoute("POST", endpoint, fn)
}

// Delete return new route by DELETE method
func Delete(endpoint string, fn interface{}) *Route {
	return NewRoute("DELETE", endpoint, fn)
}

// Patch return new route by PATCH method
func Patch(endpoint string, fn interface{}) *Route {
	return NewRoute("PATCH", endpoint, fn)
}

// Put return new route by PUT method
func Put(endpoint string, fn interface{}) *Route {
	return NewRoute("PUT", endpoint, fn)
}

func (rt Route) getKey() string {
	return fmt.Sprintf("%s:%s", rt.method, rt.endpoint)
}

func (rt Route) exec(hr *http.Request) *serveReply {
	fnArgs, err := rt.getArgs(hr)
	if err != nil {
		return &serveReply{Error: FromAnotherError(err)}
	}
	var resValue []reflect.Value
	resValue = rt.fn.Call(fnArgs)

	var replyVal interface{}
	var errVal interface{}

	if rt.outputCount == 1 {
		if rt.reply == nil {
			errVal = resValue[0].Interface()
		} else {
			replyVal = resValue[0].Interface()
		}
	} else if rt.outputCount == 2 {
		replyVal = resValue[0].Interface()
		errVal = resValue[1].Interface()
	}

	res := serveReply{}
	if errVal != nil {
		err, _ := errVal.(error)
		res.Error = FromAnotherError(err)
	}
	if replyVal != nil {
		res.Response = replyVal
	}
	return &res
}

func (rt Route) getArgs(hr *http.Request) ([]reflect.Value, error) {
	var fnArgs []reflect.Value
	var req = reflect.ValueOf(NewRequest(hr))

	var argsVal reflect.Value
	if rt.args != nil {
		argsVal = reflect.New(rt.args)
		queries := hr.URL.Query()
		query := queries.Get("query")
		var queryRaw []byte
		if query != "" {
			unQuery, _ := url.QueryUnescape(query)
			queryRaw = []byte(unQuery)
		} else if hr.ContentLength > 0 {
			queryRaw, _ = ioutil.ReadAll(hr.Body)
		}
		if queryRaw != nil && len(queryRaw) > 0 {
			err := json.Unmarshal(queryRaw, argsVal.Interface())
			if err != nil {
				return nil, NewBadRequestError()
			}
		}

		if !rt.argsIsPtr {
			argsVal = argsVal.Elem()
		}
	}

	if rt.inputCount == 1 {
		if rt.args != nil {
			fnArgs = []reflect.Value{argsVal}
		} else {
			fnArgs = []reflect.Value{req}
		}
	} else if rt.inputCount == 2 {
		fnArgs = []reflect.Value{req, argsVal}
	}
	return fnArgs, nil
}
