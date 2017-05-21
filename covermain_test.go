package main

import "testing"

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
