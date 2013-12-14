// A comprehensive Refheap API client for Go!
package goheap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// Default URL for refheap. This is the official site.
const RefheapURL = "https://www.refheap.com/api"

// There is a bit of configuration in this client, and this holds it.
// Fields:
//    URL  -- The URL to refheap's API.
//    User -- The username to authenticate with.
//    Key  -- The API key to authenticate with.
type Config struct {
	URL  string
	User string
	Key  string
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
//    ID  -- ID of paste.
//    Language -- Paste language.
//    Private  -- Whether or not the paste is private or not.
//    URL      -- URL to the paste.
//    User     -- User who owns the paste.
//    Contents -- Contents of the paste.
type Paste struct {
	// We need to tag these fields to tell the json parser what keys to
	// look for and produce. Refheap is case sensitive.
	Lines    int    `json:"lines"`
	Views    int    `json:"views"`
	Date     string `json:"date"`
	ID       string `json:"paste-id"`
	Language string `json:"language"`
	Private  bool   `json:"private"`
	URL      string `json:"url"`
	User     string `json:"user"`
	Contents string `json:"contents"`
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

func readBody(resp *http.Response) (body []byte, err error) {
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

// The body parsing code will always be the same. parseBody handles
// the JSON parsing and error handling (including handling
// RefheapError).
func parseBody(resp *http.Response, to interface{}) (err error) {
	body, err := readBody(resp)
	if err != nil {
		return
	}
	var newErr RefheapError
	if err = json.Unmarshal(body, &newErr); err != nil {
		return
	}
	if newErr.ErrorMessage != "" {
		err = newErr
	} else if to != nil {
		err = json.Unmarshal(body, to)
	}
	return
}

// If we are properly comfigured for authentication, this function
// will apply it to our data.
func addAuth(data *url.Values, config *Config) {
	if user := config.User; user != "" {
		data.Add("username", user)
		data.Add("token", config.Key)
	}
}

// Get a Paste from refheap. Result will be a Paste or an error
// if something goes wrong.
func (paste *Paste) Get(config *Config) (err error) {
	resp, err := http.Get(config.URL + "/paste/" + paste.ID)
	if err == nil {
		err = parseBody(resp, paste)
	}
	return
}

// Creating and saving are both the same thing as far as goheap
// is concerned. The only thing that changes is the endpoint to
// hit.
func (paste *Paste) createOrSave(endpoint string, config *Config) (err error) {
	data := url.Values{}
	addAuth(&data, config)
	if cont := paste.Contents; cont != "" {
		data.Add("contents", cont)
	}
	if lang := paste.Language; lang != "" {
		data.Add("language", lang)
	}
	data.Add("private", strconv.FormatBool(paste.Private))
	resp, err := http.PostForm(endpoint, data)
	if err != nil {
		return
	}
	err = parseBody(resp, paste)
	return
}

// Create a new paste from a Paste.
func (paste *Paste) Create(config *Config) error {
	return paste.createOrSave(config.URL+"/paste", config)
}

// Delete a paste. Requires you to have configured authentication.
func (paste *Paste) Delete(config *Config) (err error) {
	data := &url.Values{}
	addAuth(data, config)
	finalUrl := fmt.Sprintf("%v/paste/%v?%v", config.URL, paste.ID, data.Encode())
	request, _ := http.NewRequest("DELETE", finalUrl, nil)
	resp, err := http.DefaultClient.Do(request)
	if resp.StatusCode != 204 {
		err = parseBody(resp, nil)
	}
	return
}

// Fork a paste.
func (paste *Paste) Fork(config *Config) (err error) {
	data := url.Values{}
	addAuth(&data, config)
	data.Add("id", paste.ID)
	resp, err := http.PostForm(fmt.Sprintf("%v/paste/%v/fork", config.URL, paste.ID), data)
	if err != nil {
		return
	}
	err = parseBody(resp, paste)
	return
}

// Edit a paste. Must be authenticated.
func (paste *Paste) Save(config *Config) (err error) {
	return paste.createOrSave(config.URL+"/paste/"+paste.ID, config)
}

type highlightedPaste struct {
	Content string
}

// Get the highlighted version of a paste.
func (paste *Paste) GetHighlighted(config *Config) (s highlightedPaste, err error) {
	resp, err := http.Get(config.URL + "/paste/" + paste.ID + "/highlight")
	if err != nil {
		return
	}
	err = parseBody(resp, &s)
	return
}
