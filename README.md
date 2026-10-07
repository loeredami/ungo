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
- [Main](./application_wrapper.go)
- [Bitmap](./bitmap.go)
- [SmallMap](./small_map.go) highly recommend testing this one too, interesting performance results in _some_ cases.
  - I have made that thing when it was like 2 am, I have no Idea how this works. I would probably have been able to tell you how it works when I made it, but right the next morning after? Nah. 
  - I should read again, and figure out how it works step by step when I have the time. All I know is it uses some bit magic keying. I don't know if it either saved memory, or performance, but I remember it doing either one, and have been using it in [LoinCloth](https://github.com/loeredami/loincloth) for that reason. Use as the name suggests.

* (A lot of things have been made when i was very tired, hence the lack of documentations. I don't know the person I am when I am that tired, but I wish I could be that person on demand.

- [Lazy](lazy.go) foundation of the:
- [Registry](registry.go) , only setting keys in a map when needed can be very handy.

  - However Lazy is bugged, and works more like a Getter currently, planning on fixing that some time. -TODO

- [Result](./result.go)
- [Exception](./exception.go)
- [Default](./default.go) though doesn't seem practical
- [Optional](./optional.go) Very convenient.
