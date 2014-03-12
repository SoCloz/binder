binder
======

Go web micro-framework - http handlers without the http.Request/http.ResponseWriter clutter

binder is currently used in production.

Documentation : http://godoc.org/github.com/SoCloz/binder

Works with pat or using standard net/http handler.

Without binder :

```go
func ViewItemHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
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

http.Handle("/items", ViewItemHandler)
```

With binder :

```go
func ViewItem(i *Item) response.Response {
	if i == nil {
		return &response.Error{http.StatusNotFound, "not found, sorry"}
	}
	return &response.Json{i}
}

http.Handle("/items", binder.Wrap(ViewItem, "id"))
```

binder is compatible with pat (no need to add ":" to your bindings, binder automatically adds it) :

```go
m.Get("/items/:id", binder.Wrap(ViewItem, "id"))
```

Bindings
--------

binder binds query string or url parameters to your controller parameters.

```go
func MyAction(param1 string, param2 float32, param3 []string, param4 *int) response.Response {
	str := fmt.Sprintf("param1=%v, param2=%v, param3=%v, param4=%v", param1, param2, param3, param4)
	return &response.Basic{str}
}

http.Handle("/my_action", binder.Wrap(MyAction, "param1", "param2", "param3", "param4"))
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
* structs
* custom binders
* the request

Custom binders
--------------

You can bind an "id" url parameter to a database record using a custom binder.

```go

import(
	"github.com/SoCloz/binder"
)

var (
	ItemBinder = func(values url.Values, name string, typ reflect.Type) (reflect.Value, bool) {
		id := binder.GetValue(values, name)
		// loads the object of id "id" from the DB
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

And then :

```go
func ViewItem(i *Item) response.Response {}

m.Get("/items/:id", binder.Wrap(ViewItem, "id"))
```

Binding the request
-------------------

If your controller has a parameter of type *http.Request, it will be bound to the current http request. You can wrap it using "*" (or omit "*" if it is at the end of the call) :

```go
func ViewItem(id int, r *http.Request) response.Response {}

m.Get("/items/:id", binder.Wrap(ViewItem, "id", "*"))
m.Get("/items/:id", binder.Wrap(ViewItem, "id"))
```

Binding structs
---------------

Struct bindings are defined using field tags :

```go
type Options struct {
	OnlyNames string   `binder:"only_names"`
	Page      int      `binder:"page"`
	Tags      []string `binder:"tags"`
}

func ViewItem(id int, opt Options) response.Response {}
```

And wrapped using "*" :

```go
m.Get("/items/:id", binder.Wrap(ViewItem, "id", "*"))
```

You can omit all "*" at the end of your Wrap call :

```go
m.Get("/items/:id", binder.Wrap(ViewItem, "id"))
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

* binder : bind http headers
* more response types

License
-------

See LICENCE

Thanks
------

Binding code was largely borrowed from revel - https://github.com/robfig/revel