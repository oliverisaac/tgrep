package main

import (
	"fmt"
	"os"

	"github.com/oliverisaac/tgrep/tgrep"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	var err error
	var parsedRegex string

	for i, arg := range os.Args[1:] {
		parsedRegex, err = tgrep.Parse(arg)
		if err != nil {
			logrus.Fatal(errors.Wrapf(err, "Parsing regex at index [%d]", i))
		}
		fmt.Println(parsedRegex)
	}
}
