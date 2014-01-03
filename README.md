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

r.Handle("/items/{id}", binder.NewActionHandler(ViewItemAction, "id"))
```

Bindings
--------

binder binds query string or url parameters to func parameters.

```go
func MyAction(param1 string, param2 float32, param3 []string, param4 *int) response.Response {
	log.Printf("param1=%v, param2=%v, param3=%v, param4=%v", param1, param2, param3, param4)
}

r.Handle("/my_action", binder.NewActionHandler(MyAction, "param1", "param2", "param3", "param4"))
```

```
GET /my_action
param1=, param2=0, param3=[], param4=<nil>

GET /my_action?param1=foo&param2=12.5&param3=foo,bar,baz&param4=42
param1=foo, param2=12.5, param3=[foo,bar,baz], param4=0xc210148940
```

Binder currently binds :
* integers (signed/unsigned)
* floats
* strings
* booleans (true/false, on/off, 1/0)
* slices (comma separated lists : value,value,value)
* pointers (nil if no value found)
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