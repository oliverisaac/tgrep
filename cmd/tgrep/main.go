package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"

	arg "github.com/alexflint/go-arg"

	"github.com/oliverisaac/tgrep/pkg/templating"
	"github.com/pkg/errors"
)

type config struct {
	CaseInsensitive bool     `arg:"-i"`
	Regex           []string `arg:"-e,separate"`
	PositionalArgs  []string `arg:"positional"`
}

// Run is the primary entrypoint for the cli
func main() {
	err := runCommand(os.Args[1:], os.Stdin, os.Stdout)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func runCommand(args []string, input io.Reader, output io.Writer) error {
	config := config{}
	argParser, err := arg.NewParser(arg.Config{
		Program: "tgrep",
	}, &config)
	if err != nil {
		return errors.Wrap(err, "Creating arg parser")
	}

	err = argParser.Parse(args)
	if len(args) == 0 || err == arg.ErrHelp {
		return errors.New(generateHelpMessage(argParser))
	} else if err != nil {
		return errors.Wrap(err, "Failed to parse args")
	}

	// If the regex flag is not used, then use the first argument as the regex
	if len(config.Regex) == 0 {
		if len(config.PositionalArgs) == 0 {
			return errors.New(generateHelpMessage(argParser))
		}
		config.Regex = []string{config.PositionalArgs[0]}
		config.PositionalArgs = config.PositionalArgs[1:]
	}

	regexesToUse := []*regexp.Regexp{}
	for _, r := range config.Regex {
		var reg *regexp.Regexp
		var err error
		if config.CaseInsensitive {
			reg, err = templating.TemplatedRegexCaseInsensitive(r)
		} else {
			reg, err = templating.TemplatedRegex(r)
		}
		if err != nil {
			return errors.Wrap(err, "Templating regex")
		}
		regexesToUse = append(regexesToUse, reg)
	}

	if len(config.PositionalArgs) == 0 {
		err := runRegularExpressionsAgainstReader(regexesToUse, input, output)
		if err != nil {
			return errors.Wrap(err, "Running regex against stdin")
		}
	} else {
		for _, fname := range config.PositionalArgs {
			f, err := os.Open(fname)
			if err != nil {
				return errors.Wrap(err, "Failed to open file")
			}
			err = runRegularExpressionsAgainstReader(regexesToUse, f, output)
			if err != nil {
				return errors.Wrap(err, "Running regex against stdin")
			}
		}
	}

	return nil
}

func runRegularExpressionsAgainstReader(regexesToUse []*regexp.Regexp, read io.Reader, write io.Writer) error {

	fileScanner := bufio.NewScanner(read)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		for _, r := range regexesToUse {
			if r.MatchString(fileScanner.Text()) {
				fmt.Fprintln(write, fileScanner.Text())
				break
			}
		}
	}

	return nil
}

func generateHelpMessage(argParser *arg.Parser) string {
	helpOutput := &bytes.Buffer{}
	argParser.WriteHelp(helpOutput)

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

	return string(helpOutput.Bytes()) + `

tgrep provides templating for regular expressions. The input argument is escaped to be safe to use with grep or other CLI tools. Wrap regular expressions in {{ }} and they will not be escaped.

Examples:
    tgrep '(hello {{(world|bob)}})'                     # The {{ }} will be treated as a regex, special characters will be auto-escaped
    tgrep -e email /tmp/example.txt                     # This shorthand syntax will output the email regex
    tgrep 'My email is {{email}}' /tmp/example.txt      # This will use the shorthand syntax inside a longer string

Templates supplied:
` + templateOutput
}
