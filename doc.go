// Copyright (c) 2016, Adafruit Author: Adam Bachman
// This code is under the MIT License that can be found in the LICENSE file.

/*
Package adafruitio is a client for the Adafruit IO HTTP API (v2) at
https://io.adafruit.com.

	import adafruitio "github.com/adafruit/io-client-go/v3"

Create a client with your Adafruit IO username and key. The key is sent in
the X-AIO-Key header of every request. We recommend keeping it in an
environment variable rather than in your code:

	client := adafruitio.NewClient(os.Getenv("ADAFRUIT_IO_USERNAME"), os.Getenv("ADAFRUIT_IO_KEY"))

To set timeouts or other transport settings, pass your own http.Client with
WithHTTPClient.

Every API method takes a context.Context as its first argument:

	feeds, _, err := client.Feed.All(ctx)

# Feeds, Data, and Groups

The client's Feed, Data, and Group services cover the matching parts of the
API. Data methods take the key of the feed they read from or write to:

	_, _, err := client.Data.Create(ctx, "temperature", &adafruitio.Data{
		Value: adafruitio.Ptr("22.5"),
	})

	last, _, err := client.Data.Last(ctx, "temperature")

DataService.All returns one page of values at a time, newest first. To get
every page, set the DataFilter's Before to the NextBefore of each page's
Response, until it is empty:

	filter := &adafruitio.DataFilter{Limit: 1000, EndTime: time.Now()}
	for {
		page, resp, err := client.Data.All(ctx, "temperature", filter)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(len(page), "values")
		if resp.NextBefore == "" {
			break
		}
		filter.Before = resp.NextBefore
	}

# Pointer fields

As in go-github, every field of Feed, Group, and Data is a pointer. A nil
field is left out of create and update requests, so the server keeps its
current value, while a non-nil field is always sent, even if it is "", 0, or
false. Set fields with Ptr, and read them with the generated Get methods,
which return the zero value for nil fields:

	feed, _, err := client.Feed.Update(ctx, "temperature", &adafruitio.Feed{
		Description: adafruitio.Ptr(""),    // clear the description
		History:     adafruitio.Ptr(false), // turn off history
	})
	fmt.Println(feed.GetName())

# Times

Times such as CreatedAt are *Timestamp values. Timestamp embeds time.Time, so
any time.Time method can be used on it:

	fmt.Println(feed.GetCreatedAt().Format(time.RFC1123))

The StartTime and EndTime of a DataFilter are plain time.Time values.

# Errors

When the API returns an error status, methods return an *ErrorResponse
containing the API's error message, along with the *Response.

The Adafruit IO HTTP API documentation is at
https://io.adafruit.com/api/docs/http.html.
*/
package adafruitio
