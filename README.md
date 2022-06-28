# talk to adafruit io from Go

[![GoDoc](http://godoc.org/github.com/adafruit/io-client-go?status.svg)](http://godoc.org/github.com/adafruit/io-client-go)
[![Build Status](https://travis-ci.org/adafruit/io-client-go.svg?branch=master)](https://travis-ci.org/adafruit/io-client-go)

A go client library for talking to your io.adafruit.com account.

Requires go version 1.18 or better. Running tests uses the github.com/stretchr/testify library. To run tests, run:

```bash
$ go test ./...
```

## Usage

First, add the package:
```bash
$ go get github.com/adafruit/io-client-go/v2
```

Then import it:
```go
import "github.com/adafruit/io-client-go/v2"
```

The io-client-go repository provides the `adafruitio` package.

Authentication for Adafruit IO is managed by providing your Adafruit IO username and token. The token is sent in the head of all web requests via the `X-AIO-Key` header. This is handled for
you by the client library, which expects you API Token when it is initialized.

We recommend keeping the Token in an environment variable to avoid including it
directly in your code.

```go
client := adafruitio.NewClient(os.Getenv("ADAFRUIT_IO_USERNAME"), os.Getenv("ADAFRUIT_IO_KEY"))
feeds, _, err := adafruitio.Feed.All()
```

Some API calls expect parameters, which must be provided when making the call.

```go
feed := &aio.Feed{Name: "my-new-feed"}
feed := client.Feed.Create(newFeed)
```

Data related API calls expect a Feed to be set before the call is made.


```go
feed, _, ferr := client.Feed.Get("my-new-feed")
client.SetFeed(feed)
client.Data.Create(&adafruitio.Data{Value: 100})
```

More detailed example usage can be found in the [./examples](./examples) directory

For full package documentation, visit the godoc page at https://godoc.org/github.com/adafruit/io-client-go

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
