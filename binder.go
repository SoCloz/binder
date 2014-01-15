package binder

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// A Binder translates between string parameters and Go data structures.
type ValueBinder func(values url.Values, name string, typ reflect.Type) (reflect.Value, bool)

// The map of query string param name => struct attribute name
type StructAttributes map[string]string

var (
	// Lookup tables
	TypeBinders         = make(map[reflect.Type]ValueBinder)
	KindBinders         = make(map[reflect.Kind]ValueBinder)
	StructAttributesMap = make(map[reflect.Type]StructAttributes)

	// Binds an integer (signed)
	IntBinder = func(values url.Values, name string, typ reflect.Type) (reflect.Value, bool) {
		val := values.Get(name)
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
	UintBinder = func(values url.Values, name string, typ reflect.Type) (reflect.Value, bool) {
		val := values.Get(name)
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
	FloatBinder = func(values url.Values, name string, typ reflect.Type) (reflect.Value, bool) {
		val := GetValue(values, name)
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
	StringBinder = func(values url.Values, name string, typ reflect.Type) (reflect.Value, bool) {
		val := values.Get(name)
		return reflect.ValueOf(val), false
	}

	// Binds a boolean value, using the following formats :
	// "true" and "false"
	// "on" and "" (a checkbox)
	// "1" and "0" (why not)
	BoolBinder = func(values url.Values, name string, typ reflect.Type) (reflect.Value, bool) {
		val := values.Get(name)
		v := strings.TrimSpace(strings.ToLower(val))
		switch v {
		case "yes", "true", "on", "1":
			return reflect.ValueOf(true), false
		}
		// Return false by default.
		return reflect.ValueOf(false), false
	}

	// Binds a comma separated list to a slice of strings.
	StringSliceBinder = func(values url.Values, name string, typ reflect.Type) (reflect.Value, bool) {
		val := values.Get(name)
		if val == "" {
			return reflect.MakeSlice(typ, 0, 0), true
		}
		split := strings.Split(val, ",")
		result := reflect.MakeSlice(typ, len(split), len(split))
		for index, item := range split {
			values.Set(name, item)
			itemVal, _ := Bind(values, name, typ.Elem())
			result.Index(index).Set(itemVal)
		}
		return result, false
	}

	// Binds a set of parameters to a struct.
	StructBinder = func(values url.Values, name string, typ reflect.Type) (reflect.Value, bool) {
		result := reflect.New(typ).Elem()
		for fieldName, attrName := range StructAttributesMap[typ] {
			fieldValue := result.FieldByName(fieldName)
			boundVal, _ := Bind(values, attrName, fieldValue.Type())
			fieldValue.Set(boundVal)
		}
		return result, false
	}

	// Binds a pointer. If nothing found, returns nil.
	PointerBinder = func(values url.Values, name string, typ reflect.Type) (reflect.Value, bool) {
		pointerOf := typ.Elem()
		value, isNil := Bind(values, name, pointerOf)
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
	KindBinders[reflect.Struct] = StructBinder
}

// Takes the name and type of the desired parameter and constructs it
// from query string values.
// Returns the zero value of the type upon any sort of failure.
func GetBoundValue(values url.Values, name string, typ reflect.Type) reflect.Value {
	ret, _ := Bind(values, name, typ)
	return ret
}

// Constructs a value of a specific type
func Bind(values url.Values, name string, typ reflect.Type) (reflect.Value, bool) {
	if binder, found := TypeBinders[typ]; found {
		return binder(values, name, typ)
	} else {
		if binder, found := KindBinders[typ.Kind()]; found {
			return binder(values, name, typ)
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

// Get a single value from url.Values (query string or pat url parameter)
func GetValue(values url.Values, name string) string {
	// pat url parameter
	if _, found := values[":"+name]; found {
		return values.Get(":" + name)
	}
	// query string
	return values.Get(name)
}
