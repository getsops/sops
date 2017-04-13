// +build js

package cookiejar_test

import "fmt"

func ExampleNew() {
	// network access not supported by GopherJS, and this test depends on httptest.NewServer

	fmt.Println(`After 1st request:
  Flavor: Chocolate Chip
After 2nd request:
  Flavor: Oatmeal Raisin`)
}
