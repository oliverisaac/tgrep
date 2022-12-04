package templating

import (
	"regexp"

	"github.com/pkg/errors"
)

var specialRegexCharacters = map[string]bool{
	`[`: true,
	`(`: true,
	`)`: true,
	`{`: true,
	`}`: true,
	`$`: true,
	`^`: true,
	`\`: true,
}

func Parse(input string) (string, error) {
	if templateContent, ok := templates[input]; ok {
		return templateContent, nil
	}

	var output string
	var thisChar, nextTwo string
	var templateContent string
	insideTemplate := false

	inputArr := []rune(input)

	for i := 0; i < len(inputArr); i++ {
		thisChar = string(inputArr[i : i+1])
		if i+2 <= len(inputArr) {
			nextTwo = string(inputArr[i : i+2])
		}

		if insideTemplate {
			if nextTwo == "}}" {
				i++
				insideTemplate = false
				if resolveTemplate, ok := templates[templateContent]; ok {
					output = output + resolveTemplate
				} else {
					output = output + templateContent
				}
				continue
			}
			templateContent = templateContent + thisChar
			continue
		}

		if nextTwo == "{{" {
			i++
			insideTemplate = true
			templateContent = ""
			continue
		}

		if nextTwo == `\{` {
			i++
			thisChar = `{`
		}

		if _, ok := specialRegexCharacters[thisChar]; ok {
			output = output + `\` + thisChar
		} else {
			output = output + thisChar
		}
	}

	if insideTemplate {
		return "", errors.New("Unclosed template {{")
	}

	return output, nil
}

func TemplatedRegexCaseInsensitive(input string) (*regexp.Regexp, error) {
	return TemplatedRegex("{{(?i)}}" + input)
}

func TemplatedRegex(input string) (*regexp.Regexp, error) {
	parsedInput, err := Parse(input)
	if err != nil {
		return nil, errors.Wrap(err, "Parsing input template")
	}
	r, err := regexp.Compile(parsedInput)
	if err != nil {
		return nil, errors.Wrapf(err, "Compiling parsed regex: %s", parsedInput)
	}
	return r, nil
}

func matchesRegex(stringToCheck, searchRegex string) (bool, error) {
	r, err := TemplatedRegex(searchRegex)
	if err != nil {
		return false, errors.Wrapf(err, "Getting tempalted regex of %s", searchRegex)
	}

	return r.MatchString(stringToCheck), nil
}

func Wrap(s string) string {
	return "{{" + s + "}}"
}
