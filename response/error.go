package response

import (
	"net/http"
)

// An Error response
type Error struct {
	// HTTP Status code
	StatusCode int
	// Message
	Message string
}

func (r *Error) ApplyTo(w http.ResponseWriter) {
	w.WriteHeader(r.StatusCode)
	w.Write([]byte(r.Message))
}
