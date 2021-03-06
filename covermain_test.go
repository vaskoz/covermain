package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

var ccTestcases = []struct {
	in, out string
}{
	{"AGreatProject", "a_great_project"},
	{"SomethingElseGreat", "something_else_great"},
	{"HttpProject", "http_project"},
	{"SomeHttpProject", "some_http_project"},
	{"", ""},
}

func TestCamelcaseToSnakecase(t *testing.T) {
	t.Parallel()
	for _, c := range ccTestcases {
		if out := camelcaseToSnakecase(c.in); out != c.out {
			t.Errorf("%s does not equal expected %s", out, c.out)
		}
	}
}

func BenchmarkCamelcaseToSnakecase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, c := range ccTestcases {
			camelcaseToSnakecase(c.in)
		}
	}
}

// Can't be run in parallel due to global state mutation
func TestCoverMainBadName(t *testing.T) {
	buff := new(bytes.Buffer)
	stderr = buff
	os.Args = append(os.Args[:1], "lower")
	main()
	out := buff.String()
	expected := "start with uppercase"
	if !strings.Contains(out, expected) {
		t.Errorf("Expected a different usage message for a bad name")
	}
}

// Can't be run in parallel due to global state mutation
func TestCoverMainBadSource(t *testing.T) {
	source = `{{ nope }}`
	buff := new(bytes.Buffer)
	stderr = buff
	os.Args = append(os.Args[:1], "Upper")
	main()
	out := buff.String()
	expected := "Can't parse source file template\n"
	if out != expected {
		t.Errorf("%s does not equal expected %s", out, expected)
	}
	source = sourceString
}

// Can't be run in parallel due to global state mutation
func TestCoverMainBadTests(t *testing.T) {
	tests = `{{ nope }}`
	buff := new(bytes.Buffer)
	stderr = buff
	os.Args = append(os.Args[:1], "Upper")
	main()
	out := buff.String()
	expected := "Can't parse test file template\n"
	if out != expected {
		t.Errorf("%s does not equal expected %s", out, expected)
	}
	tests = testString
}

// Can't be run in parallel due to global state mutation
func TestCoverMainTooManyArguments(t *testing.T) {
	os.Args = append(os.Args, "one too many")
	buff := new(bytes.Buffer)
	stderr = buff
	main()
	out := buff.String()
	if !strings.Contains(out, "CamelCaseNoSpacesName") {
		t.Errorf("Too many args should result in usage message")
	}
}

// Can't be run in parallel due to global state mutation
func TestCoverMainTooFewArguments(t *testing.T) {
	os.Args = os.Args[:1]
	buff := new(bytes.Buffer)
	stderr = buff
	main()
	out := buff.String()
	if !strings.Contains(out, "CamelCaseNoSpacesName") {
		t.Errorf("Too few args should result in usage message")
	}
}

// Can't be run in parallel due to global state mutation
func TestCoverMainCantCreateDirectory(t *testing.T) {
	os.Args = append(os.Args[:1], "CantCreateSourceFile")
	buff := new(bytes.Buffer)
	originalMkdir := mkdir
	mkdir = func(dirname string) error {
		return errors.New("can't create directory")
	}
	originalCreateFile := createFile
	createFile = func(filename string) (io.Writer, error) {
		return stdout, nil
	}
	stderr = buff
	main()
	out := buff.String()
	if !strings.Contains(out, "Could not create directory") {
		t.Errorf("Can't create directory")
	}
	createFile = originalCreateFile
	mkdir = originalMkdir
}

// Can't be run in parallel due to global state mutation
func TestCoverMainCantCreateSourceFile(t *testing.T) {
	os.Args = append(os.Args[:1], "CantCreateSourceFile")
	buff := new(bytes.Buffer)
	originalMkdir := mkdir
	mkdir = func(dirname string) error {
		return nil
	}
	originalCreateFile := createFile
	createFile = func(filename string) (io.Writer, error) {
		if !strings.Contains(filename, "_test.go") {
			return nil, errors.New("Can't create source file")
		}
		return stdout, nil
	}
	stderr = buff
	main()
	out := buff.String()
	if !strings.Contains(out, "Can't create source file. Redirecting output to STDOUT") {
		t.Errorf("Can't create source file")
	}
	createFile = originalCreateFile
	mkdir = originalMkdir
}

