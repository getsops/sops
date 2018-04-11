package testingtsupport_test

import (
	. "github.com/onsi/gomega"

	"fmt"
	"testing"
)

type FakeT struct {
	LastCall string
}

func (f *FakeT) Fatalf(format string, args ...interface{}) {
	f.LastCall = fmt.Sprintf(format, args...)
}

func TestTestingT(t *testing.T) {
	RegisterTestingT(t)
	Expect(true).Should(BeTrue())
}

func TestGomegaWithT(t *testing.T) {
	g := NewGomegaWithT(t)

	f := &FakeT{}
	testG := NewGomegaWithT(f)

	testG.Expect("foo").To(Equal("foo"))
	g.Expect(f.LastCall).To(BeZero())

	testG.Expect("foo").To(Equal("bar"))
	g.Expect(f.LastCall).To(ContainSubstring("<string>: foo"))
}
