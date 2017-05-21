package main

import (
	"fmt"
	"io"
	"os"
	"text/template"
	"unicode"
)

var stderr io.Writer = os.Stderr
var source string = sourceString
var tests string = testString

type TestName struct {
	Name string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(stderr, "Usage: %s CamelCaseNoSpacesName\n", os.Args[0])
		return
	}
	var err error
	sourceTemplate := template.New("source")
	sourceTemplate, err = sourceTemplate.Parse(source)
	if err != nil {
		fmt.Fprintln(stderr, "Can't parse source file template")
		return
	}

	testTemplate := template.New("test")
	testTemplate, err = testTemplate.Parse(tests)
	if err != nil {
		fmt.Fprintln(stderr, "Can't parse test file template")
		return
	}
	dirname := CamelcaseToLowercase(os.Args[1])
	os.Mkdir(dirname, os.ModePerm)
	sourceFile, err := os.Create(fmt.Sprintf("%[1]s/%[1]s.go", dirname))
	testFile, err := os.Create(fmt.Sprintf("%[1]s/%[1]s_test.go", dirname))
	sourceTemplate.Execute(sourceFile, TestName{os.Args[1]})
	testTemplate.Execute(testFile, TestName{os.Args[1]})
}

func CamelcaseToLowercase(camel string) string {
	var result []rune
	for i, r := range camel {
		if i != 0 && unicode.IsUpper(r) {
			result = append(result, '_')
		}
		r := unicode.ToLower(r)
		result = append(result, r)
	}
	return string(result)
}
