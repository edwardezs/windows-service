package process

import (
	"os"
	"os/signal"
	"strings"
)

var dev_null, _ = os.Open(os.DevNull)

// DevNull returns the NULL file and can be used to suppress
// input and output on processes. It can be used as a regular
// os.File
func DevNull() *os.File {
	res := new(os.File)
	*res = *dev_null
	return res
}

func ListenForCTRLC() chan os.Signal {
	exitC := make(chan os.Signal, 10)
	signal.Notify(exitC, os.Interrupt)
	return exitC
}

// ParseCommandArgs gets a list of strings and parses their content
// splitting them into separated strings. The characters used to parse
// the commands are, in the relevant order, <'>, <"> and < >
func ParseCommandArgs(args ...string) []string {
	a := make([]string, 0)
	for _, s := range args {
		for i, s1 := range strings.Split(s, "'") {
			if i%2 == 1 {
				a = append(a, s1)
				continue
			}

			for j, s2 := range strings.Split(s1, "\"") {
				if j%2 == 1 {
					a = append(a, s2)
					continue
				}

				for _, s3 := range strings.Split(s2, " ") {
					if s3 != "" {
						a = append(a, s3)
					}
				}
			}
		}
	}

	

	return a
}

// FastCommandParse
func FastCommandParse(args ...string) (wd string, execName string, argV []string) {
	a := ParseCommandArgs(args...)
	return "", a[0], a[1:]
}
