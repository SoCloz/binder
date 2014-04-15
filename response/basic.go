package response

import (
	"net/http"
)

// An basic http response
type Basic struct {
	// Content
	Content string
	Base
}

func (r *Basic) ApplyTo(w http.ResponseWriter) {
	r.Base.ApplyTo(w)
	w.Write([]byte(r.Content))
}
