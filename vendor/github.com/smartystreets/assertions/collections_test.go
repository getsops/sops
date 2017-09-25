package assertions

import (
	"fmt"
	"time"
)

func (this *AssertionsFixture) TestShouldContainKey() {
	this.fail(so(map[int]int{}, ShouldContainKey), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so(map[int]int{}, ShouldContainKey, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	this.fail(so(Thing1{}, ShouldContainKey, 1), "You must provide a valid map type (was assertions.Thing1)!")
	this.fail(so(nil, ShouldContainKey, 1), "You must provide a valid map type (was <nil>)!")
	this.fail(so(map[int]int{1: 41}, ShouldContainKey, 2), "Expected the map[int]int to contain the key: [2] (but it didn't)!")

	this.pass(so(map[int]int{1: 41}, ShouldContainKey, 1))
	this.pass(so(map[int]int{1: 41, 2: 42, 3: 43}, ShouldContainKey, 2))
}

func (this *AssertionsFixture) TestShouldNotContainKey() {
	this.fail(so(map[int]int{}, ShouldNotContainKey), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so(map[int]int{}, ShouldNotContainKey, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	this.fail(so(Thing1{}, ShouldNotContainKey, 1), "You must provide a valid map type (was assertions.Thing1)!")
	this.fail(so(nil, ShouldNotContainKey, 1), "You must provide a valid map type (was <nil>)!")
	this.fail(so(map[int]int{1: 41}, ShouldNotContainKey, 1), "Expected the map[int]int NOT to contain the key: [1] (but it did)!")
	this.pass(so(map[int]int{1: 41}, ShouldNotContainKey, 2))
}

func (this *AssertionsFixture) TestShouldContain() {
	this.fail(so([]int{}, ShouldContain), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so([]int{}, ShouldContain, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	this.fail(so(Thing1{}, ShouldContain, 1), "You must provide a valid container (was assertions.Thing1)!")
	this.fail(so(nil, ShouldContain, 1), "You must provide a valid container (was <nil>)!")
	this.fail(so([]int{1}, ShouldContain, 2), "Expected the container ([]int) to contain: '2' (but it didn't)!")
	this.fail(so([][]int{{1}}, ShouldContain, []int{2}), "Expected the container ([][]int) to contain: '[2]' (but it didn't)!")

	this.pass(so([]int{1}, ShouldContain, 1))
	this.pass(so([]int{1, 2, 3}, ShouldContain, 2))
	this.pass(so([][]int{{1}, {2}, {3}}, ShouldContain, []int{2}))
}

func (this *AssertionsFixture) TestShouldNotContain() {
	this.fail(so([]int{}, ShouldNotContain), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so([]int{}, ShouldNotContain, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	this.fail(so(Thing1{}, ShouldNotContain, 1), "You must provide a valid container (was assertions.Thing1)!")
	this.fail(so(nil, ShouldNotContain, 1), "You must provide a valid container (was <nil>)!")

	this.fail(so([]int{1}, ShouldNotContain, 1), "Expected the container ([]int) NOT to contain: '1' (but it did)!")
	this.fail(so([]int{1, 2, 3}, ShouldNotContain, 2), "Expected the container ([]int) NOT to contain: '2' (but it did)!")
	this.fail(so([][]int{{1}, {2}, {3}}, ShouldNotContain, []int{2}), "Expected the container ([][]int) NOT to contain: '[2]' (but it did)!")

	this.pass(so([]int{1}, ShouldNotContain, 2))
	this.pass(so([][]int{{1}, {2}, {3}}, ShouldNotContain, []int{4}))
}

func (this *AssertionsFixture) TestShouldBeIn() {
	this.fail(so(4, ShouldBeIn), needNonEmptyCollection)

	container := []int{1, 2, 3, 4}
	this.pass(so(4, ShouldBeIn, container))
	this.pass(so(4, ShouldBeIn, 1, 2, 3, 4))
	this.pass(so([]int{4}, ShouldBeIn, [][]int{{1}, {2}, {3}, {4}}))
	this.pass(so([]int{4}, ShouldBeIn, []int{1}, []int{2}, []int{3}, []int{4}))

	this.fail(so(4, ShouldBeIn, 1, 2, 3), "Expected '4' to be in the container ([]interface {}), but it wasn't!")
	this.fail(so(4, ShouldBeIn, []int{1, 2, 3}), "Expected '4' to be in the container ([]int), but it wasn't!")
	this.fail(so([]int{4}, ShouldBeIn, []int{1}, []int{2}, []int{3}), "Expected '[4]' to be in the container ([]interface {}), but it wasn't!")
	this.fail(so([]int{4}, ShouldBeIn, [][]int{{1}, {2}, {3}}), "Expected '[4]' to be in the container ([][]int), but it wasn't!")
}

func (this *AssertionsFixture) TestShouldNotBeIn() {
	this.fail(so(4, ShouldNotBeIn), needNonEmptyCollection)

	container := []int{1, 2, 3, 4}
	this.pass(so(42, ShouldNotBeIn, container))
	this.pass(so(42, ShouldNotBeIn, 1, 2, 3, 4))
	this.pass(so([]int{42}, ShouldNotBeIn, []int{1}, []int{2}, []int{3}, []int{4}))
	this.pass(so([]int{42}, ShouldNotBeIn, [][]int{{1}, {2}, {3}, {4}}))

	this.fail(so(2, ShouldNotBeIn, 1, 2, 3), "Expected '2' NOT to be in the container ([]interface {}), but it was!")
	this.fail(so(2, ShouldNotBeIn, []int{1, 2, 3}), "Expected '2' NOT to be in the container ([]int), but it was!")
	this.fail(so([]int{2}, ShouldNotBeIn, []int{1}, []int{2}, []int{3}), "Expected '[2]' NOT to be in the container ([]interface {}), but it was!")
	this.fail(so([]int{2}, ShouldNotBeIn, [][]int{{1}, {2}, {3}}), "Expected '[2]' NOT to be in the container ([][]int), but it was!")
}

func (this *AssertionsFixture) TestShouldBeEmpty() {
	this.fail(so(1, ShouldBeEmpty, 2, 3), "This assertion requires exactly 0 comparison values (you provided 2).")

	this.pass(so([]int{}, ShouldBeEmpty))           // empty slice
	this.pass(so([][]int{}, ShouldBeEmpty))         // empty slice
	this.pass(so([]interface{}{}, ShouldBeEmpty))   // empty slice
	this.pass(so(map[string]int{}, ShouldBeEmpty))  // empty map
	this.pass(so("", ShouldBeEmpty))                // empty string
	this.pass(so(&[]int{}, ShouldBeEmpty))          // pointer to empty slice
	this.pass(so(&[0]int{}, ShouldBeEmpty))         // pointer to empty array
	this.pass(so(nil, ShouldBeEmpty))               // nil
	this.pass(so(make(chan string), ShouldBeEmpty)) // empty channel

	this.fail(so([]int{1}, ShouldBeEmpty), "Expected [1] to be empty (but it wasn't)!")                      // non-empty slice
	this.fail(so([][]int{{1}}, ShouldBeEmpty), "Expected [[1]] to be empty (but it wasn't)!")                // non-empty slice
	this.fail(so([]interface{}{1}, ShouldBeEmpty), "Expected [1] to be empty (but it wasn't)!")              // non-empty slice
	this.fail(so(map[string]int{"hi": 0}, ShouldBeEmpty), "Expected map[hi:0] to be empty (but it wasn't)!") // non-empty map
	this.fail(so("hi", ShouldBeEmpty), "Expected hi to be empty (but it wasn't)!")                           // non-empty string
	this.fail(so(&[]int{1}, ShouldBeEmpty), "Expected &[1] to be empty (but it wasn't)!")                    // pointer to non-empty slice
	this.fail(so(&[1]int{1}, ShouldBeEmpty), "Expected &[1] to be empty (but it wasn't)!")                   // pointer to non-empty array
	c := make(chan int, 1)                                                                                   // non-empty channel
	go func() { c <- 1 }()
	time.Sleep(time.Millisecond)
	this.fail(so(c, ShouldBeEmpty), fmt.Sprintf("Expected %+v to be empty (but it wasn't)!", c))
}

func (this *AssertionsFixture) TestShouldNotBeEmpty() {
	this.fail(so(1, ShouldNotBeEmpty, 2, 3), "This assertion requires exactly 0 comparison values (you provided 2).")

	this.fail(so([]int{}, ShouldNotBeEmpty), "Expected [] to NOT be empty (but it was)!")             // empty slice
	this.fail(so([]interface{}{}, ShouldNotBeEmpty), "Expected [] to NOT be empty (but it was)!")     // empty slice
	this.fail(so(map[string]int{}, ShouldNotBeEmpty), "Expected map[] to NOT be empty (but it was)!") // empty map
	this.fail(so("", ShouldNotBeEmpty), "Expected  to NOT be empty (but it was)!")                    // empty string
	this.fail(so(&[]int{}, ShouldNotBeEmpty), "Expected &[] to NOT be empty (but it was)!")           // pointer to empty slice
	this.fail(so(&[0]int{}, ShouldNotBeEmpty), "Expected &[] to NOT be empty (but it was)!")          // pointer to empty array
	this.fail(so(nil, ShouldNotBeEmpty), "Expected <nil> to NOT be empty (but it was)!")              // nil
	c := make(chan int, 0)                                                                            // non-empty channel
	this.fail(so(c, ShouldNotBeEmpty), fmt.Sprintf("Expected %+v to NOT be empty (but it was)!", c))  // empty channel

	this.pass(so([]int{1}, ShouldNotBeEmpty))                // non-empty slice
	this.pass(so([]interface{}{1}, ShouldNotBeEmpty))        // non-empty slice
	this.pass(so(map[string]int{"hi": 0}, ShouldNotBeEmpty)) // non-empty map
	this.pass(so("hi", ShouldNotBeEmpty))                    // non-empty string
	this.pass(so(&[]int{1}, ShouldNotBeEmpty))               // pointer to non-empty slice
	this.pass(so(&[1]int{1}, ShouldNotBeEmpty))              // pointer to non-empty array
	c = make(chan int, 1)
	go func() { c <- 1 }()
	time.Sleep(time.Millisecond)
	this.pass(so(c, ShouldNotBeEmpty))
}

func (this *AssertionsFixture) TestShouldHaveLength() {
	this.fail(so(1, ShouldHaveLength, 2), "You must provide a valid container (was int)!")
	this.fail(so(nil, ShouldHaveLength, 1), "You must provide a valid container (was <nil>)!")
	this.fail(so("hi", ShouldHaveLength, float64(1.0)), "You must provide a valid integer (was float64)!")
	this.fail(so([]string{}, ShouldHaveLength), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so([]string{}, ShouldHaveLength, 1, 2), "This assertion requires exactly 1 comparison values (you provided 2).")
	this.fail(so([]string{}, ShouldHaveLength, -10), "You must provide a valid positive integer (was -10)!")

	this.fail(so([]int{}, ShouldHaveLength, 1), "Expected [] (length: 0) to have length equal to '1', but it wasn't!")             // empty slice
	this.fail(so([]interface{}{}, ShouldHaveLength, 1), "Expected [] (length: 0) to have length equal to '1', but it wasn't!")     // empty slice
	this.fail(so(map[string]int{}, ShouldHaveLength, 1), "Expected map[] (length: 0) to have length equal to '1', but it wasn't!") // empty map
	this.fail(so("", ShouldHaveLength, 1), "Expected  (length: 0) to have length equal to '1', but it wasn't!")                    // empty string
	this.fail(so(&[]int{}, ShouldHaveLength, 1), "Expected &[] (length: 0) to have length equal to '1', but it wasn't!")           // pointer to empty slice
	this.fail(so(&[0]int{}, ShouldHaveLength, 1), "Expected &[] (length: 0) to have length equal to '1', but it wasn't!")          // pointer to empty array
	c := make(chan int, 0)                                                                                                         // non-empty channel
	this.fail(so(c, ShouldHaveLength, 1), fmt.Sprintf("Expected %+v (length: 0) to have length equal to '1', but it wasn't!", c))
	c = make(chan int) // empty channel
	this.fail(so(c, ShouldHaveLength, 1), fmt.Sprintf("Expected %+v (length: 0) to have length equal to '1', but it wasn't!", c))

	this.pass(so([]int{1}, ShouldHaveLength, 1))                // non-empty slice
	this.pass(so([]interface{}{1}, ShouldHaveLength, 1))        // non-empty slice
	this.pass(so(map[string]int{"hi": 0}, ShouldHaveLength, 1)) // non-empty map
	this.pass(so("hi", ShouldHaveLength, 2))                    // non-empty string
	this.pass(so(&[]int{1}, ShouldHaveLength, 1))               // pointer to non-empty slice
	this.pass(so(&[1]int{1}, ShouldHaveLength, 1))              // pointer to non-empty array
	c = make(chan int, 1)
	go func() { c <- 1 }()
	time.Sleep(time.Millisecond)
	this.pass(so(c, ShouldHaveLength, 1))
	this.pass(so(c, ShouldHaveLength, uint(1)))

}
