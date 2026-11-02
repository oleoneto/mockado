package httpclient_test

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/oleoneto/mockado/pkg/envloader"
	"github.com/oleoneto/mockado/pkg/helpers"
	"github.com/oleoneto/mockado/pkg/httpclient"
	"github.com/oleoneto/mockado/pkg/server"
	"github.com/stretchr/testify/assert"
)

/*
func Test_GetComments(t *testing.T) {
	client := httpclient.Client()

	resp, err := client.Get("http://localhost:4567/comments")

	if err != nil {
		t.Errorf("Error: %s", err)
	}
	fmt.Printf("%d - %s", resp.StatusCode, resp.Header.Get("Date"))

	// --------------------------------------------

	resp, err = client.Get("http://localhost:4567/comments")

	if err != nil {
		t.Errorf("Error: %s", err)
	}
	fmt.Printf("%d - %s", resp.StatusCode, resp.Header.Get("Date"))
}

func Test_GetNothing(t *testing.T) {
	client := httpclient.Client()

	resp, err := client.Get("http://localhost:4567/nothing")

	if err != nil {
		t.Errorf("Error: %s", err)
	}
	fmt.Printf("%d - %s", resp.StatusCode, resp.Header.Get("Date"))

	// --------------------------------------------

	resp, err = client.Get("http://localhost:4567/nothing")

	if err != nil {
		t.Errorf("Error: %s", err)
	}
	fmt.Printf("%d - %s", resp.StatusCode, resp.Header.Get("Date"))
}

func Test_Post(t *testing.T) {
	client := httpclient.Client()

	resp, err := client.Post("http://localhost:4567", "application/json", http.NoBody)

	if err != nil {
		t.Errorf("Error: %s", err)
	}
	fmt.Printf("%d - %s", resp.StatusCode, resp.Header.Get("Date"))

	// --------------------------------------------

	resp, err = client.Post("http://localhost:4567", "text/plain", http.NoBody)

	if err != nil {
		t.Errorf("Error: %s", err)
	}
	fmt.Printf("%d - %s", resp.StatusCode, resp.Header.Get("Date"))
}
*/

func TestIndexRoute(t *testing.T) {
	tests := []struct {
		description     string
		route           string
		method          string
		expectedLatency *time.Duration
		expectedError   bool
		expectedCode    int
		expectedBody    []byte
		expectedHeaders map[string]string
	}{
		{
			description:   "ping",
			route:         "/ping",
			method:        "POST",
			expectedError: false,
			expectedCode:  202,
			expectedBody:  []byte(`{"event": "ping", "server_response": "pong"}`),
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			description:     "me",
			route:           "/v1/me",
			method:          "GET",
			expectedError:   false,
			expectedCode:    200,
			expectedLatency: helpers.PointerTo(102 * time.Millisecond),
			expectedBody:    []byte(`{"user_id": "1"}`),
			expectedHeaders: map[string]string{
				"Content-Type":       "application/json",
				"X-Response-Latency": "102ms",
			},
		},
	}

	app := server.New(
		server.WithEnvLoader(envloader.New()),
		server.WithRoutes(
			httpclient.Route{
				Name:    "Current User",
				Methods: []string{"GET"},
				Path:    "/v1/me",
				Accepts: "application/json",
				Response: struct {
					Engine     httpclient.RouteEngine "yaml:\"engine,omitempty\" json:\"engine,omitempty\""
					Latency    *time.Duration         "yaml:\"latency,omitempty\" json:\"latency,omitempty\""
					StatusCode int                    "yaml:\"status,omitempty\" json:\"status,omitempty\""
					Body       any                    "yaml:\"body,omitempty\" json:\"body,omitempty\""
					EnvAsBody  bool                   "yaml:\"envAsBody\" json:\"env_as_body\""
					Headers    map[string]string      "yaml:\"headers,omitempty\" json:\"headers,omitempty\""
				}{
					Latency:    helpers.PointerTo(102 * time.Millisecond),
					StatusCode: http.StatusOK,
					Body:       []byte(`{"user_id": "1"}`),
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
			},
			httpclient.Route{
				Name:    "Ping",
				Methods: []string{"POST"},
				Path:    "/ping",
				Accepts: "application/json",
				Response: struct {
					Engine     httpclient.RouteEngine "yaml:\"engine,omitempty\" json:\"engine,omitempty\""
					Latency    *time.Duration         "yaml:\"latency,omitempty\" json:\"latency,omitempty\""
					StatusCode int                    "yaml:\"status,omitempty\" json:\"status,omitempty\""
					Body       any                    "yaml:\"body,omitempty\" json:\"body,omitempty\""
					EnvAsBody  bool                   "yaml:\"envAsBody\" json:\"env_as_body\""
					Headers    map[string]string      "yaml:\"headers,omitempty\" json:\"headers,omitempty\""
				}{
					StatusCode: http.StatusAccepted,
					Body:       []byte(`{"event": "ping", "server_response": "pong"}`),
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
			},
		),
	).App()

	for _, test := range tests {
		req, err := http.NewRequest(test.method, test.route, nil)
		res, err := app.Test(req, fiber.TestConfig{Timeout: 0, FailOnTimeout: true})

		assert.Equal(t, test.expectedError, err != nil, test.description)
		if test.expectedError {
			continue
		}
		defer res.Body.Close()

		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
		body, err := io.ReadAll(res.Body)

		assert.Equal(t, test.expectedBody, body, test.description)

		for key := range test.expectedHeaders {
			assert.Equal(
				t,
				test.expectedHeaders[key],
				res.Header.Get(key),
				test.description,
			)
		}
	}
}
