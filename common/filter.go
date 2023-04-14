package common

import "net/http"

type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

type Filter struct {
	// store URI to be filtered
	filterMap map[string]FilterHandle
}

func NewFilter() *Filter {
	return &Filter{
		filterMap: make(map[string]FilterHandle),
	}
}

func (f *Filter) RegisterFilterUri(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

func (f *Filter) GetFilterHandle(uri string) FilterHandle {
	return f.filterMap[uri]
}

type WebHandle func(rw http.ResponseWriter, req *http.Request)

func (f *Filter) Handle(webHandle WebHandle) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		for path, handle := range f.filterMap {
			if path == r.RequestURI {
				// filter
				err := handle(rw, r)
				if err != nil {
					rw.Write([]byte(err.Error()))
					return
				}
				break
			}
		}
		// Execute register function
		webHandle(rw, r)
	}
}
