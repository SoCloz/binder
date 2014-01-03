package binder

import (
	"log"
	"reflect"
	"strconv"
	"strings"
	"net/url"
)

// A Binder translates between string parameters and Go data structures.
type ValueBinder func(val string, typ reflect.Type) (reflect.Value, bool)

var (
	// Lookup tables
	TypeBinders = make(map[reflect.Type]ValueBinder)
	KindBinders = make(map[reflect.Kind]ValueBinder)

	// Binds an integer (signed)
	IntBinder = func(val string, typ reflect.Type) (reflect.Value, bool) {
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
	UintBinder = func(val string, typ reflect.Type) (reflect.Value, bool) {
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
	FloatBinder = func(val string, typ reflect.Type) (reflect.Value, bool) {
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
	StringBinder = func(val string, typ reflect.Type) (reflect.Value, bool) {
		return reflect.ValueOf(val), false
	}

	// Binds a boolean value, using the following formats :
	// "true" and "false"
	// "on" and "" (a checkbox)
	// "1" and "0" (why not)
	BoolBinder = func(val string, typ reflect.Type) (reflect.Value, bool) {
		v := strings.TrimSpace(strings.ToLower(val))
		switch v {
		case "true", "on", "1":
			return reflect.ValueOf(true), false
		}
		// Return false by default.
		return reflect.ValueOf(false), false
	}

	// Binds a comma separated list to a slice of srtangs.
	StringSliceBinder = func(val string, typ reflect.Type) (reflect.Value, bool) {
		split := strings.Split(val, ",")
		return SliceBinder(split, typ)
	}

	// Binds a slice
	SliceBinder = func(val []string, typ reflect.Type) (reflect.Value, bool) {
		resultArray := reflect.MakeSlice(typ, len(val), len(val))
		for index, item := range val {
			itemVal, _ := Bind(item, typ.Elem())
			resultArray.Index(index).Set(itemVal)
		}
		return resultArray, false
	}

	// Binds a pointer. If nothing found, returns nil.
	PointerBinder = func(val string, typ reflect.Type) (reflect.Value, bool) {
		pointerOf := typ.Elem()
		value, isNil := Bind(val, pointerOf)
		if isNil {
			return reflect.Zero(typ), true
		} else {
			return value.Addr(), false
		}

		return reflect.Zero(typ), true
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
}

// Takes the name and type of the desired parameter and constructs it
// from query string values.
// Returns the zero value of the type upon any sort of failure.
func BindFromValues(values url.Values, name string, typ reflect.Type) reflect.Value {
	if typ.Kind() == reflect.Slice {
		if val, found := values[name]; found {
			ret, _ := SliceBinder(val, typ)
			return ret
		}
		return reflect.Zero(typ)
	}
	ret, _ := Bind(values.Get(name), typ)
	return ret
}

// Takes the name and type of the desired parameter and constructs it
// from a map of values (eg: Gorilla Mux vars).
// Returns the zero value of the type upon any sort of failure.
func BindFromMap(values map[string]string, name string, typ reflect.Type) reflect.Value {
	if val, found := values[name]; found {
		ret, _ := Bind(val, typ)
		return ret
	}
	return reflect.Zero(typ)
}

// Constructs a value of a specific type
func Bind(val string, typ reflect.Type) (reflect.Value, bool) {
	if binder, found := TypeBinders[typ]; found {
		return binder(val, typ)
	} else {
		if binder, found := KindBinders[typ.Kind()]; found {
			return binder(val, typ)
		}
	}
	return reflect.Zero(typ), true
}

// Registers a custom binder for a specific type
func RegisterBinder(i interface{}, binder ValueBinder) {
	typ := reflect.ValueOf(i).Type()
	TypeBinders[typ] = binder
}