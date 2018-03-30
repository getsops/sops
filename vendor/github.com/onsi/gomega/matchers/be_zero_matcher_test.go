package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BeZero", func() {
	It("should succeed if the passed in object is the zero value for its type", func() {
		Expect(nil).Should(BeZero())

		Expect("").Should(BeZero())
		Expect(" ").ShouldNot(BeZero())

		Expect(0).Should(BeZero())
		Expect(1).ShouldNot(BeZero())

		Expect(0.0).Should(BeZero())
		Expect(0.1).ShouldNot(BeZero())

		// Expect([]int{}).Should(BeZero())
		Expect([]int{1}).ShouldNot(BeZero())

		// Expect(map[string]int{}).Should(BeZero())
		Expect(map[string]int{"a": 1}).ShouldNot(BeZero())

		Expect(myCustomType{}).Should(BeZero())
		Expect(myCustomType{s: "a"}).ShouldNot(BeZero())
	})
})
