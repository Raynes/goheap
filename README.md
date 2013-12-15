# Goheap

Goheap is a [refheap](https://www.refheap.com) API wrapper for Go. It supports
refheap's entire API and has a full test suite.

## Usage

You can find the godocs [here](http://godoc.org/github.com/Raynes/goheap).

This library is meant to be pretty easy to use. Let's go over the API!

### Config

Some refheap endpoints require you be authenticated with the service, and
others work if authenticated to make changes targeting your user account.
Furthermore, what if you want to use a different refheap instance than the one
Anthony Grimes (me!) runs? `Config` dat! First of all, let's import the package:

```go
import "github.com/raynes/goheap"
```

Now, let's explore how we can configure goheap. You have two options:

* Use the `NewConfig` convenience function.
* Construct the `Config` manually.

Either of these are fine. I've included `NewConfig` solely for convenience. It
takes variadic arguments (one, two, or three are accepted) and has the following
rules:

* If only one argument is passed, it is assumed to be a URL to the instance of
  refheap you'd like to use.
* If two arguments are passed, they are expected to be refheap username and api
  token, in that order.
* If three arguments are passed, they are expected to be the URL to refheap's
  API, username, and api token, in that order.

Examples:

```go
config := goheap.NewConfig("localhost:8080")
```

The above is the same as this:

```go
config := goheap.Config{"localhost:8080", "", ""}
```

And this:

```go
config := goheap.NewConfig("foo", "bar")
```

is the same as:

```go
config := goheap.Config{goheap.RefheapURL, "foo", "bar")
```

That's about it for configuration. Keep in mind that most endpoints allow for
anonymous usage as well. If no username is configured pastes are created
anonymously, etc, but `Save()` and `Delete()` won't as they both only work on
non-anon pastes. Let's move on to those fancy pastes!

### Pastes

The most interesting thing you'll see is the `Paste` type. It represents a
refheap paste and matches the JSON blob that refheap returns when you create,
edit, or fork a paste. This is the most important part of goheap! Let's see what
we can do with it. We'll start by creating a new `Paste`.

```go
paste := goheap.Paste{Private: true, Language: "Go", Contents: "2 + 2"}
```

Now we can call methods on this `Paste`!

#### Creating Pastes

Just because we made a new `Paste` doesn't mean it actually exists yet! We
should create it on refheap now.

```go
err = paste.Create(&config)
```

And that should do it! If the creation failed for some reason, `err` will be
non-nil. If it was refheap itself that returned the error as json, it'll be a
`RefheapError`.

Now `paste` should be filled in with the info about the new paste. I'm not going
to talk about all of these attributes (see the godocs), but the one you're
probably most interested in is the URL to the new paste. We can easily retrieve
that:

```go
paste.URL
```

This will be a string to the new paste.


#### Saving Pastes

The `Paste` type, as you would expect, is mutable. You can change its fields. It
makes perfect sense then that in order to edit a paste on refheap, you can just
save an existing paste whose fields you have changed since creation!

```go
paste.Contents = "3 + 3" // Bigger and better things.
err = paste.Save(&config)
```

Done. Paste edited. If you want to edit a paste that was not created by goheap
(and thus you do not have a `Paste`), just create a dummy `Paste` with the `ID`
field pointing to the right paste!

```go
paste := goheap.Paste{ID: "30000", Contents: "4 + 4"}
paste.Save()
```


#### Deleting Pastes

Creating and editing pastes is fun, but perhaps we're sick of our latest paste
and just want to delete it? Simple enough:

```go
paste.Delete(&config)
```

As with saving pastes, if you want to delete a paste you didn't create with
goheap, just create a dummy `Paste` with an ID field pointing to the correct
paste.

#### Getting Pastes

Assume there is a paste on Refheap with id '30000'. Let's see how we'd fetch
this paste.

```go
paste := goheap.Paste{ID: "30000"}
paste.Get(&config)
```

You do not have to be authenticated to get a private paste, since the URL is
what is secret. The reason `&config` is passed in is in case the url to the
refheap instance is not the default.

#### Pygmentized

An interesting feature of refheap is that its API allows you to get the
pygmentized HTML of a paste.

```go
paste := goheap.Paste{ID: "30000"}
s, err := paste.GetHighlighted(&config)
```

If this succeeded, `s` will be a string containing the HTML.

## FIN

Well, I sure hope goheap turns out to be useful for ya. If you have any
suggestions, feedback, issues, etc, just pop open a ticket or contact me via
email/IRC and I'll get back to you asap. Pull requests highly appreciated.

