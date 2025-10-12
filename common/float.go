package common

import (
	"fmt"
	"strconv"
)

func ToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

// StringToFloat converts a string to float64, returns the converted value and any error that occurred
func StringToFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func FloatToString(value float64) string {
	return fmt.Sprintf("%f", value)
}
