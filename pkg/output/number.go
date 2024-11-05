package output

import (
	"fmt"
	"strings"
)

type Integer interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64
}

func FormatNumber[T Integer](num T) string {
	str := fmt.Sprintf("%d", num)
	n := len(str)

	var result strings.Builder

	for i := 0; i < n; i++ {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteString("_")
		}

		result.WriteByte(str[i])
	}

	return result.String()
}
