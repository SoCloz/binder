package response

import (
	"encoding/json"
	"net/http"
)

// A JSON Response
type Json struct {
	// Payload
	Data interface{}
	Base
}

func (r *Json) ApplyTo(w http.ResponseWriter) {
	jsonResult, err := json.MarshalIndent(r.Data, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	r.Base.ApplyTo(w)
	w.Write(jsonResult)
}
