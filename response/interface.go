package response

import (
	"net/http"
)

// Interface for HTTP responses
type Response interface {
	SetStatusCode(int)
	ApplyTo(http.ResponseWriter)
}
