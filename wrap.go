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

// A controller action wrapper
type Wrapper struct {
	Call       reflect.Value
	IsVariadic bool
	Args       []*ActionArgs
}

// Wraps a controller action
//
// Example :
//   func MyAction(id int, param string) {}
//   http.Handle("/my_action", binder.Wrap(MyAction, "id", "param")
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

// The http handler for the wrapped action
func (wr *Wrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	methodArgs := make([]reflect.Value, len(wr.Args))
	for i, a := range wr.Args {
		methodArgs[i] = GetBoundValue(r, a.Name, a.Type)
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
