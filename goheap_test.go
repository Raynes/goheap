package goheap

import (
	"testing"
	"os"
)

func devConfig() (config Config) {
	url   := os.Getenv("RH_URL")
	user  := os.Getenv("RH_USER")
	token := os.Getenv("RH_TOKEN")
	if url == "" {
		config.URL = RefheapURL
	}
	config.User = user
	config.Key  = token
	return
}

func cError(t *testing.T, config *Config, expected interface{}, err *error, call string) {
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
		cError(t, &config, &zero, &err, "NewConfig()")
	}

	if config, err := NewConfig("foo"); err != nil || config != one {
		cError(t, &config, &one, &err, "NewConfig(\"foo\")")
	}

	if config, err := NewConfig("raynes", "123"); err != nil || config != two {
		cError(t, &config, &two, &err, "NewConfig(\"raynes\", \"123\")")
	}

	if config, err := NewConfig("foo", "raynes", "123"); err != nil || config != three {
		cError(t, &config, &three, &err, "NewConfig(\"foo\", \"raynes\", \"123\", )")
	}

	if config, err := NewConfig("", "", "", ""); err == nil {
		cError(t, &config, &error, &err, "NewConfig(\"\", \"\", \"\", \"\")")
	}
}

// This will be set to whatever the current expression is for
// gpError() messages. It is a convenience because validating
// individual paste fields manually is already tedious and
// passing the current expression each time would be a massive
// pain in the rear. It pokes at my FP nerves, but these are
// merely tests after all. We're allowed a bit of leeway. When
// changing this variable we should always document what we're
// doing with a comment.
var expression string

func gpError(t *testing.T, missing string, missingValue interface{}, expected interface{}) {
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
	config := devConfig()
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
		gpError(t, "Lines", lines, 1)
	}

	if views := paste.Views; views <= 0 {
		gpError(t, "Views", views, "a number greater than zero")
	}

	const dateValue = "2012-01-04T01:44:22.964Z"
	if date := paste.Date; date != dateValue {
		gpError(t, "Date", date, dateValue)
	}

	if ID := paste.ID; ID != "1" {
		gpError(t, "ID", ID, "1")
	}

	if language := paste.Language; language != "Clojure" {
		gpError(t, "Language", language, "Clojure")
	}

	if private := paste.Private; private {
		gpError(t, "Private", private, !private)
	}

	const expectedUrl = "https://www.refheap.com/1"
	if url := paste.URL; url != expectedUrl {
		gpError(t, "Url", url, expectedUrl)
	}

	if user := paste.User; user != "raynes" {
		gpError(t, "User", user, "raynes")
	}

	if contents := paste.Contents; contents != "(begin)" {
		gpError(t, "Contents", contents, "(begin)")
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

func deletePaste(config *Config, paste *Paste) {
	DeletePaste(config, paste.ID)
}

// Sadly, TestCreatePaste and TestDeletePaste are rather interleaved, since we
// can't delete a paste without creating it (and thus TestCreatePaste must
// pass) and you don't want to create a paste without deleting it after
// because nobody likes a litterbug. As such, these tests depend on one
// another.

func TestCreatePaste(t *testing.T) {
	config := devConfig()
	expression = "CreatePaste(&config, Paste{Private: true, Contents: \"hi\", Language: \"Go\"})"
	paste := Paste{Private: true, Contents: "hi", Language: "Go"}
	defer deletePaste(&config, &paste)
	err := CreatePaste(&config, &paste)
	if err != nil {
		t.Errorf("Error creating paste with expression %v: %v", expression, err)
	}

	if pUser, cUser := paste.User, config.User; pUser != cUser {
		t.Errorf("Expected creating user to be %v. It was %v.", cUser, pUser)
	}

	if lang := paste.Language; lang != "Go" {
		t.Errorf("Expected language to be Go. It was %v.", lang)
	}

	if priv := paste.Private; !priv {
		t.Error("Expected paste to be private!")
	}
}

func TestDeletePaste(t *testing.T) {
	// TODO: Implement.
}
