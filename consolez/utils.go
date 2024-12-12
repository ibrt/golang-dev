package consolez

import (
	"fmt"
	"strings"
)

func alignRight(s string, width int) string {
	if width < 4 {
		width = 4
	}

	r := []rune(s)

	if len(r) > width {
		return fmt.Sprintf("...%v", string(r[len(r)-width+3:]))
	}

	if len(r) < width {
		return strings.Repeat(".", width-len(r)) + string(r)
	}

	return string(r)
}

func truncateLeft(s string, maxWidth int) string {
	if maxWidth < 4 {
		maxWidth = 4
	}

	if r := []rune(s); len(r) > maxWidth {
		return fmt.Sprintf("...%v", string(r[len(r)-maxWidth+3:]))
	}

	return s
}
