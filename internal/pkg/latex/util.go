package latex

import (
	"strings"
)

var oldNew = []string{
	"\\", "\\\\",
	"$", "\\$",
	"[", "\\[",
	"]", "\\]",
	"{", "\\{",
	"}", "\\}",
	"%", "\\%",
	"_", "\\_",
	"^", "\\^",
	"#", "\\#",
}

func escapeText(text string) string {
	return strings.NewReplacer(oldNew...).Replace(text)
}