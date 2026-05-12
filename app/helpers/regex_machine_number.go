package helpers

import (
	"regexp"
	"strings"
)

var MachineNumberRegex = regexp.MustCompile(`^[A-Z0-9]{8,25}$`)

func NormalizeMachineNumber(machine string) string {
	machine = regexp.MustCompile(`\s+`).ReplaceAllString(machine, "")
	machine = strings.ToUpper(machine)
	return machine
}
