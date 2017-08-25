package examples

type EnvironmentController struct {
	hardware HVAC
}

func NewController(hardware HVAC) *EnvironmentController {
	hardware.DeactivateBlower()
	hardware.DeactivateHeater()
	hardware.DeactivateCooler()
	hardware.DeactivateHighTemperatureAlarm()
	hardware.DeactivateLowTemperatureAlarm()
	return &EnvironmentController{hardware: hardware}
}

func (this *EnvironmentController) Regulate() {
	temperature := this.hardware.CurrentTemperature()

	if temperature >= WAY_TOO_HOT {
		this.hardware.ActivateHighTemperatureAlarm()
	} else if temperature <= WAY_TOO_COLD {
		this.hardware.ActivateLowTemperatureAlarm()
	}

	if temperature >= TOO_HOT {
		this.hardware.DeactivateHeater()
		this.hardware.ActivateBlower()
		this.hardware.ActivateCooler()
	} else if temperature <= TOO_COLD {
		this.hardware.DeactivateCooler()
		this.hardware.ActivateBlower()
		this.hardware.ActivateHeater()
	}
}

const (
	WAY_TOO_HOT  = 80
	TOO_HOT      = 70
	TOO_COLD     = 60
	WAY_TOO_COLD = 50
	COMFORTABLE  = (TOO_HOT + TOO_COLD) / 2
)
