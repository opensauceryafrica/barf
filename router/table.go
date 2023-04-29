package router

import "net/http"

var table = map[string]map[string]func(http.ResponseWriter, *http.Request){}
