package templating

import (
	"github.com/pkg/errors"
	"regexp"
)

var specialRegexCharacters = map[string]bool{
	`[`: true,
	`(`: true,
	`)`: true,
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

		if nextTwo == "{{" {
			i++
			insideTemplate = true
			templateContent = ""
			continue
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

		if _, ok := specialRegexCharacters[thisChar]; ok {
			output = output + `\` + thisChar
		} else {
			output = output + thisChar
		}
	}

	return output, nil
}

func templatedRegex(input string) (*regexp.Regexp, error) {
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
	r, err := templatedRegex(searchRegex)
	if err != nil {
		return false, errors.Wrapf(err, "Getting tempalted regex of %s", searchRegex)
	}

	return r.MatchString(stringToCheck), nil
}
