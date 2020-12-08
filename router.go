package restik

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router is main router in rest
type Router struct {
	routes                  routes
	muxRouter               *mux.Router
	middlewares             []Middleware
	notFoundHandler         func(ResponseWriter, *Request)
	methodNotAllowedHandler func(ResponseWriter, *Request)
	replyImpl               Reply
}

// NewRouter create new Router
func NewRouter() *Router {
	r := &Router{
		routes:      routes{},
		muxRouter:   mux.NewRouter(),
		middlewares: make([]Middleware, 0),
		replyImpl:   &serveReply{},
	}
	r.muxRouter.Handle("/", r)
	r.muxRouter.MethodNotAllowedHandler = methodNotAllowedHandler{r}
	r.muxRouter.NotFoundHandler = notFoundHandler{r}
	return r
}

// Handler return http handler
func (r *Router) Handler() *mux.Router {
	return r.muxRouter
}

// Add add new routers
func (r *Router) Add(rts ...*Route) *Router {
	for _, rt := range rts {
		r.routes[rt.getKey()] = rt
		r.muxRouter.Handle(rt.Endpoint, r).Methods(rt.Method)
	}
	return r
}

// Get add new route with GET method to router and return
func (r *Router) Get(endpoint string, fn interface{}) *Route {
	rt := NewRoute("GET", endpoint, fn)
	r.Add(rt)
	return rt
}

// Post add new route with POST method to router and return
func (r *Router) Post(endpoint string, fn interface{}) *Route {
	rt := NewRoute("POST", endpoint, fn)
	r.Add(rt)
	return rt
}

// Delete add new route with DELETE method to router and return
func (r *Router) Delete(endpoint string, fn interface{}) *Route {
	rt := NewRoute("DELETE", endpoint, fn)
	r.Add(rt)
	return rt
}

// Patch add new route with PATCH method to router and return
func (r *Router) Patch(endpoint string, fn interface{}) *Route {
	rt := NewRoute("PATCH", endpoint, fn)
	r.Add(rt)
	return rt
}

// Put add new route with PUT method to router and return
func (r *Router) Put(endpoint string, fn interface{}) *Route {
	rt := NewRoute("PUT", endpoint, fn)
	r.Add(rt)
	return rt
}

// SetNotFoundHandlerFunc set not found handler
func (r *Router) SetNotFoundHandlerFunc(handler func(ResponseWriter, *Request)) {
	r.notFoundHandler = handler
}

// SetMethodNotAllowedHandlerFunc set not found handler
func (r *Router) SetMethodNotAllowedHandlerFunc(handler func(ResponseWriter, *Request)) {
	r.methodNotAllowedHandler = handler
}

func (r *Router) Use(middlewares ...Middleware) *Router {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

func (r *Router) ServeHTTP(hw http.ResponseWriter, hr *http.Request) {
	handle := r.routeHandler
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handle = r.middlewares[i].Middleware(handle)
	}
	rt, _ := r.getCurrentRoute(hr)
	handle(NewResponseWriter(hw, r.replyImpl), NewRequest(hr, rt))
}

func (r *Router) routeHandler(rw ResponseWriter, rr *Request) {
	rt := rr.Route
	if rt == nil {
		if r.notFoundHandler != nil {
			r.notFoundHandler(rw, rr)
		} else {
			rw.WriteError(ErrNotFoundEndpoint)
		}
		return
	}

	if rt.handlerType == httpHandlerType {
		rt.httpHandler(rw.ResponseWriter, rr.Request)
		return
	}

	if rt.handlerType == restHandlerType {
		rt.restHandler(rw, rr)
		return
	}

	rpl := r.replyImpl.New()
	rt.exec(rr, rpl)
	rw.WriteReply(rpl)
}

func (r *Router) getCurrentRoute(hr *http.Request) (*Route, bool) {
	muxRoute := mux.CurrentRoute(hr)
	if muxRoute == nil {
		return nil, false
	}
	pathTemplate, err := muxRoute.GetPathTemplate()
	if err != nil {
		return nil, false
	}
	return r.routes.find(hr.Method, pathTemplate)
}

func (r *Router) SetCustomReply(rpl Reply) {
	r.replyImpl = rpl
}

type notFoundHandler struct {
	r *Router
}

func (h notFoundHandler) ServeHTTP(hw http.ResponseWriter, hr *http.Request) {
	if h.r.notFoundHandler != nil {
		h.r.notFoundHandler(NewResponseWriter(hw, h.r.replyImpl), NewRequest(hr, nil))
		return
	}
	rw := NewResponseWriter(hw, h.r.replyImpl)
	rw.WriteError(NewNotFoundError())
}

type methodNotAllowedHandler struct {
	r *Router
}

func (h methodNotAllowedHandler) ServeHTTP(hw http.ResponseWriter, hr *http.Request) {
	if h.r.methodNotAllowedHandler != nil {
		h.r.methodNotAllowedHandler(NewResponseWriter(hw, h.r.replyImpl), NewRequest(hr, nil))
		return
	}
	rw := NewResponseWriter(hw, h.r.replyImpl)
	rw.WriteError(NewError(http.StatusMethodNotAllowed, "not_allowed", "Method not allowed"))
}
