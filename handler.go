package binder

import(
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"

	"github.com/SoCloz/binder/response"
)

// Action Parameters
type ActionArgs struct {
	Name string
	Type reflect.Type
}

// A handler wrapping an action
type Handler struct {
	Call reflect.Value
	IsVariadic bool
	Args []*ActionArgs
}

// Creates an handler wrapping an action
// func MyAction(id int, param string) {}
// Example: binder.NewActionHandler(MyAction, "id", "param")
func NewActionHandler(call interface{}, params ...string) *Handler {
	w := new(Wrapper)
	w.Call = reflect.ValueOf(call)
	callType := w.Call.Type()
	w.IsVariadic = callType.IsVariadic()
	w.Args = make([]*ActionArgs, callType.NumIn())

	if callType.NumIn() != len(params) {
		panic("Wrong number of params")
	}
	for i := 0; i < callType.NumIn(); i++ {
		typ := callType.In(i)
		w.Args[i] = &ActionArgs{Name: params[i], Type: typ}
		log.Printf("%+v %+v", typ, w.Args[i])
	}
	log.Printf("%+v", w)
	return w
}

// The handler
func (wr *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	values := r.URL.Query()

	methodArgs := make([]reflect.Value, len(wr.Args))
	for i, a := range wr.Args {
		if a.Name == "*" {
			methodArgs[i] = reflect.ValueOf(values)
		} else {
			if _, found := vars[a.Name]; found {
				methodArgs[i] = BindFromMap(vars, a.Name, a.Type)
			} else {
				methodArgs[i] = BindFromValues(values, a.Name, a.Type)
			}
		}
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