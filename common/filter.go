package common

import (
	"net/http"
	"strings"
)

type FilterHandler func(resp http.ResponseWriter, req *http.Request) error
type WebHandler func(resp http.ResponseWriter, req *http.Request)

type Filter struct {
	filterMap map[string]FilterHandler
}

func NewFilter() *Filter {
	return &Filter{filterMap: make(map[string]FilterHandler)}
}

func (filter *Filter) RegisterFilterUri(uri string, handle FilterHandler) {
	filter.filterMap[uri] = handle
}

func (filter *Filter) GetFilterHandler(uri string) FilterHandler {
	return filter.filterMap[uri]
}

func (filter *Filter) Handler(w WebHandler) func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		for uri, f := range filter.filterMap {
			if strings.Contains(req.RequestURI, uri) {
				err := f(resp, req)
				if err != nil {
					resp.Write([]byte(err.Error()))
					return
				}
				break
			}
		}
		w(resp, req)
	}
}
