package goheap

import "testing"

func CError(t *testing.T, config *Config, expected interface{}, err *error, call string) {
	t.Errorf("%v failed! Returned config %#v and err %#v; Wanted %#v",
		call, config, err, expected)
}

// This function is by nature pretty fickle because of the magic
// that it does with variadic arguments. As such, we're going to
// very thoroughly test it!
func TestNewConfig(t *testing.T) {
	zero := Config{RefheapURL, "", ""}
	one := Config{"foo", "", ""}
	two := Config{RefheapURL, "raynes", "123"}
	three := Config{"foo", "raynes", "123"}
	error := ConfigError{[]string{"", "", "", ""}}

	if config, err := NewConfig(); err != nil || config != zero {
		CError(t, &config, &zero, &err, "NewConfig()")
	}

	if config, err := NewConfig("foo"); err != nil || config != one {
		CError(t, &config, &one, &err, "NewConfig(\"foo\")")
	}

	if config, err := NewConfig("raynes", "123"); err != nil || config != two {
		CError(t, &config, &two, &err, "NewConfig(\"raynes\", \"123\")")
	}

	if config, err := NewConfig("foo", "raynes", "123"); err != nil || config != three {
		CError(t, &config, &three, &err, "NewConfig(\"foo\", \"raynes\", \"123\", )")
	}

	if config, err := NewConfig("", "", "", ""); err == nil {
		CError(t, &config, &error, &err, "NewConfig(\"\", \"\", \"\", \"\")")
	}
}

// This will be set to whatever the current expression is for
// GPError() messages. It is a convenience because validating
// individual paste fields manually is already tedious and
// passing the current expression each time would be a massive
// pain in the rear. It pokes at my FP nerves, but these are
// merely tests after all. We're allowed a bit of leeway. When
// changing this variable we should always document what we're
// doing with a comment.
var expression string

func GPError(t *testing.T, missing string, missingValue interface{}, expected interface{}) {
	err := `
		Expression %v failing.
		Paste field %v was not as expected.
		Got %#v; Expected %v.
		`
	t.Errorf(err, expression, missing, missingValue, expected)
}

// TODO: Allow for test configuration for calls like this with
// environment variables to set refheap url, user, pass, etc.
func TestGetPaste(t *testing.T) {
	// Set what the current expression is for error messages.
	expression = "GetPaste(&config, \"1\")"
	config, _ := NewConfig()
	paste, err := GetPaste(&config, "1")
	if err != nil {
		t.Errorf("%v failed because of error %v", expression, err)
		return
	}

	// Unfortunately we cannot just create a dummy object to
	// compare against because views is dynamic. Technically
	// all of this is dynamic, but views is the only thing
	// a person other than me (Raynes) can change. Anyways,
	// because of this we have to validate each field one by
	// one manually. At least we get nice failure messages.
	if lines := paste.Lines; lines != 1 {
		GPError(t, "Lines", lines, 1)
	}

	if views := paste.Views; views <= 0 {
		GPError(t, "Views", views, "a number greater than zero")
	}

	const dateValue = "2012-01-04T01:44:22.964Z"
	if date := paste.Date; date != dateValue {
		GPError(t, "Date", date, dateValue)
	}

	if pasteID := paste.PasteID; pasteID != "1" {
		GPError(t, "PasteID", pasteID, "1")
	}

	if language := paste.Language; language != "Clojure" {
		GPError(t, "Language", language, "Clojure")
	}

	if private := paste.Private; private {
		GPError(t, "Private", private, !private)
	}

	const expectedUrl = "https://www.refheap.com/1"
	if url := paste.Url; url != expectedUrl {
		GPError(t, "Url", url, expectedUrl)
	}

	if user := paste.User; user != "raynes" {
		GPError(t, "User", user, "raynes")
	}

	if contents := paste.Contents; contents != "(begin)" {
		GPError(t, "Contents", contents, "(begin)")
	}

	// Set expression for next round of tests.
	expression = "GetPaste(&config, \"0\")"
	expectedErr := RefheapError{"Paste does not exist."}
	paste, err = GetPaste(&config, "@#R##")
	if err != expectedErr {
		msg := `
		Expression %v did not fail as expected.
		err was %#v.
		Expected err to be %#v.
		`
		t.Errorf(msg, expression, err, expectedErr)
	}
}

