package myapp

import (
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter
}
