package matchers_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

func Erroring() error {
	return errors.New("bam")
}

func NotErroring() error {
	return nil
}

type AnyType struct{}

func Invalid() *AnyType {
	return nil
}

var _ = Describe("Succeed", func() {
	It("should succeed if the function succeeds", func() {
		Expect(NotErroring()).Should(Succeed())
	})

	It("should succeed (in the negated) if the function errored", func() {
		Expect(Erroring()).ShouldNot(Succeed())
	})

	It("should not if passed a non-error", func() {
		success, err := (&SucceedMatcher{}).Match(Invalid())
		Expect(success).Should(BeFalse())
		Expect(err).Should(MatchError("Expected an error-type.  Got:\n    <*matchers_test.AnyType | 0x0>: nil"))
	})

	It("doesn't support non-error type", func() {
		success, err := (&SucceedMatcher{}).Match(AnyType{})
		Expect(success).Should(BeFalse())
		Expect(err).Should(MatchError("Expected an error-type.  Got:\n    <matchers_test.AnyType>: {}"))
	})

	It("doesn't support non-error pointer type", func() {
		success, err := (&SucceedMatcher{}).Match(&AnyType{})
		Expect(success).Should(BeFalse())
		Expect(err).Should(MatchError(MatchRegexp(`Expected an error-type.  Got:\n    <*matchers_test.AnyType | 0x[[:xdigit:]]+>: {}`)))
	})

	It("should not succeed with pointer types that conform to error interface", func() {
		err := &CustomErr{"ohai"}
		Expect(err).ShouldNot(Succeed())
	})

	It("should succeed with nil pointers to types that conform to error interface", func() {
		var err *CustomErr = nil
		Expect(err).Should(Succeed())
	})

})
