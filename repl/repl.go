package repl

import (
	"strings"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func() error
}

func CleanInput(text string) []string {
	text = strings.ToLower(text)

	return strings.Split(text, " ")
}
