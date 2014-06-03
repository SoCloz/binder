package response

import (
	"net/http"
)

type Base struct {
	StatusCode int
}

func (r *Base) ApplyTo(w http.ResponseWriter) {
	if r.StatusCode != 0 {
		w.WriteHeader(r.StatusCode)
	}
}

// Set the HTTP status code
func (r *Base) SetStatusCode(s int) {
	r.StatusCode = s
}
