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
	Base
}

func (r *Error) ApplyTo(w http.ResponseWriter) {
	r.SetStatusCode(r.StatusCode)
	r.Base.ApplyTo(w)
	w.Write([]byte(r.Message))
}
