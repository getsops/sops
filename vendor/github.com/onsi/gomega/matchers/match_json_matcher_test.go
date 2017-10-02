package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("MatchJSONMatcher", func() {
	Context("When passed stringifiables", func() {
		It("should succeed if the JSON matches", func() {
			Ω("{}").Should(MatchJSON("{}"))
			Ω(`{"a":1}`).Should(MatchJSON(`{"a":1}`))
			Ω(`{
			             "a":1
			         }`).Should(MatchJSON(`{"a":1}`))
			Ω(`{"a":1, "b":2}`).Should(MatchJSON(`{"b":2, "a":1}`))
			Ω(`{"a":1}`).ShouldNot(MatchJSON(`{"b":2, "a":1}`))

			Ω(`{"a":"a", "b":"b"}`).ShouldNot(MatchJSON(`{"a":"a", "b":"b", "c":"c"}`))
			Ω(`{"a":"a", "b":"b", "c":"c"}`).ShouldNot(MatchJSON(`{"a":"a", "b":"b"}`))

			Ω(`{"a":null, "b":null}`).ShouldNot(MatchJSON(`{"c":"c", "d":"d"}`))
			Ω(`{"a":null, "b":null, "c":null}`).ShouldNot(MatchJSON(`{"a":null, "b":null, "d":null}`))
		})

		It("should work with byte arrays", func() {
			Ω([]byte("{}")).Should(MatchJSON([]byte("{}")))
			Ω("{}").Should(MatchJSON([]byte("{}")))
			Ω([]byte("{}")).Should(MatchJSON("{}"))
		})
	})

	Context("when a key mismatch is found", func() {
		It("reports the first found mismatch", func() {
			subject := MatchJSONMatcher{JSONToMatch: `5`}
			actual := `7`
			subject.Match(actual)

			failureMessage := subject.FailureMessage(`7`)
			Ω(failureMessage).ToNot(ContainSubstring("first mismatched key"))

			subject = MatchJSONMatcher{JSONToMatch: `{"a": 1, "b.g": {"c": 2, "1": ["hello", "see ya"]}}`}
			actual = `{"a": 1, "b.g": {"c": 2, "1": ["hello", "goodbye"]}}`
			subject.Match(actual)

			failureMessage = subject.FailureMessage(actual)
			Ω(failureMessage).To(ContainSubstring(`first mismatched key: "b.g"."1"[1]`))
		})
	})

	Context("when the expected is not valid JSON", func() {
		It("should error and explain why", func() {
			success, err := (&MatchJSONMatcher{JSONToMatch: `{}`}).Match(`oops`)
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("Actual 'oops' should be valid JSON"))
		})
	})

	Context("when the actual is not valid JSON", func() {
		It("should error and explain why", func() {
			success, err := (&MatchJSONMatcher{JSONToMatch: `oops`}).Match(`{}`)
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("Expected 'oops' should be valid JSON"))
		})
	})

	Context("when the expected is neither a string nor a stringer nor a byte array", func() {
		It("should error", func() {
			success, err := (&MatchJSONMatcher{JSONToMatch: 2}).Match("{}")
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("MatchJSONMatcher matcher requires a string, stringer, or []byte.  Got expected:\n    <int>: 2"))

			success, err = (&MatchJSONMatcher{JSONToMatch: nil}).Match("{}")
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("MatchJSONMatcher matcher requires a string, stringer, or []byte.  Got expected:\n    <nil>: nil"))
		})
	})

	Context("when the actual is neither a string nor a stringer nor a byte array", func() {
		It("should error", func() {
			success, err := (&MatchJSONMatcher{JSONToMatch: "{}"}).Match(2)
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("MatchJSONMatcher matcher requires a string, stringer, or []byte.  Got actual:\n    <int>: 2"))

			success, err = (&MatchJSONMatcher{JSONToMatch: "{}"}).Match(nil)
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("MatchJSONMatcher matcher requires a string, stringer, or []byte.  Got actual:\n    <nil>: nil"))
		})
	})
})
