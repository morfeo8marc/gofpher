package functools

import "reflect"

// Composed wraps a reflect.Value with a new function call syntax
type Composed struct {
	f reflect.Value
}

// Call calls the inner wrapped composed function
func (c Composed) Call(i interface{}) interface{} {
	return c.f.Call([]reflect.Value{reflect.ValueOf(i)})[0].Interface()
}

// Compose :: (a -> b) -> (b -> c) -> (a -> c); panics if a and be aren't
// functions that each take and return a single value
func Compose(a, b interface{}) Composed {
	if c, ok := a.(Composed); ok {
		return Compose(c.f.Interface(), b)
	}
	if c, ok := b.(Composed); ok {
		return Compose(a, c.f.Interface())
	}
	aVal := reflect.ValueOf(a)
	aType := reflect.TypeOf(a)
	bVal := reflect.ValueOf(b)
	bType := reflect.TypeOf(b)
	newFuncIn := []reflect.Type{aType.In(0)}
	newFuncOut := []reflect.Type{bType.Out(0)}
	newFType := reflect.FuncOf(newFuncIn, newFuncOut, false)
	newF := reflect.MakeFunc(newFType, func(args []reflect.Value) []reflect.Value {
		return bVal.Call(aVal.Call(args))
	})
	return Composed{newF}
}
