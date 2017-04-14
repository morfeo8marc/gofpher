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

package either_test

import (
	"fmt"

	"github.com/rebeccaskinner/gofpher/either"
	"github.com/rebeccaskinner/gofpher/monad"
)

func toNum(i interface{}) either.EitherM {
	if num, ok := i.(int); ok {
		return either.RightM(num)
	}
	return either.LeftM(i)
}

func addOne(src interface{}) monad.Monad {
	num := src.(int)
	if num == 7 {
		return toNum("seven")
	}
	return toNum(num + 1)
}

func show(a, b either.EitherM) {
	fmt.Println(a)
	fmt.Println(b)
}

func Example_join() {
	a := either.RightM(either.RightM("foo"))
	fmt.Println(a)
	fmt.Println(monad.Join(a))
	a = either.LeftM(either.LeftM("foo"))
	fmt.Println(a)
	fmt.Println(monad.Join(a))
	//Output:
	// Right (Right (foo))
	// Right (foo)
	// Left (Left (foo))
	// Left (Left (foo))

}

func Example_either() {
	a := either.RightM(1)
	b := either.LeftM("foo")
	show(a, b)
	a1 := a.AndThen(addOne)
	b1 := b.AndThen(addOne)
	show(a1.(either.EitherM), b1.(either.EitherM))
	//Output:
	// Right (1)
	// Left (foo)
	// Right (2)
	// Left (foo)
}

func Example_chainAndThen() {
	plusOne := func(i interface{}) interface{} { return 1 + i.(int) }
	a := either.RightM(1).AndThen(addOne).AndThen(addOne).AndThen(addOne)
	fmt.Println(a)
	b := monad.FMap(plusOne, a)
	fmt.Println(b)
	//Output:
	// Right (4)
	// Right (5)
}

func Example_eitherFMap() {
	plusOne := func(i interface{}) interface{} { return 1 + i.(int) }
	a := either.ReturnM(1)
	b := monad.FMap(plusOne, a)
	fmt.Println(b)
	b = monad.FMap(plusOne, b)
	fmt.Println(b)
	b = monad.FMap(plusOne, b)
	fmt.Println(b)
	//Output:
	// Right (2)
	// Right (3)
	// Right (4)
}
