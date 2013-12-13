// A comprehensive Refheap API client for Go!
package goheap

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

// Default URL for refheap. This is the official site.
const RefheapURL = "https://refheap.com/api"

// There is a bit of configuration in this client, and this holds it.
// Fields:
//    Url  -- The Url to refheap's API.
//    User -- The username to authenticate with.
//    Key  -- The API key to authenticate with.
type Config struct {
	Url string
	User string
	Key string
}

// If there is an error in the NewConfig function as a result of
// too many arguments, this error will be returned.
type ConfigError struct {
	Args []string
}

func (e *ConfigError) Error() string {
	msg := "Config could not be constructed from these args: "
	return msg + fmt.Sprint(e.Args)
}

// NewConfig is a convenience function for creating a new configuration
// struct for goheap. It takes variadic arguments and is meant to take
// up to three strings. If it receives one argument, it is assumed that
// this argument is a custom URL (for example, a local refheap instance).
// If two arguments are passed, they are assumed to be username and
// API Key. Official refheap URL is used. If three arguments are pased
// they are expected to be a refheap URL, username, and api key, in that
// order. If zero arguments are passed, you get an anonymous default
// config object. Pass more than that and you're going to get an error
// value back back as the second return value. Pretty cool, huh? You can
// also just create a Config struct the old fashioned way if you'd like,
// of course!
func NewConfig(args ...string) (config Config, err error) {
	switch len(args) {
	default:
		err = &ConfigError{args}
	case 0:
		config = Config{RefheapURL, "", ""}
	case 1:
		config = Config{args[0], "", ""}
	case 2:
		config = Config{RefheapURL, args[0], args[1]}
	case 3:
		config = Config{args[0], args[1], args[2]}
	}
	return
}

// A struct for holding a paste response.
// Fields:
//    Lines    -- Number of lines in paste.
//    Views    -- Number of views paste has.
//    Date     -- Date that paste was created.
//    PasteID  -- ID of paste.
//    Language -- Paste language.
//    Private  -- Whether or not the paste is private or not.
//    Url      -- URL to the paste.
//    User     -- User who owns the paste.
//    Contents -- Contents of the paste.
type Paste struct {
	Lines int
	Views int
	Date string
	PasteID string `json:"paste-id"`
	Language string
	Private bool
	Url string
	User string
	Contents string
}

// When Refheap gives us back a json object with an 'error'
// key, we return an error of this type. It has an
// ErrorMessage key to hold the error from Refheap.
type RefheapError struct {
	// Can't name this 'Error' because of the Error() function.
	ErrorMessage string `json:"error"`
}

func (e RefheapError) Error() string {
	return e.ErrorMessage
}

// The body parsing code will always be the same. parseBody handles
// the JSON parsing and error handling (including handling
// RefheapError).
func parseBody(resp *http.Response, to interface{}) (err error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var newErr RefheapError
	if err = json.Unmarshal(body, &newErr); err != nil {
		return
	}
	if newErr.ErrorMessage != "" {
		err = newErr
	} else {
		err = json.Unmarshal(body, to)
	}
	return
}

// Get a Paste from refheap. Result will be a Paste or an error
// if something goes wrong.
func GetPaste(config *Config, id string) (paste Paste, err error) {
	resp, err := http.Get(config.Url + "/paste/" + id)
	if err == nil {
		err = parseBody(resp, &paste)
	}
	return
}

