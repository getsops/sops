// +build js

package reflect_test

import "fmt"

func ExampleStructOf() {
	// GopherJS does not implement reflect.addReflectOff needed for this test.
	// See https://github.com/gopherjs/gopherjs/issues/499

	fmt.Println(`value: &{Height:0.4 Age:2}
json:  {"height":0.4,"age":2}
value: &{Height:1.5 Age:10}`)
}
