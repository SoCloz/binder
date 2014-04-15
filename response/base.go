package response

import (
	"net/http"
)

type Base struct {
	StatusCode int
}

func (r *Base) ApplyTo(w http.ResponseWriter) {
	w.WriteHeader(r.StatusCode)
}

// Set the HTTP status code
func (r *Base) SetStatusCode(s int) {
	r.StatusCode = s
}
