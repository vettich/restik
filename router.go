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
		r.muxRouter.Handle(rt.endpoint, r).Methods(rt.method)
	}
	return r
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
