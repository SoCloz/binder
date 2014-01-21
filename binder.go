package binder

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// A Binder translates between string parameters and Go data structures.
type ValueBinder func(r *http.Request, name string, typ reflect.Type) (reflect.Value, bool)

// The map of query string param name => struct attribute name
type StructAttributes map[string]string

var (
	// Lookup tables
	TypeBinders         = make(map[reflect.Type]ValueBinder)
	KindBinders         = make(map[reflect.Kind]ValueBinder)
	StructAttributesMap = make(map[reflect.Type]StructAttributes)

	// Binds an integer (signed)
	IntBinder = func(r *http.Request, name string, typ reflect.Type) (reflect.Value, bool) {
		val := GetValue(r, name)
		if len(val) == 0 {
			return reflect.Zero(typ), true
		}
		intValue, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return reflect.Zero(typ), true
		}
		pValue := reflect.New(typ)
		pValue.Elem().SetInt(intValue)
		return pValue.Elem(), false
	}

	// Binds an integer (unsigned)
	UintBinder = func(r *http.Request, name string, typ reflect.Type) (reflect.Value, bool) {
		val := GetValue(r, name)
		if len(val) == 0 {
			return reflect.Zero(typ), true
		}
		uintValue, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return reflect.Zero(typ), true
		}
		pValue := reflect.New(typ)
		pValue.Elem().SetUint(uintValue)
		return pValue.Elem(), false
	}

	// Binds a float
	FloatBinder = func(r *http.Request, name string, typ reflect.Type) (reflect.Value, bool) {
		val := GetValue(r, name)
		if len(val) == 0 {
			return reflect.Zero(typ), true
		}
		floatValue, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return reflect.Zero(typ), true
		}
		pValue := reflect.New(typ)
		pValue.Elem().SetFloat(floatValue)
		return pValue.Elem(), false
	}

	// Binds a string
	StringBinder = func(r *http.Request, name string, typ reflect.Type) (reflect.Value, bool) {
		val := GetValue(r, name)
		return reflect.ValueOf(val), false
	}

	// Binds a boolean value, using the following formats :
	// "true" and "false"
	// "on" and "" (a checkbox)
	// "1" and "0" (why not)
	BoolBinder = func(r *http.Request, name string, typ reflect.Type) (reflect.Value, bool) {
		val := GetValue(r, name)
		v := strings.TrimSpace(strings.ToLower(val))
		switch v {
		case "yes", "true", "on", "1":
			return reflect.ValueOf(true), false
		}
		// Return false by default.
		return reflect.ValueOf(false), false
	}

	// Binds a comma separated list to a slice
	StringSliceBinder = func(r *http.Request, name string, typ reflect.Type) (reflect.Value, bool) {
		val := GetValue(r, name)
		if val == "" {
			return reflect.MakeSlice(typ, 0, 0), true
		}
		split := strings.Split(val, ",")
		result := reflect.MakeSlice(typ, len(split), len(split))
		old := r.URL.RawQuery
		for index, item := range split {
			r.URL.RawQuery = fmt.Sprintf("%s=%s", url.QueryEscape(name), url.QueryEscape(item))
			itemVal, _ := Bind(r, name, typ.Elem())
			result.Index(index).Set(itemVal)
		}
		r.URL.RawQuery = old
		return result, false
	}

	// Binds a set of parameters to a struct.
	StructBinder = func(r *http.Request, name string, typ reflect.Type) (reflect.Value, bool) {
		result := reflect.New(typ).Elem()
		for fieldName, attrName := range StructAttributesMap[typ] {
			fieldValue := result.FieldByName(fieldName)
			boundVal, _ := Bind(r, attrName, fieldValue.Type())
			fieldValue.Set(boundVal)
		}
		return result, false
	}

	// Binds a pointer. If nothing found, returns nil.
	PointerBinder = func(r *http.Request, name string, typ reflect.Type) (reflect.Value, bool) {
		pointerOf := typ.Elem()
		value, isNil := Bind(r, name, pointerOf)
		if isNil {
			return reflect.Zero(typ), true
		} else {
			return value.Addr(), false
		}

		return reflect.Zero(typ), true
	}

	// Binds the http request.
	RequestBinder = func(r *http.Request, name string, typ reflect.Type) (reflect.Value, bool) {
		return reflect.ValueOf(r), false
	}
)

// Builds the lookup table
func init() {
	KindBinders[reflect.Int] = IntBinder
	KindBinders[reflect.Int8] = IntBinder
	KindBinders[reflect.Int16] = IntBinder
	KindBinders[reflect.Int32] = IntBinder
	KindBinders[reflect.Int64] = IntBinder

	KindBinders[reflect.Uint] = UintBinder
	KindBinders[reflect.Uint8] = UintBinder
	KindBinders[reflect.Uint16] = UintBinder
	KindBinders[reflect.Uint32] = UintBinder
	KindBinders[reflect.Uint64] = UintBinder

	KindBinders[reflect.Float32] = FloatBinder
	KindBinders[reflect.Float64] = FloatBinder

	KindBinders[reflect.String] = StringBinder
	KindBinders[reflect.Bool] = BoolBinder
	KindBinders[reflect.Slice] = StringSliceBinder
	KindBinders[reflect.Ptr] = PointerBinder
	KindBinders[reflect.Struct] = StructBinder

	TypeBinders[reflect.TypeOf(&http.Request{})] = RequestBinder
}

// Takes the name and type of the desired parameter and constructs it
// from query string values.
// Returns the zero value of the type upon any sort of failure.
func GetBoundValue(r *http.Request, name string, typ reflect.Type) reflect.Value {
	ret, _ := Bind(r, name, typ)
	return ret
}

// Constructs a value of a specific type
func Bind(r *http.Request, name string, typ reflect.Type) (reflect.Value, bool) {
	if binder, found := TypeBinders[typ]; found {
		return binder(r, name, typ)
	} else {
		if binder, found := KindBinders[typ.Kind()]; found {
			return binder(r, name, typ)
		}
	}
	return reflect.Zero(typ), true
}

// Registers a custom binder for a specific type
func RegisterBinder(i interface{}, binder ValueBinder) {
	typ := reflect.ValueOf(i).Type()
	TypeBinders[typ] = binder
}

// Register the mapping of params to struct fields
func RegisterStructAttributes(typ reflect.Type) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return
	}
	StructAttributesMap[typ] = make(StructAttributes)
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Type.Kind() == reflect.Struct {
			RegisterStructAttributes(f.Type)
			StructAttributesMap[typ][f.Name] = "*"
		} else {
			tag := f.Tag.Get("binder")
			if tag != "" {
				StructAttributesMap[typ][f.Name] = tag
			}
		}
	}
}

// Get a single value from the request query (plain or pat url parameter)
func GetValue(r *http.Request, name string) string {
	values := r.URL.Query()
	// pat url parameter
	if _, found := values[":"+name]; found {
		return values.Get(":" + name)
	}
	// query string
	return values.Get(name)
}
