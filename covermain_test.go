package main

import (
	"bytes"
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

func TestCamelcaseToLowercase(t *testing.T) {
	t.Parallel()
	for _, c := range ccTestcases {
		if out := CamelcaseToLowercase(c.in); out != c.out {
			t.Errorf("%s does not equal expected %s", out, c.out)
		}
	}
}

func BenchmarkCamelcaseToLowercase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, c := range ccTestcases {
			CamelcaseToLowercase(c.in)
		}
	}
}

// Can't be run in parallel due to global state mutation
func TestCoverMainBadSource(t *testing.T) {
	source = `{{ nope }}`
	buff := new(bytes.Buffer)
	stderr = buff
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

func BenchmarkCoverMain(b *testing.B) {
}
