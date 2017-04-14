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

package monad

import "reflect"

// BindFunc is a function type used for Bind
type BindFunc func(interface{}) Monad

// Monad represents a generic monad
type Monad interface {
	AndThen(func(interface{}) Monad) Monad
	Return(i interface{}) Monad
}

// MPlus defines the interface for the MonadPlus class
type MPlus interface {
	Monad
	MZero() Monad
}

// Join takes a Monad (Monad (interface{})) and returns Monad (interface{})
func Join(m Monad) Monad {
	return m.AndThen(func(i interface{}) Monad { return i.(Monad) })
}

// Kleisli composition for monadic functions
func Kleisli(a, b func(i interface{}) Monad) func(i interface{}) Monad {
	return func(i interface{}) Monad { return b(i).AndThen(a) }
}

// FMap :: (a -> b) -> m a -> m b; it will panic if f is not a func with
// a single input and a single output
func FMap(f interface{}, m Monad) Monad {
	return m.AndThen(func(i interface{}) Monad {
		v := reflect.ValueOf(f).Call([]reflect.Value{reflect.ValueOf(i)})
		return m.Return(v[0].Interface())
	})
}
