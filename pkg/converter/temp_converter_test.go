package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCelsiusToFahrenheitIsValid(t *testing.T) {
	celsius := 0.0
	expected := 32.0
	result := CelsiusToFahrenheit(celsius)
	assert.Equal(t, expected, result)
}

func TestCelsiusToKelvin(t *testing.T) {
	celsius := 0.0
	expected := 273.15
	result := CelsiusToKelvin(celsius)
	assert.Equal(t, expected, result)
}
