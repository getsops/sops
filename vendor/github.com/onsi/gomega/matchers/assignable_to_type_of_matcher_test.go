package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("AssignableToTypeOf", func() {
	Context("When asserting assignability between types", func() {
		It("should do the right thing", func() {
			Expect(0).Should(BeAssignableToTypeOf(0))
			Expect(5).Should(BeAssignableToTypeOf(-1))
			Expect("foo").Should(BeAssignableToTypeOf("bar"))
			Expect(struct{ Foo string }{}).Should(BeAssignableToTypeOf(struct{ Foo string }{}))

			Expect(0).ShouldNot(BeAssignableToTypeOf("bar"))
			Expect(5).ShouldNot(BeAssignableToTypeOf(struct{ Foo string }{}))
			Expect("foo").ShouldNot(BeAssignableToTypeOf(42))
		})
	})

	Context("When asserting nil values", func() {
		It("should error", func() {
			success, err := (&AssignableToTypeOfMatcher{Expected: nil}).Match(nil)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
		})
	})
})
