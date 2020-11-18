package restik

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router is main router in rest
type Router struct {
	routes    routes
	muxRouter *mux.Router
}

// NewRouter create new Router
func NewRouter() *Router {
	return &Router{
		routes:    routes{},
		muxRouter: mux.NewRouter(),
	}
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

func (r *Router) ServeHTTP(hw http.ResponseWriter, hr *http.Request) {
	rt, ok := r.getCurrentRoute(hr)
	if !ok {
		writeError(hw, ErrNotFoundEndpoint)
		return
	}

	if rt.handlerType == httpHandlerType {
		rt.httpHandler(hw, hr)
		return
	}

	if rt.handlerType == restHandlerType {
		rt.restHandler(NewResponseWriter(hw), NewRequest(hr))
		return
	}

	sr := rt.exec(hr)
	writeReply(hw, sr)
}

func (r *Router) getCurrentRoute(hr *http.Request) (*Route, bool) {
	muxRouter := mux.CurrentRoute(hr)
	pathTemplate, err := muxRouter.GetPathTemplate()
	if err != nil {
		return nil, false
	}
	return r.routes.find(hr.Method, pathTemplate)
}
