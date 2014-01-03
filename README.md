binder
======

http handlers without the http.Request/http.ResponseWriter clutter

Uses Gorilla Mux.

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

Bindings
--------

Binder currently binds :
* integers (signed/unsigned)
* floats
* strings
* booleans (true/false, on/off, 1/0)
* slices (param=value1&param=value2&param=value3 or param=value1,value2,value3)
* pointers (if no value found, binds nil)
* custom binders

Custom binders
--------------

You can bind an "id" url parameter to a database record using a custom binder.

```go

var (
	ItemBinder = func(val string, typ reflect.Type) (reflect.Value, bool) {
		// loads the object from the DB
		i, err := [...]
		if err != nil {
			return reflect.Zero(typ), true
		}
		return reflect.ValueOf(i), false
	}

)

func init() {
	binder.RegisterBinder(new(Item), ItemBinder)
}
```

Responses
---------

The following responses are possible :
* basic
```go
return &response.Basic{"content"}
```
* json
```go
return &response.Json{data}
```
* error
```go
return &response.Error{http.StatusNotFound, "content"}
```

Roadmap
-------

* support non gorilla mux http handlers
* struct support
* more response types