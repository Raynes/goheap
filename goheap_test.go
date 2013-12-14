package goheap

import (
	"os"
	"strings"
	"testing"
)

func devConfig() (config Config) {
	url := os.Getenv("RH_URL")
	user := os.Getenv("RH_USER")
	token := os.Getenv("RH_TOKEN")
	if url == "" {
		config.URL = RefheapURL
	}
	config.User = user
	config.Key = token
	return
}

func cError(t *testing.T, config *Config, expected interface{}, err *error, call string) {
	msg := `
	%v failed!
	Expected %#v.
	Got err %#v and config %#v
	`
	t.Errorf(msg, call, config, err, expected)
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

func TestCreate(t *testing.T) {
	config := devConfig()
	paste := Paste{Private: true, Contents: "hi", Language: "Go"}
	err := paste.Create(&config)
	if err != nil {
		t.Errorf("Error creating paste: %v", err)
	}

	defer paste.Delete(&config)

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

func TestDelete(t *testing.T) {
	config := devConfig()
	paste := Paste{Contents: "foo", Private: true}
	if err := paste.Create(&config); err != nil {
		t.Errorf("Something went wrong creating a paste: %v", err)
	}

	if err := paste.Delete(&config); err != nil {
		t.Errorf("Something went wrong deleting a paste: %v", err)
	}

	err := paste.Get(&config)
	if _, ok := err.(RefheapError); !ok {
		t.Errorf("Paste %v still exists after trying to delete!", paste.ID)
	}
}

func gpError(t *testing.T, missing string, missingValue interface{}, expected interface{}) {
	err := `
		Paste field %v was not as expected.
		Got %#v; Expected %v.
		`
	t.Errorf(err, missing, missingValue, expected)
}

func TestGet(t *testing.T) {
	config := devConfig()
	testPaste := Paste{Private: true, Contents: "hi", Language: "Go"}
	defer testPaste.Delete(&config)
	if err := testPaste.Create(&config); err != nil {
		t.Errorf("Something went wrong creating a paste: %v", err)
	}

	paste := Paste{ID: testPaste.ID}
	if err := paste.Get(&config); err != nil {
		t.Errorf("Something went wrong getting a paste: %v", err)
	}

	if lines := paste.Lines; lines != 1 {
		gpError(t, "Lines", lines, 1)
	}

	if date := paste.Date; date == "" {
		gpError(t, "Date", "a date", "no date")
	}

	if id := paste.ID; id == "" {
		gpError(t, "ID", "no id", "an id")
	}

	if language := paste.Language; language != "Go" {
		gpError(t, "Language", language, "Go")
	}

	if private := paste.Private; !private {
		gpError(t, "Private", !private, private)
	}

	if url := paste.URL; url == "" {
		gpError(t, "Url", url, "no url")
	}

	if user := paste.User; user != config.User {
		gpError(t, "User", user, config.User)
	}

	if contents := paste.Contents; contents != "hi" {
		gpError(t, "Contents", contents, "hi")
	}

	expectedErr := RefheapError{"Paste does not exist."}
	paste = Paste{ID: "@D("}
	err := paste.Get(&config)
	if err != expectedErr {
		msg := `
		err was %#v.
		Expected err to be %#v.
		`
		t.Errorf(msg, err, expectedErr)
	}
}

func TestGetHighlighted(t *testing.T) {
	config := devConfig()
	paste := Paste{Private: true, Contents: "hi"}

	if err := paste.Create(&config); err != nil {
		t.Errorf("Something went wrong saving a paste: %v", err)
	}

	defer paste.Delete(&config)

	highlighted, err := paste.GetHighlighted(&config)
	if err != nil {
		return
	}

	if !strings.HasPrefix(highlighted.Content, "<table") {
		t.Errorf("Expected string to begin with '<table'. Got: %#v", highlighted.Content)
	}
}

func TestFork(t *testing.T) {
	config := devConfig()
	anonConfig, _ := NewConfig(config.URL)

	testPaste := Paste{Private: true, Contents: "hi"}

	// We can't delete this one because it is anonymous, and we
	// can't fork pastes that we own. It was either this or make
	// tests require two different refheap user configs.
	if err := testPaste.Create(&anonConfig); err != nil {
		t.Errorf("Something went wrong creating a paste: %v", err)
	}

	pasteCopy := testPaste
	if err := pasteCopy.Fork(&config); err != nil {
		t.Errorf("Something went wrong forking a paste: %v", err)
	}
	defer pasteCopy.Delete(&config)

	pasteCopy.ID = testPaste.ID
	// We want to make sure that the new paste is under our account.
	testPaste.User = config.User
	pasteCopy.Date = testPaste.Date
	pasteCopy.URL = testPaste.URL
	if pasteCopy != testPaste {
		t.Errorf("Expected %#v, got %#v.", testPaste, pasteCopy)
	}
}

func TestSave(t *testing.T) {
	config := devConfig()
	paste := Paste{Private: true, Contents: "hi"}

	if err := paste.Create(&config); err != nil {
		t.Errorf("Something went wrong creating a paste: %v", err)
	}

	defer paste.Delete(&config)
	const newContents = "hi there"
	paste.Contents = newContents
	if err := paste.Save(&config); err != nil {
		t.Errorf("Something went wrong saving a paste: %v", err)
	}

	newPaste := Paste{ID: paste.ID}

	if err := newPaste.Get(&config); err != nil {
		t.Errorf("Something went wrong getting a paste %v", err)
	}

	if contents := newPaste.Contents; contents != "hi there" {
		t.Errorf("Expected %#v but got %#v.", newContents, contents)
	}
}
