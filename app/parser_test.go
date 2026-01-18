package main

import (
	"testing"
)

type TestCase struct {
	value  []string
	length int
}

func NewTestCase(length int, value []string) TestCase {
	return TestCase{
		value:  value,
		length: length,
	}
}

func TestSplitWithNormalString(t *testing.T) {

	input := `Hello World`
	want := []string{`Hello`, `World`}

	got, err := Split(input)

	if err != nil {
		t.Error(err)
	}

	if len(got) != len(want) {
		t.Errorf("Split(%q). Got=[%v], Want=[%v]", input, got, want)
	}

	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Split(%q). Got=[%v], Want=[%v]", input, got[i], want[i])
		}
	}

}

func TestSplitWithSingleQuoteString(t *testing.T) {

	input := []string{`Hello World`, `'Hello     World'`, `Hello''World`, `Hello    World`, `'"Hello      World"'`}
	want := []TestCase{
		NewTestCase(2, []string{`Hello`, `World`}),
		NewTestCase(1, []string{`Hello     World`}),
		NewTestCase(1, []string{`HelloWorld`}),
		NewTestCase(2, []string{`Hello`, `World`}),
		NewTestCase(1, []string{`"Hello      World"`}),
	}

	for tt := range input {

		got, err := Split(input[tt])
		if err != nil {
			t.Error(err)
		}

		if len(got) != want[tt].length {
			t.Errorf("Split(%q),  Got: %v, Want: %v", input[tt], got, want[tt])
		}

		for tti := range got {
			if want[tt].value[tti] != got[tti] {
				t.Errorf("Split(%q),  Got: %v, Want: %v", input[tt], got, want[tt])
			}
		}

	}

}

func TestSplitWithDoubleQuoteString(t *testing.T) {

	input := []string{`"Hello     World"`, `Hello""World`, `"Hello"    "World"`, `"'Hello      World'"`, `"Hello""World"`}
	want := []TestCase{
		NewTestCase(1, []string{`Hello     World`}),
		NewTestCase(1, []string{`HelloWorld`}),
		NewTestCase(2, []string{`Hello`, `World`}),
		NewTestCase(1, []string{`'Hello      World'`}),
		NewTestCase(1, []string{`HelloWorld`}),
	}

	for tt := range input {

		got, err := Split(input[tt])
		if err != nil {
			t.Error(err)
		}

		if len(got) != want[tt].length {
			t.Errorf("Split(%q),  Got: %v, Want: %v", input[tt], got, want[tt])
		}

		for tti := range got {
			if want[tt].value[tti] != got[tti] {
				t.Errorf("Split(%q),  Got: %v, Want: %v", input[tt], got, want[tt])
			}
		}

	}

}

func TestSplitWithEscapeCharacterString(t *testing.T) {

	input := []string{
		`world\ \ \ \ \ \ script`,
		`"shell'hello'\\'example"`,
		`"/tmp/ant/\"f 38\"" "/tmp/ant/\"f\\93\""`,
		`"/tmp/ant/'f 17'" "/tmp/ant/'f  \\39'" "/tmp/ant/'f \\16\\'"`,
	}
	want := []TestCase{
		NewTestCase(1, []string{`world      script`}),
		NewTestCase(1, []string{`shell'hello'\'example`}),
		NewTestCase(2, []string{`/tmp/ant/"f 38"`, `/tmp/ant/"f\93"`}),
		NewTestCase(3, []string{`/tmp/ant/'f 17'`, `/tmp/ant/'f  \39'`, `/tmp/ant/'f \16\'`}),
	}

	for tt := range input {

		got, err := Split(input[tt])
		if err != nil {
			t.Error(err)
		}

		if len(got) != want[tt].length {
			t.Errorf("Split(%q),  Got: %v, Want: %v", input[tt], got, want[tt])
		}

		for tti := range got {
			if want[tt].value[tti] != got[tti] {
				t.Errorf("Split(%q),  Got: %v, Want: %v", input[tt], got, want[tt])
			}
		}

	}

}
