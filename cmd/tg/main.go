package main

import (
	"fmt"
	"github.com/oliverisaac/tgrep/pkg/templating"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	var parsedRegex string

	for i, arg := range os.Args[1:] {
		parsedRegex, err = templating.Parse(arg)
		if err != nil {
			logrus.Fatal(errors.Wrapf(err, "Parsing regex at index [%d]", i))
		}
		fmt.Println(parsedRegex)
	}
}
