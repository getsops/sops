package examples

type HVAC interface {
	ActivateHeater()
	ActivateBlower()
	ActivateCooler()
	ActivateHighTemperatureAlarm()
	ActivateLowTemperatureAlarm()

	DeactivateHeater()
	DeactivateBlower()
	DeactivateCooler()
	DeactivateHighTemperatureAlarm()
	DeactivateLowTemperatureAlarm()

	IsHeating() bool
	IsBlowing() bool
	IsCooling() bool
	HighTemperatureAlarm() bool
	LowTemperatureAlarm() bool

	CurrentTemperature() int
}
