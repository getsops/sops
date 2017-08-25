package examples

import "strings"

type FakeHVAC struct {
	state       map[string]bool
	temperature int
}

func NewFakeHardware() *FakeHVAC {
	return &FakeHVAC{
		state: map[string]bool{
			"heater":          false,
			"blower":          false,
			"cooler":          false,
			"high-temp-alarm": false,
			"low-temp-alarm":  false,
		},
	}
}

func (this *FakeHVAC) ActivateHeater()               { this.state["heater"] = true }
func (this *FakeHVAC) ActivateBlower()               { this.state["blower"] = true }
func (this *FakeHVAC) ActivateCooler()               { this.state["cooler"] = true }
func (this *FakeHVAC) ActivateHighTemperatureAlarm() { this.state["high"] = true }
func (this *FakeHVAC) ActivateLowTemperatureAlarm()  { this.state["low"] = true }

func (this *FakeHVAC) DeactivateHeater()               { this.state["heater"] = false }
func (this *FakeHVAC) DeactivateBlower()               { this.state["blower"] = false }
func (this *FakeHVAC) DeactivateCooler()               { this.state["cooler"] = false }
func (this *FakeHVAC) DeactivateHighTemperatureAlarm() { this.state["high"] = false }
func (this *FakeHVAC) DeactivateLowTemperatureAlarm()  { this.state["low"] = false }

func (this *FakeHVAC) IsHeating() bool            { return this.state["heater"] }
func (this *FakeHVAC) IsBlowing() bool            { return this.state["blower"] }
func (this *FakeHVAC) IsCooling() bool            { return this.state["cooler"] }
func (this *FakeHVAC) HighTemperatureAlarm() bool { return this.state["high"] }
func (this *FakeHVAC) LowTemperatureAlarm() bool  { return this.state["low"] }

func (this *FakeHVAC) SetCurrentTemperature(value int) { this.temperature = value }
func (this *FakeHVAC) CurrentTemperature() int         { return this.temperature }

// String returns the status of each hardware component encoded in a single space-delimited string.
// UPPERCASE components are activated.
// lowercase components are deactivated.
func (this *FakeHVAC) String() string {
	current := []string{"heater", "blower", "cooler", "low", "high"}
	for i, component := range current {
		if this.state[component] {
			current[i] = strings.ToUpper(current[i])
		}
	}
	return strings.Join(current, " ")
}
