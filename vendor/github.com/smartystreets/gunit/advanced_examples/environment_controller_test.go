package examples

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestEnvironmentControllerFixture(t *testing.T) {
	gunit.Run(new(EnvironmentControllerFixture), t)
}

type EnvironmentControllerFixture struct {
	*gunit.Fixture
	hardware   *FakeHVAC
	controller *EnvironmentController
}

func (this *EnvironmentControllerFixture) Setup() {
	this.hardware = NewFakeHardware()
	this.controller = NewController(this.hardware)
}

func (this *EnvironmentControllerFixture) TestEverythingTurnedOffAtStartup() {
	this.activateAllHardwareComponents()
	this.controller = NewController(this.hardware)
	this.assertAllHardwareComponentsAreDeactivated()
}

func (this *EnvironmentControllerFixture) TestEverythingOffWhenComfortable() {
	this.setupComfortableEnvironment()
	this.assertAllHardwareComponentsAreDeactivated()
}

func (this *EnvironmentControllerFixture) TestCoolerAndBlowerWhenHot() {
	this.setupHotEnvironment()
	this.AssertEqual("heater BLOWER COOLER low high", this.hardware.String())
}

func (this *EnvironmentControllerFixture) TestHeaterAndBlowerWhenCold() {
	this.setupColdEnvironment()
	this.AssertEqual("HEATER BLOWER cooler low high", this.hardware.String())
}

func (this *EnvironmentControllerFixture) TestHighAlarmOnIfAtThreshold() {
	this.setupBlazingEnvironment()
	this.AssertEqual("heater BLOWER COOLER low HIGH", this.hardware.String())
}

func (this *EnvironmentControllerFixture) TestLowAlarmOnIfAtThreshold() {
	this.setupFreezingEnvironment()
	this.AssertEqual("HEATER BLOWER cooler LOW high", this.hardware.String())
}

func (this *EnvironmentControllerFixture) setupComfortableEnvironment() {
	this.hardware.SetCurrentTemperature(COMFORTABLE)
	this.controller.Regulate()
}
func (this *EnvironmentControllerFixture) setupHotEnvironment() {
	this.hardware.SetCurrentTemperature(TOO_HOT)
	this.controller.Regulate()
}
func (this *EnvironmentControllerFixture) setupBlazingEnvironment() {
	this.hardware.SetCurrentTemperature(WAY_TOO_HOT)
	this.controller.Regulate()
}
func (this *EnvironmentControllerFixture) setupColdEnvironment() {
	this.hardware.SetCurrentTemperature(TOO_COLD)
	this.controller.Regulate()
}
func (this *EnvironmentControllerFixture) setupFreezingEnvironment() {
	this.hardware.SetCurrentTemperature(WAY_TOO_COLD)
	this.controller.Regulate()
}

func (this *EnvironmentControllerFixture) activateAllHardwareComponents() {
	this.hardware.ActivateBlower()
	this.hardware.ActivateHeater()
	this.hardware.ActivateCooler()
	this.hardware.ActivateHighTemperatureAlarm()
	this.hardware.ActivateLowTemperatureAlarm()
}

func (this *EnvironmentControllerFixture) assertAllHardwareComponentsAreDeactivated() {
	this.AssertEqual("heater blower cooler low high", this.hardware.String())
}
