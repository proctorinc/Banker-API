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
