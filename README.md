binder
======

A thin wrapper around http handlers to ease handler coding.

Works with Gorilla Mux or without.

Without binder :

```go
func ViewItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
    // Fetch item from DB
	item, err := ...
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonResult, err := json.Marshal(item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}

r.HandleFunc("/items/{id}", ViewItemHandler)
```

With binder :

```go
func ViewItemAction(i *Item) response.Response {
	if i == nil {
		return &response.Error{http.StatusNotFound, "not found, sorry"}
	}
	return &response.Json{i}
}
r.Handle("/items/{id}", binder.NewActionWrapper(ViewItemAction, "id"))
```

