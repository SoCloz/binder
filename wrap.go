package binder

import (
	"net/http"
	"reflect"

	"github.com/SoCloz/binder/response"
)

// Action Parameters
type ActionArgs struct {
	Name string
	Type reflect.Type
}

// A handler wrapping an action
type Wrapper struct {
	Call       reflect.Value
	IsVariadic bool
	Args       []*ActionArgs
}

// Creates an handler wrapping an action
// func MyAction(id int, param string) {}
// Example: binder.NewActionHandler(MyAction, "id", "param")
func Wrap(call interface{}, params ...string) *Wrapper {
	w := new(Wrapper)
	w.Call = reflect.ValueOf(call)
	callType := w.Call.Type()
	w.IsVariadic = callType.IsVariadic()
	w.Args = make([]*ActionArgs, callType.NumIn())

	if callType.NumIn() < len(params) {
		panic("Wrong number of params")
	}
	for i := 0; i < callType.NumIn(); i++ {
		typ := callType.In(i)
		var paramName string
		 if i < len(params) {
		 	paramName = params[i]
		 } else {
		 	paramName = "*"
		 }
		w.Args[i] = &ActionArgs{Name: paramName, Type: typ}
		if paramName == "*" {
			RegisterStructAttributes(typ)
		}
	}
	return w
}

// The http handler
func (wr *Wrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	methodArgs := make([]reflect.Value, len(wr.Args))
	for i, a := range wr.Args {
		methodArgs[i] = GetBoundValue(values, a.Name, a.Type)
	}
	var resultValue reflect.Value
	if wr.IsVariadic {
		resultValue = wr.Call.CallSlice(methodArgs)[0]
	} else {
		resultValue = wr.Call.Call(methodArgs)[0]
	}
	resp := resultValue.Interface().(response.Response)
	resp.ApplyTo(w)
}