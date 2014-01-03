package response

import (
	"net/http"
)

// An basic http response
type Basic struct {
	// Content
	Content string
}

func (r *Basic) ApplyTo(w http.ResponseWriter) {
	w.Write([]byte(r.Content))
}
