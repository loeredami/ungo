# ungo
Go prides itself on what it lacks. ungo provides what it needs. It is a collection of common types and utilities that prioritize developer utility over the strict minimalist philosophy  of the standard library.

This is a nightmare to work with, and not fully tested, however I do plan on testing all of the things I have written, as most things are just things I could come up with but didn't bother using personally yet.

[I work with this a lot, though, I do not work with it enough.](https://en.wikipedia.org/wiki/Satire)

# install
```bash
go get github.com/loeredami/ungo
```

# update
```bash
GOPROXY=direct go get -u github.com/loeredami/ungo
```

What has been tested so far:
[Main](./application_wrapper.go)
[Bitmap](./bitmap.go)
[SmallMap](./small_map.go) highly recommend testing this one too, interesting performance results in _some_ cases.
[Result](./result.go)
[]
