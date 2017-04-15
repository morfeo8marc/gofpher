// Copyright © 2016 Asteris, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package maybe

import (
	"fmt"
	"reflect"

	"github.com/rebeccaskinner/gofpher/monad"
)

// Just is a plain value
type just struct {
	Val interface{}
}

// Nothing is the nil value for the maybe
type nothing struct{}

// Maybe provides an implementation of monad.Monad
type Maybe struct {
	internal interface{}
}

// AndThen provides the monadic implementaiton of AndThen
func (m Maybe) AndThen(f func(interface{}) monad.Monad) monad.Monad {
	if asJust, ok := m.internal.(just); ok {
		return f(asJust.Val)
	}
	return m
}

// LogAndThen provides the monadic implementaiton of LogAndThen
func (m Maybe) LogAndThen(f func(interface{}) monad.Monad, logger func(interface{})) monad.Monad {
	logger(m)
	if asJust, ok := m.internal.(just); ok {
		return f(asJust.Val)
	}
	return m
}

func (m Maybe) LiftM(f interface{}) monad.Monad {
	return monad.FMap(f, m)
}

func (m Maybe) Next(fnc interface{}) monad.Monad {
	f := reflect.ValueOf(fnc)
	if f.Type().Kind() != reflect.Func {
		return Nothing()
	}
	if f.Type().NumOut() == 1 {
		return m.LiftM(fnc)
	}
	if f.Type().NumOut() == 2 {
		return m.AndThen(WrapMaybe(fnc))
	}
	return Nothing()
}

func WrapMaybe(f interface{}) func(interface{}) monad.Monad {
	errF := func(s string) func(interface{}) monad.Monad {
		return func(interface{}) monad.Monad {
			return Nothing()
		}
	}
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return errF(fmt.Sprintf("expected function but got %T", f))
	}

	if t.NumIn() != 1 {
		return errF(fmt.Sprintf("function should have input arity of 1"))
	}

	if t.NumOut() != 2 {
		return errF(fmt.Sprintf("function should have an output arity of 2"))
	}

	errorType := reflect.TypeOf((*error)(nil)).Elem()

	if !t.Out(1).Implements(errorType) {
		return errF(fmt.Sprintf("function's second return value should be an error"))
	}

	return func(i interface{}) monad.Monad {
		res := reflect.ValueOf(f).Call([]reflect.Value{reflect.ValueOf(i)})
		if res[1].IsNil() {
			return Just(res[0].Interface())
		}
		return Nothing()
	}
}

// Return provides the monadic implementation of Return
func (m Maybe) Return(i interface{}) monad.Monad {
	return Maybe{internal: just{Val: i}}
}

// Nothing creates a Nothing value
func Nothing() monad.Monad {
	return Maybe{internal: nothing{}}
}

// Just creates a just value
func Just(i interface{}) monad.Monad {
	return Maybe{internal: just{Val: i}}
}

// IsJust returns true if the value is a just value
func (m Maybe) IsJust() bool {
	_, ok := m.internal.(just)
	return ok
}

// FromJust gets the just value and panics if it's nothing
func (m Maybe) FromJust() interface{} {
	return m.internal.(just).Val
}

// FromMaybe gets the just value or returns the default
func (m Maybe) FromMaybe(defaultVal interface{}) interface{} {
	if m.IsJust() {
		return m.FromJust()
	}
	return defaultVal
}

func (m Maybe) String() string {
	if m.IsJust() {
		return fmt.Sprintf("Just %v", m.FromJust())
	}
	return "Nothing"
}
