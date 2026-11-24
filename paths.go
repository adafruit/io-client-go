package adafruitio

import (
	"errors"
	"net/url"
	"strings"
)

var (
	errEmptyFeedKey  = errors.New("feed key must not be empty")
	errEmptyGroupKey = errors.New("group key must not be empty")
	errEmptyID       = errors.New("id must not be empty")
)

// feedPath builds the path to the Feed identified by key, followed by any
// further segments, e.g. feedPath("temperature", "data", "last") returns
// "feeds/temperature/data/last". Every segment is escaped.
func feedPath(key string, segments ...string) (string, error) {
	if key == "" {
		return "", errEmptyFeedKey
	}
	return joinPath(append([]string{"feeds", key}, segments...)...), nil
}

// groupPath builds the path to the Group identified by key, followed by any
// further segments. Every segment is escaped.
func groupPath(key string, segments ...string) (string, error) {
	if key == "" {
		return "", errEmptyGroupKey
	}
	return joinPath(append([]string{"groups", key}, segments...)...), nil
}

// dataPath builds the path to the data of the Feed identified by feedKey,
// followed by any further segments, e.g. dataPath("temperature", "last")
// returns "feeds/temperature/data/last".
func dataPath(feedKey string, segments ...string) (string, error) {
	return feedPath(feedKey, append([]string{"data"}, segments...)...)
}

// dataRecordPath builds the path to the Data record identified by id.
func dataRecordPath(feedKey, id string) (string, error) {
	if feedKey != "" && id == "" {
		return "", errEmptyID
	}
	return dataPath(feedKey, id)
}

// joinPath escapes each segment and joins them with slashes.
func joinPath(segments ...string) string {
	escaped := make([]string, len(segments))
	for i, segment := range segments {
		escaped[i] = escapePathSegment(segment)
	}
	return strings.Join(escaped, "/")
}

// escapePathSegment escapes s so it is used as a single path segment, even if
// it contains characters such as "/", "?", "#", or spaces. "." and ".." are
// escaped too, since URL resolution would otherwise treat them as references
// to the current and parent paths.
func escapePathSegment(s string) string {
	switch s {
	case ".":
		return "%2E"
	case "..":
		return "%2E%2E"
	}
	return url.PathEscape(s)
}
