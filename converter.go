package main

func convertDrive(value int) string {
	drive := map[int]string{0: "OFF", 1: "ON", 2: "TESTRUN", 4: "ON", 5: "ON", 6: "OFF"}
	if drive[value] == "" {
		return "OFF"
	}
	return drive[value]
}

func convertMode(value int) string {
	mode := map[int]string{0: "FAN", 1: "COOL", 2: "HEAT", 3: "DRY", 4: "AUTO", 5: "BAHP", 6: "AUTOCOOL",
		7: "AUTOHEAT", 8: "VENTILATE", 9: "PANECOOL", 10: "PANEHEAT", 11: "OUTCOOL", 12: "DEFLOST",
		20: "SETBACK", 21: "SETBACKCOOL", 22: "SETBACKHEAT", 128: "HEATRECOVERY", 129: "BYPASS",
		130: "LC_AUTO", 144: "HEATING", 145: "HEATING_ECO", 146: "HOT_WATER", 147: "ANTI_FREEZE", 148: "COOLING"}
	return mode[value]
}

func convertFanSpeed(value int) string {
	speed := map[int]string{0: "LOW", 1: "MID2", 2: "MID1", 3: "HIGH", 6: "AUTO", 7: "EXLOW"}
	return speed[value]
}

func convertAirDirection(value int) string {
	direction := map[int]string{0: "SWING", 1: "VERTICAL", 2: "MID2", 3: "MID1", 4: "HORIZONTAL", 5: "MID0", 6: "AUTO"}
	return direction[value]
}

func convertTemp(value1, value2 int) float64 {
	return float64(value1) + (0.1 * float64(value2))
}

func convertTempX(value1, value2 int) float64 {
	value1 &= 0xFF
	value1 <<= 8
	value1 |= (value2 & 0xFF)
	return float64(value1) / 10.0
}

// not sure about this one
func convertUseTemp(value1 int) bool {
	if value1 == 0 {
		return true
	}
	return false
}
