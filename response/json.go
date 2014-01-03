package response

import (
	"encoding/json"
	"net/http"
)

// A JSON Response
type Json struct {
	Data interface{}
}

func (r *Json) ApplyTo(w http.ResponseWriter) {
	jsonResult, err := json.Marshal(r.Data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsonResult)
}
