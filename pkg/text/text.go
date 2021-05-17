package text

import (
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var (
	bold  = color.New(color.Bold)
	green = color.New(color.FgGreen)
	red   = color.New(color.FgRed)
)

func Bold(a ...interface{}) string  { return bold.Sprint(a...) }
func Green(a ...interface{}) string { return green.Sprint(a...) }
func Red(a ...interface{}) string   { return red.Sprint(a...) }
func GreenF(format string, a ...interface{}) string {
	return green.Sprintf(format, a...)
}
func RedF(format string, a ...interface{}) string {
	return red.Sprintf(format, a...)
}

var lineRE = regexp.MustCompile(`(?m)^`)

func Indent(s, indent string) string {
	if strings.TrimSpace(s) == "" {
		return s
	}
	return lineRE.ReplaceAllLiteralString(s, indent)
}
