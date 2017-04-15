package functools

import (
	"fmt"
	"reflect"
)

// Function wraps a reflect.Value with a new function call syntax. It only
// supports functions with a single input value and a single return value. Use
// `Curry` to convert multi-arity functions into Function types.
//
// If your function returns two values, you may consider if it's appropriate to
// convert it to an EitherM or MaybeM.
type Function struct {
	f reflect.Value
}

// Call calls the inner wrapped composed function
func (f Function) Call(i interface{}) interface{} {
	return f.f.Call([]reflect.Value{reflect.ValueOf(i)})[0].Interface()
}

// Compose :: (a -> b) -> (b -> c) -> (a -> c); panics if a and be aren't
// functions that each take and return a single value
func Compose(a, b interface{}) Function {
	if c, ok := a.(Function); ok {
		return Compose(c.f.Interface(), b)
	}
	if c, ok := b.(Function); ok {
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
	return Function{newF}
}

// Wrap lifts a go function into a Function.  It panics if f is not a function.
func Wrap(f interface{}) Function {
	fTyp := reflect.TypeOf(f)
	if fTyp.Kind() != reflect.Func {
		panic(fmt.Sprintf("attempt to wrap non-function type %T", f))
	}
	if fTyp.NumIn() != 1 || fTyp.NumOut() != 1 {
		panic("attempt to wrap function with incorrect arity (must be 1/1)")
	}
	return Function{reflect.ValueOf(f)}
}

// Compose joins two functions
func (f Function) Compose(newF Function) Function {
	oldOut := f.f.Type().Out(0)
	newIn := newF.f.Type().In(0)
	if oldOut != newIn {
		panic(fmt.Sprintf("expected: a -> %v -> b; got (a -> %v) -> (%v -> b)", oldOut, oldOut, newIn))
	}
	if newF.f.Type().NumIn() != 1 || newF.f.Type().NumOut() != 1 {
		panic("compose got non-arity-1 values for in or out")
	}
	inTyp := []reflect.Type{f.f.Type().In(0)}
	outTyp := []reflect.Type{newF.f.Type().Out(0)}
	fTyp := reflect.FuncOf(inTyp, outTyp, false)
	composed := reflect.MakeFunc(fTyp, func(args []reflect.Value) []reflect.Value {
		outer := f.Call(args[0].Interface())
		inner := newF.Call(outer)
		return []reflect.Value{reflect.ValueOf(inner)}
	})
	return Function{composed}
}
