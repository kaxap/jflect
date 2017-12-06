# Jflect

Forked from mrosset/jflect.
Added `-tag` flag. Adds custom flag additional to "json". Example: ```jflect -tag db -tag rmq``` will produce struct with `json`, `db` and `rmq` tags.

Takes JSON from stdin and outputs go structs to stdout.

## Documentation
[godoc.org](http:godoc.org/github.com/str1ngs/jflect)
