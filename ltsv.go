package loghttpltsv

import (
	"fmt"
	"strings"
)

func replaceColon(s string) string {
	return strings.ReplaceAll(s, ":", ";")
}

func replaceTab(s string) string {
	return strings.ReplaceAll(s, "\t", " ")
}

func replaceNewLine(s string) string {
	return strings.ReplaceAll(s, "\n", " ")
}

func emptyToHyphen(s string) string {
	if s == "" {
		return "-"
	}

	return s
}

func lv(label string, value interface{}) string {
	l := replaceColon(label)
	v := emptyToHyphen(fmt.Sprint(value))
	return replaceNewLine(replaceTab(l + ":" + v))
}

func ltsv(lvs []string) string {
	return strings.Join(lvs, "\t")
}
