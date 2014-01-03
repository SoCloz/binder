package response

import (
	"net/http"
)

// Interface for HTTP responses
type Response interface {
	ApplyTo(http.ResponseWriter)
}
