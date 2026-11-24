# talk to adafruit io from Go

[![Go Reference](https://pkg.go.dev/badge/github.com/adafruit/io-client-go/v3.svg)](https://pkg.go.dev/github.com/adafruit/io-client-go/v3)
[![CI](https://github.com/adafruit/io-client-go/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/adafruit/io-client-go/actions/workflows/ci.yml)

A Go client library for the [Adafruit IO HTTP API (v2)](https://io.adafruit.com/api/docs/http.html).

Requires Go 1.18 or later.

## Installation

```bash
go get github.com/adafruit/io-client-go/v3
```

```go
import adafruitio "github.com/adafruit/io-client-go/v3"
```

## Usage

### Creating a client

Create a client with your Adafruit IO username and key. We recommend keeping
your key in an environment variable rather than in your code:

```go
client := adafruitio.NewClient(os.Getenv("ADAFRUIT_IO_USERNAME"), os.Getenv("ADAFRUIT_IO_KEY"))
```

The client sends your key in the `X-AIO-Key` header of every request.

By default, the client uses `http.DefaultClient`. To set timeouts, a proxy, or
other transport settings, pass your own `http.Client`:

```go
client := adafruitio.NewClient(username, key,
	adafruitio.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
)
```

### Contexts

Every API method takes a `context.Context` as its first argument, so you can
cancel requests or give them deadlines. If the context is cancelled or times
out, the method returns `context.Canceled` or `context.DeadlineExceeded`.

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

feeds, _, err := client.Feed.All(ctx)
if err != nil {
	log.Fatal(err)
}
fmt.Println(len(feeds), "feeds")
```

### Feeds

```go
// Create a feed
feed, _, err := client.Feed.Create(ctx, &adafruitio.Feed{
	Name:        adafruitio.Ptr("Temperature"),
	Description: adafruitio.Ptr("Living room temperature"),
})
if err != nil {
	log.Fatal(err)
}
fmt.Println("created", feed.GetKey())

// Update it. Only the fields you set are sent.
feed, _, err = client.Feed.Update(ctx, feed.GetKey(), &adafruitio.Feed{
	History: adafruitio.Ptr(false),
})
if err != nil {
	log.Fatal(err)
}

// Delete it
_, err = client.Feed.Delete(ctx, feed.GetKey())
if err != nil {
	log.Fatal(err)
}
```

### Data

Data methods take the key of the feed they read from or write to:

```go
// Send a value to the "temperature" feed
data, _, err := client.Data.Create(ctx, "temperature", &adafruitio.Data{
	Value: adafruitio.Ptr("22.5"),
})
if err != nil {
	log.Fatal(err)
}
fmt.Println("created", data.GetID())

// Get the most recent value
last, _, err := client.Data.Last(ctx, "temperature")
if err != nil {
	log.Fatal(err)
}
fmt.Println(last.GetValue(), last.GetCreatedAt())

// List the values in a time range
values, _, err := client.Data.All(ctx, "temperature", &adafruitio.DataFilter{
	StartTime: "2026-11-01T00:00:00Z",
	EndTime:   "2026-11-02T00:00:00Z",
})
if err != nil {
	log.Fatal(err)
}
fmt.Println(len(values), "values")
```

### Groups

```go
groups, _, err := client.Group.All(ctx)
if err != nil {
	log.Fatal(err)
}
for _, group := range groups {
	fmt.Println(group.GetName(), len(group.Feeds), "feeds")
}
```

### Pointer fields

Like [go-github](https://github.com/google/go-github), every field of `Feed`,
`Group`, and `Data` is a pointer. This lets the library tell the difference
between a field you didn't set and one you set to an empty value:

- A `nil` field is left out of create and update requests, so the server keeps
  its current value.
- A non-`nil` field is always sent, even if it is `""`, `0`, or `false`.

Use `adafruitio.Ptr` to set fields. For example, to clear a feed's description
and turn off its history:

```go
_, _, err := client.Feed.Update(ctx, "temperature", &adafruitio.Feed{
	Description: adafruitio.Ptr(""),
	History:     adafruitio.Ptr(false),
})
if err != nil {
	log.Fatal(err)
}
```

To read fields, use the `Get` methods, such as `feed.GetName()`. They return
the zero value when a field (or the struct itself) is `nil`, so you don't need
nil checks.

### Errors

When the API returns an error status, methods return an
`*adafruitio.ErrorResponse` with the API's error message, along with the
`*adafruitio.Response` so you can inspect it further:

```go
_, _, err := client.Feed.Get(ctx, "does-not-exist")

var apiErr *adafruitio.ErrorResponse
if errors.As(err, &apiErr) {
	fmt.Println(apiErr.Response.StatusCode, apiErr.Message)
}
```

### More examples

The [examples](./examples) directory has complete programs for feeds, groups,
and data, and a small web server for browsing your feeds.

Full documentation is on [pkg.go.dev](https://pkg.go.dev/github.com/adafruit/io-client-go/v3).

## Upgrading from v2

v3 changes the API in several ways:

- **Import path:** `github.com/adafruit/io-client-go/v2` is now
  `github.com/adafruit/io-client-go/v3`.
- **Contexts:** every API method takes a `context.Context` as its first
  argument, and `Client.Do` is now `Do(ctx, req, v)`.
- **No more `SetFeed`:** Data methods take the feed key instead, for example
  `client.Data.Last(ctx, "temperature")`. `Client.SetFeed`,
  `FeedService.CurrentFeed`, and `FeedService.Path` are removed.
- **Pointer fields:** every field of `Feed`, `Group`, `Data`, and `Owner` is a
  pointer, and IDs are `int64`. Set fields with `adafruitio.Ptr` and read them
  with the `Get` methods. See [Pointer fields](#pointer-fields).
- **Feed fields:** `Feed.Shared` is removed, because the API doesn't return it
  for feeds. `Feed.Writable`, `Feed.Group`, and `Feed.Groups` are new.
- **Fixes:** `Data.Search` now applies its filter, `Delete` methods return
  the response alongside an error, and the client no longer prints to stdout
  when the API returns an error.

## Development

Run the tests with:

```bash
go test ./...
```

The `Get` methods in `accessors.go` are generated. After changing the fields of
`Feed`, `Group`, `Data`, or `Owner`, regenerate them with:

```bash
go generate ./...
```

## License

Copyright (c) 2016 Adafruit Industries. Licensed under the MIT license.

## Contributing

- Fork it ( http://github.com/adafruit/io-client-go/fork )
- Create your feature branch (git checkout -b my-new-feature)
- Commit your changes (git commit -am 'Add some feature')
- Push to the branch (git push origin my-new-feature)
- Create new Pull Request

---

[adafruit](https://adafruit.com) invests time and resources providing this open source code. please support adafruit and open-source hardware by purchasing products from [adafruit](https://adafruit.com).
