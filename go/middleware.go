package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
)

type ModifierMiddleware struct {
	handler http.Handler
}

func (m *ModifierMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rec := httptest.NewRecorder()

	// passing a ResponseRecorder instead of original RW
	m.handler.ServeHTTP(rec, r)

	// after this finishe, we have the response recorded
	// and can modify it before copying it to the original RW

	//we copy the origina headers first
	for k, v := range rec.Header() {
		w.Header()[k] = v
	}

	// and set an additional one
	w.Header().Set("X-My-Own-Header", "Yo !")

	// only then the status code, as this call writes the header as well
	w.WriteHeader(418)

	// The body hasn't been written (to the real RW) yet,
	// so we can prepend some data

	data := []byte("Middleware says hello again\n")

	// But the Content-Length might have been set already
	// we should modify it by adding length
	// of our own data
	// Ignoring the error is fine here:
	// if Content-Lenght is empty or otherwise invalid,
	// Atoi() will return zero,
	// which is just what we'd want in that case

	clen, _ := strconv.Atoi(r.Header.Get("Content-Length"))
	clen += len(data)

	w.Header().Set("Content-Length", strconv.Itoa(clen))

	//finally, write out our data

	w.Write(data)

	// then write out the original body
	w.Write(rec.Body.Bytes())
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Success!\n"))
}
func main() {
	mid := &ModifierMiddleware{http.HandlerFunc(myHandler)}
	println("Listening on port 8080")
	http.ListenAndServe(":8080", mid)
}
