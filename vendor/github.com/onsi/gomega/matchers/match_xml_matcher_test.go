package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("MatchXMLMatcher", func() {

	var (
		sample_01 = readFileContents("test_data/xml/sample_01.xml")
		sample_02 = readFileContents("test_data/xml/sample_02.xml")
		sample_03 = readFileContents("test_data/xml/sample_03.xml")
		sample_04 = readFileContents("test_data/xml/sample_04.xml")
		sample_05 = readFileContents("test_data/xml/sample_05.xml")
		sample_06 = readFileContents("test_data/xml/sample_06.xml")
		sample_07 = readFileContents("test_data/xml/sample_07.xml")
		sample_08 = readFileContents("test_data/xml/sample_08.xml")
		sample_09 = readFileContents("test_data/xml/sample_09.xml")
		sample_10 = readFileContents("test_data/xml/sample_10.xml")
		sample_11 = readFileContents("test_data/xml/sample_11.xml")
	)

	Context("When passed stringifiables", func() {
		It("should succeed if the XML matches", func() {
			Ω(sample_01).Should(MatchXML(sample_01))    // same XML
			Ω(sample_01).Should(MatchXML(sample_02))    // same XML with blank lines
			Ω(sample_01).Should(MatchXML(sample_03))    // same XML with different formatting
			Ω(sample_01).ShouldNot(MatchXML(sample_04)) // same structures with different values
			Ω(sample_01).ShouldNot(MatchXML(sample_05)) // different structures
			Ω(sample_06).ShouldNot(MatchXML(sample_07)) // same xml names with different namespaces
			Ω(sample_07).ShouldNot(MatchXML(sample_08)) // same structures with different values
			Ω(sample_09).ShouldNot(MatchXML(sample_10)) // same structures with different attribute values
			Ω(sample_11).Should(MatchXML(sample_11))    // with non UTF-8 encoding
		})

		It("should work with byte arrays", func() {
			Ω([]byte(sample_01)).Should(MatchXML([]byte(sample_01)))
			Ω([]byte(sample_01)).Should(MatchXML(sample_01))
			Ω(sample_01).Should(MatchXML([]byte(sample_01)))
		})
	})

	Context("when the expected is not valid XML", func() {
		It("should error and explain why", func() {
			success, err := (&MatchXMLMatcher{XMLToMatch: sample_01}).Match(`oops`)
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("Actual 'oops' should be valid XML"))
		})
	})

	Context("when the actual is not valid XML", func() {
		It("should error and explain why", func() {
			success, err := (&MatchXMLMatcher{XMLToMatch: `oops`}).Match(sample_01)
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("Expected 'oops' should be valid XML"))
		})
	})

	Context("when the expected is neither a string nor a stringer nor a byte array", func() {
		It("should error", func() {
			success, err := (&MatchXMLMatcher{XMLToMatch: 2}).Match(sample_01)
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("MatchXMLMatcher matcher requires a string, stringer, or []byte.  Got expected:\n    <int>: 2"))

			success, err = (&MatchXMLMatcher{XMLToMatch: nil}).Match(sample_01)
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("MatchXMLMatcher matcher requires a string, stringer, or []byte.  Got expected:\n    <nil>: nil"))
		})
	})

	Context("when the actual is neither a string nor a stringer nor a byte array", func() {
		It("should error", func() {
			success, err := (&MatchXMLMatcher{XMLToMatch: sample_01}).Match(2)
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("MatchXMLMatcher matcher requires a string, stringer, or []byte.  Got actual:\n    <int>: 2"))

			success, err = (&MatchXMLMatcher{XMLToMatch: sample_01}).Match(nil)
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("MatchXMLMatcher matcher requires a string, stringer, or []byte.  Got actual:\n    <nil>: nil"))
		})
	})
})