// Can't be run in parallel due to global state mutation
func TestCoverMainCantCreateTestFile(t *testing.T) {
	os.Args = append(os.Args[:1], "CantCreateTestFile")
	buff := new(bytes.Buffer)
	originalMkdir := mkdir
	mkdir = func(dirname string) error {
		return nil
	}
	originalCreateFile := createFile
	createFile = func(filename string) (io.Writer, error) {
		if strings.Contains(filename, "_test.go") {
			return nil, errors.New("Can't create test file")
		}
		return stdout, nil
	}
	stderr = buff
	main()
	out := buff.String()
	if !strings.Contains(out, "Can't create test file") {
		t.Errorf("Can't create test file")
	}
	createFile = originalCreateFile
	mkdir = originalMkdir
}

// Can't be run in parallel due to global state mutation
func TestCoverMainCantWriteTemplates(t *testing.T) {
	os.Args = append(os.Args[:1], "CantWriteTemplates")
	buff := new(bytes.Buffer)
	originalCreateFile := createFile
	originalMkdir := mkdir
	mkdir = func(dirname string) error {
		return nil
	}
	createFile = func(filename string) (io.Writer, error) {
		return nil, nil
	}
	originalSource := source
	source = "{{.Foobar}}"
	originalTest := tests
	tests = "{{.Foobar}}"
	stderr = buff
	main()
	out := buff.String()
	if !strings.Contains(out, "Couldn't write to source file") {
		t.Errorf("Can't parse source file template")
	}
	if !strings.Contains(out, "Couldn't write to test file") {
		t.Errorf("Can't parse test file template")
	}
	createFile = originalCreateFile
	mkdir = originalMkdir
	source = originalSource
	tests = originalTest
}

var testcases = []struct {
	name string
}{
	{"AGreatProject"},
	{"SomeProject"},
	{"HttpProject"},
}

func TestCoverMain(t *testing.T) {
	originalCreateFile := createFile
	originalMkdir := mkdir
	createFile = func(filename string) (io.Writer, error) {
		return stderr, nil
	}
	mkdir = func(dirname string) error {
		return nil
	}
	for _, c := range testcases {
		buff := new(bytes.Buffer)
		stderr = buff
		os.Args = append(os.Args[:1], c.name)
		main()
		out := buff.String()
		if !strings.Contains(out, fmt.Sprintf("Test%s", c.name)) {
			t.Errorf("%s does not contain Test%s", out, c.name)
		}
		if !strings.Contains(out, fmt.Sprintf("Benchmark%s", c.name)) {
			t.Errorf("%s does not contain Benchmark%s", out, c.name)
		}
	}
	mkdir = originalMkdir
	createFile = originalCreateFile
}

func BenchmarkCoverMain(b *testing.B) {
	originalCreateFile := createFile
	originalMkdir := mkdir
	createFile = func(filename string) (io.Writer, error) {
		return stderr, nil
	}
	mkdir = func(dirname string) error {
		return nil
	}
	for i := 0; i < b.N; i++ {
		for _, c := range testcases {
			buff := new(bytes.Buffer)
			stderr = buff
			os.Args = append(os.Args[:1], c.name)
			main()
		}
	}
	mkdir = originalMkdir
	createFile = originalCreateFile
}

func TestCoverMainIntegration(t *testing.T) {
	for _, c := range testcases {
		buff := new(bytes.Buffer)
		stderr = buff
		os.Args = append(os.Args[:1], c.name)
		main()
		out := buff.String()
		if out != "" {
			t.Errorf("STDERR should be empty on successful run")
		}
		dirname := camelcaseToSnakecase(c.name)
		filename := fmt.Sprintf("%[1]s/%[1]s.go", dirname)
		filenameTest := fmt.Sprintf("%[1]s/%[1]s_test.go", dirname)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Errorf("Source file should exist")
		}
		if _, err := os.Stat(filenameTest); os.IsNotExist(err) {
			t.Errorf("Test file should exist")
		}
		err := os.RemoveAll(dirname)
		if err != nil {
			t.Errorf("Couldn't remove directory used for integration testing")
		}
	}
}
