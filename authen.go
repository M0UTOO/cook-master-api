package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/foo", foo)
	http.HandleFunc("/bar", bar)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="My Realm"`)
	w.WriteHeader(401)
	w.Write([]byte("401 Unauthorized"))
}

func foo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("foo"))
}

func bar(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bar"))
}
