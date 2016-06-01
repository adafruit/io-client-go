// Copyright (c) 2016, Adafruit Author: Adam Bachman
// This code is under the MIT License that can be found in the LICENSE file.

/*
The adafruitio package is a simple HTTP client for accessing v1 of the Adafruit IO REST API at https://io.adafruit.com.

	import "github.com/adafruit/io-client-go"

Authentication for Adafruit IO is managed by providing your Adafruit IO token
in the head of all web requests via the `X-AIO-Key` header. This is handled for
you by the client library, which expects you API Token when it is initialized.

We recommend keeping the Token in an environment variable to avoid including it
directly in your code.

	client := adafruitio.NewClient(os.Getenv("ADAFRUIT_IO_KEY"))
	feeds, _, err := adafruitio.Feed.All()

Some API calls expect parameters, which must be provided when making the call.

	feed := &aio.Feed{Name: "my-new-feed"}
	client.Feed.Create(newFeed)

Data related API calls expect a Feed to be set before the call is made.

**NOTE:** the Feed doesn't have to exist yet if you're using the `Data.Send()`
method, but it still needs to be set. If you're relying on the Data API to
create the Feed, make sure you set the `Key` attribute on the new Feed.

	feed := &aio.Feed{Name: "My New Feed", Key: "my-new-feed"}
	client.SetFeed(newFeed)
	client.Data.Send(&adafruitio.Data{Value: 100})

You can see the v1 Adafruit IO REST API documentation online at https://io.adafruit.com/api/docs/
*/
package adafruitio
