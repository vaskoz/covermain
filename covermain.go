package main

import (
	"fmt"
	"io"
	"os"
	"text/template"
	"unicode"
	"unicode/utf8"
)

var (
	stderr io.Writer = os.Stderr
	stdout io.Writer = os.Stdout
	source           = sourceString
	tests            = testString
	mkdir            = func(dirname string) error {
		return os.Mkdir(dirname, os.ModePerm)
	}
	createFile = func(filename string) (io.Writer, error) {
		return os.Create(filename)
	}
)

type testName struct {
	Name string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(stderr, "Usage: %s CamelCaseNoSpacesName\n", os.Args[0])
		return
	}
	r, _ := utf8.DecodeRuneInString(os.Args[1])
	if !unicode.IsUpper(r) {
		fmt.Fprintf(stderr, "Usage: %s CamelCaseNoSpacesName (start with uppercase)\n", os.Args[0])
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
	dirname := camelcaseToSnakecase(os.Args[1])
	err = mkdir(dirname)
	if err != nil {
		fmt.Fprintln(stderr, "Could not create directory")
	}
	var sourceFile, testFile io.Writer
	sourceFile, err = createFile(fmt.Sprintf("%[1]s/%[1]s.go", dirname))
	if err != nil {
		fmt.Fprintln(stderr, "Can't create source file. Redirecting output to STDOUT")
		sourceFile = stdout
	}
	testFile, err = createFile(fmt.Sprintf("%[1]s/%[1]s_test.go", dirname))
	if err != nil {
		fmt.Fprintln(stderr, "Can't create test file")
		testFile = stdout
	}
	err = sourceTemplate.Execute(sourceFile, testName{os.Args[1]})
	if err != nil {
		fmt.Fprintln(stderr, "Couldn't write to source file")
	}
	err = testTemplate.Execute(testFile, testName{os.Args[1]})
	if err != nil {
		fmt.Fprintln(stderr, "Couldn't write to test file")
	}
}

func camelcaseToSnakecase(camel string) string {
	if camel == "" {
		return ""
	}
	var result []rune
	for _, r := range camel {
		if unicode.IsUpper(r) {
			result = append(result, '_')
			r = unicode.ToLower(r)
		}
		result = append(result, r)
	}
	return string(result[1:])
}
