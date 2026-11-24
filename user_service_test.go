package adafruitio

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserGet(t *testing.T) {
	// The API documents a plain user record but returns a wrapped one, so
	// check both. /api/v2/user is not under the username.
	for name, body := range map[string]string{
		"wrapped": `{"user": {"id": 2, "username": "jwcooper", "time_zone": "America/Chicago", "created_at": "2026-08-27T14:42:48Z"}, "profile": {}, "throttle": {}}`,
		"plain":   `{"id": 2, "username": "jwcooper", "time_zone": "America/Chicago", "created_at": "2026-08-27T14:42:48Z"}`,
	} {
		t.Run(name, func(t *testing.T) {
			setup()
			defer teardown()

			mux.HandleFunc(APIPath+"/user",
				func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					fmt.Fprint(w, body)
				},
			)

			user, response, err := client.User.Get(context.Background())

			assert.Nil(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, int64(2), user.GetID())
			assert.Equal(t, "jwcooper", user.GetUsername())
			assert.Equal(t, "America/Chicago", user.GetTimeZone())
			assert.Equal(t, 2026, user.GetCreatedAt().Year())
		})
	}
}

func TestUserThrottle(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("throttle"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, `{"data_rate_limit": 60, "active_data_rate": 3, "mqtt_ban_error_message": null}`)
		},
	)

	throttle, _, err := client.User.Throttle(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, 60, throttle.GetDataRateLimit())
	assert.Equal(t, 3, throttle.GetActiveDataRate())
	assert.Nil(t, throttle.MQTTBanErrorMessage)
}
