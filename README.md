# talk to adafruit io from Go

A go client library for talking to your io.adafruit.com account.

## Usage

```go
import "github.com/adafruit/io-client-go"
```

The io-client-go repository provides the `adafruitio` package.

Authentication for Adafruit IO is managed by providing your Adafruit IO token
in the head of all web requests via the `X-AIO-Key` header. This is handled for
you by the client library, which expects you API Token when it is initialized.

We recommend keeping the Token in an environment variable to avoid including it
directly in your code.

```go
client := adafruitio.NewClient(os.Getenv("ADAFRUIT_IO_KEY"))
feeds, _, err := adafruitio.Feed.All()
```

Some API calls expect parameters, which must be provided when making the call.

```go
feed := &aio.Feed{Name: "my-new-feed"}
client.Feed.Create(newFeed)
```

Data related API calls expect a Feed to be set before the call is made.

**NOTE:** the Feed doesn't have to exist yet if you're using the `Data.Send()`
method, but it still needs to be set. If you're relying on the Data API to
create the Feed, make sure you set the `Key` attribute on the new Feed.

```go
feed := &aio.Feed{Name: "My New Feed", Key: "my-new-feed"}
client.SetFeed(newFeed)
client.Data.Send(&adafruitio.Data{Value: 100})
```

For full package documentation, visit the godoc page at https://godoc.org/github.com/adafruit/io-client-go

## License

Copyright (c) 2016 Adafruit Industries. Licensed under the MIT license.

## Contributing

- Fork it ( http://github.com/adafruit/io-client-go/fork )
- Create your feature branch (git checkout -b my-new-feature)
- Commit your changes (git commit -am 'Add some feature')
- Push to the branch (git push origin my-new-feature)
- Create new Pull Request

