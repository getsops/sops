package assertions

import (
	"testing"

	"github.com/smartystreets/assertions/internal/unit"
)

func TestEqualityFixture(t *testing.T) {
	unit.Run(new(EqualityFixture), t)
}

type EqualityFixture struct {
	*unit.Fixture
}

func (this *EqualityFixture) TestNilNil() {
	spec := newEqualityMethodSpecification(nil, nil)
	this.So(spec.IsSatisfied(), ShouldBeFalse)
}

func (this *EqualityFixture) TestEligible1() {
	a := Eligible1{"hi"}
	b := Eligible1{"hi"}
	specification := newEqualityMethodSpecification(a, b)
	this.So(specification.IsSatisfied(), ShouldBeTrue)
	this.So(specification.AreEqual(), ShouldBeTrue)
}

func (this *EqualityFixture) TestAreEqual() {
	a := Eligible1{"hi"}
	b := Eligible1{"hi"}
	specification := newEqualityMethodSpecification(a, b)
	this.So(specification.IsSatisfied(), ShouldBeTrue)
	this.So(specification.AreEqual(), ShouldBeTrue)
}

func (this *EqualityFixture) TestAreNotEqual() {
	a := Eligible1{"hi"}
	b := Eligible1{"bye"}
	specification := newEqualityMethodSpecification(a, b)
	this.So(specification.IsSatisfied(), ShouldBeTrue)
	this.So(specification.AreEqual(), ShouldBeFalse)
}

func (this *EqualityFixture) TestEligible2() {
	a := Eligible2{"hi"}
	b := Eligible2{"hi"}
	specification := newEqualityMethodSpecification(a, b)
	this.So(specification.IsSatisfied(), ShouldBeTrue)
}

func (this *EqualityFixture) TestEligible1_PointerReceiver() {
	a := &Eligible1{"hi"}
	b := Eligible1{"hi"}
	this.So(a.Equal(b), ShouldBeTrue)
	specification := newEqualityMethodSpecification(a, b)
	this.So(specification.IsSatisfied(), ShouldBeTrue)
}

func (this *EqualityFixture) TestIneligible_PrimitiveTypes() {
	specification := newEqualityMethodSpecification(1, 1)
	this.So(specification.IsSatisfied(), ShouldBeFalse)
}

func (this *EqualityFixture) TestIneligible_DisparateTypes() {
	a := Eligible1{"hi"}
	b := Eligible2{"hi"}
	specification := newEqualityMethodSpecification(a, b)
	this.So(specification.IsSatisfied(), ShouldBeFalse)
}

func (this *EqualityFixture) TestIneligible_NoEqualMethod() {
	a := Ineligible_NoEqualMethod{}
	b := Ineligible_NoEqualMethod{}
	specification := newEqualityMethodSpecification(a, b)
	this.So(specification.IsSatisfied(), ShouldBeFalse)
}

func (this *EqualityFixture) TestIneligible_EqualMethodReceivesNoInput() {
	a := Ineligible_EqualMethodNoInputs{}
	b := Ineligible_EqualMethodNoInputs{}
	specification := newEqualityMethodSpecification(a, b)
	this.So(specification.IsSatisfied(), ShouldBeFalse)
}

func (this *EqualityFixture) TestIneligible_EqualMethodReceivesTooManyInputs() {
	a := Ineligible_EqualMethodTooManyInputs{}
	b := Ineligible_EqualMethodTooManyInputs{}
	specification := newEqualityMethodSpecification(a, b)
	this.So(specification.IsSatisfied(), ShouldBeFalse)
}

func (this *EqualityFixture) TestIneligible_EqualMethodReceivesWrongInput() {
	a := Ineligible_EqualMethodWrongInput{}
	b := Ineligible_EqualMethodWrongInput{}
	specification := newEqualityMethodSpecification(a, b)
	this.So(specification.IsSatisfied(), ShouldBeFalse)
}

func (this *EqualityFixture) TestIneligible_EqualMethodReturnsNoOutputs() {
	a := Ineligible_EqualMethodNoOutputs{}
	b := Ineligible_EqualMethodNoOutputs{}
	specification := newEqualityMethodSpecification(a, b)
	this.So(specification.IsSatisfied(), ShouldBeFalse)
}

func (this *EqualityFixture) TestIneligible_EqualMethodReturnsTooManyOutputs() {
	a := Ineligible_EqualMethodTooManyOutputs{}
	b := Ineligible_EqualMethodTooManyOutputs{}
	specification := newEqualityMethodSpecification(a, b)
	this.So(specification.IsSatisfied(), ShouldBeFalse)
}

func (this *EqualityFixture) TestIneligible_EqualMethodReturnsWrongOutputs() {
	a := Ineligible_EqualMethodWrongOutput{}
	b := Ineligible_EqualMethodWrongOutput{}
	specification := newEqualityMethodSpecification(a, b)
	this.So(specification.IsSatisfied(), ShouldBeFalse)
}

func (this *EqualityFixture) TestEligibleAsymmetric_EqualMethodResultDiffersWhenArgumentsInverted() {
	a := EligibleAsymmetric{a: 0}
	b := EligibleAsymmetric{a: 1}
	specification := newEqualityMethodSpecification(a, b)
	this.So(specification.IsSatisfied(), ShouldBeTrue)
	this.So(specification.AreEqual(), ShouldBeFalse)
}

/**************************************************************************/

type (
	Eligible1                            struct{ a string }
	Eligible2                            struct{ a string }
	EligibleAsymmetric                   struct{ a int }
	Ineligible_NoEqualMethod             struct{}
	Ineligible_EqualMethodNoInputs       struct{}
	Ineligible_EqualMethodNoOutputs      struct{}
	Ineligible_EqualMethodTooManyInputs  struct{}
	Ineligible_EqualMethodTooManyOutputs struct{}
	Ineligible_EqualMethodWrongInput     struct{}
	Ineligible_EqualMethodWrongOutput    struct{}
)

func (this Eligible1) Equal(that Eligible1) bool { return this.a == that.a }
func (this Eligible2) Equal(that Eligible2) bool { return this.a == that.a }
func (this EligibleAsymmetric) Equal(that EligibleAsymmetric) bool {
	return this.a == 0
}
func (this Ineligible_EqualMethodNoInputs) Equal() bool                                    { return true }
func (this Ineligible_EqualMethodNoOutputs) Equal(that Ineligible_EqualMethodNoOutputs)    {}
func (this Ineligible_EqualMethodTooManyInputs) Equal(a, b bool) bool                      { return true }
func (this Ineligible_EqualMethodTooManyOutputs) Equal(bool) (bool, bool)                  { return true, true }
func (this Ineligible_EqualMethodWrongInput) Equal(a string) bool                          { return true }
func (this Ineligible_EqualMethodWrongOutput) Equal(Ineligible_EqualMethodWrongOutput) int { return 0 }
