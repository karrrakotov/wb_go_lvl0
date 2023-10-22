package transport

import "net/http"

type HandlerOrder interface {
	Init(router *http.ServeMux)
}
