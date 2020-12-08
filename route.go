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
	customHandlerType routeHandlerType = iota
	httpHandlerType
	restHandlerType
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
	Method   string
	Endpoint string

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
			Method:      method,
			Endpoint:    endpoint,
			handlerType: httpHandlerType,
			httpHandler: httpHndl,
		}
	}

	if restHndl, ok := fn.(func(ResponseWriter, *Request)); ok {
		return &Route{
			Method:      method,
			Endpoint:    endpoint,
			handlerType: restHandlerType,
			restHandler: restHndl,
		}
	}

	fnTypeOf := reflect.TypeOf(fn)
	args, argsIsPtr := parseInput(fnTypeOf)
	return &Route{
		Method:      method,
		Endpoint:    endpoint,
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
	if elem.String() != "restik.Request" {
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
	return elem
}

func (rt Route) getKey() string {
	return fmt.Sprintf("%s:%s", rt.Method, rt.Endpoint)
}

func (rt Route) exec(rr *Request, rpl Reply) {
	fnArgs, err := rt.getArgs(rr)
	if err != nil {
		rpl.SetError(err)
		return
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

	if errVal != nil {
		err, _ := errVal.(error)
		rpl.SetError(err)
	}
	if replyVal != nil {
		rpl.SetResponse(replyVal)
	}
	return
}

func (rt *Route) getArgs(rr *Request) ([]reflect.Value, error) {
	var fnArgs []reflect.Value
	var req = reflect.ValueOf(rr)

	var argsVal reflect.Value
	if rt.args != nil {
		argsVal = reflect.New(rt.args)
		queries := rr.URL.Query()
		query := queries.Get("query")
		var queryRaw []byte
		if query != "" {
			unQuery, _ := url.QueryUnescape(query)
			queryRaw = []byte(unQuery)
		} else if rr.ContentLength > 0 {
			queryRaw, _ = ioutil.ReadAll(rr.Body)
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
