package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"

	"github.com/oliverisaac/tgrep/pkg/templating"
	"github.com/pkg/errors"
)

// Run is the primary entrypoint for the cli
func main() {
	err := runCommand(os.Args[1:], os.Stdout)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func runCommand(args []string, output io.Writer) error {
	var err error
	var parsedRegex string

	if len(args) == 0 {
		return errors.New(generateHelpMessage())
	}

	for i, arg := range args {
		if arg == "--help" {
			return errors.New(generateHelpMessage())
		}
		parsedRegex, err = templating.Parse(arg)
		if err != nil {
			return errors.Wrapf(err, "Parsing regex at index [%d]", i)
		}
		fmt.Fprintf(output, "%s\n", parsedRegex)
	}
	return nil
}

func generateHelpMessage() string {
	templates := templating.GetTemplates()
	templateKeys := []string{}
	for k := range templates {
		templateKeys = append(templateKeys, k)
	}
	sort.Strings(templateKeys)

	templateOutput := fmt.Sprintf("  %-12s  %s\n", "Template", "Regex")
	templateOutput = templateOutput + regexp.MustCompile("[^ \n]").ReplaceAllString(templateOutput, "=")
	for _, k := range templateKeys {
		templateOutput = templateOutput + fmt.Sprintf("- %-12s: %s\n", k, templates[k])
	}

	return `tg: Templated reGex

Usage: tg [REGEX_TEMPLATE... | REGEX_SHORTCUT...]

tg provides templating for regular expressions. The input argument is escaped to be safe to use with grep or other CLI tools. Wrap regular expressions in {{ }} and they will not be escaped.

Examples:
    tg '(hello {{(world|bob)}})' # The {{ }} will be treated as a regex, special characters will be auto-escaped
    tg email                     # This shorthand syntax will output the email regex
    tg 'My email is {{email}}'   # This will use the shorthand syntax inside a longer string

Templates supplied:
` + templateOutput
}
