package utils

import (
	"bytes"
	"strings"
)

func MaskData(data string) string {
	length := len(data)

	if length > 4 {
		var buffer bytes.Buffer

		// Write stars
		buffer.WriteString(strings.Repeat("*", length-4))

		// Write last 4 characters
		buffer.WriteString(data[length-5 : length-1])

		return buffer.String()
	}

	return data
}

func FormatCurrencyFloat32(amount int32) float32 {
	return float32(amount) / 100
}

func FormatCurrencyFloat64(amount int32) float64 {
	return float64(amount) / 100
}

func FormatCurrencyInt(amount float32) int32 {
	return int32(amount * 100)
}
