package httpclient

import "time"

type Route struct {
	// A name to identify the route (not checked for uniqueness)
	Name string `yaml:"name,omitempty" json:"name,omitempty"`

	// HTTP request path
	//
	// Default: none
	Path string `yaml:"path,omitempty" json:"path,omitempty"`

	// HTTP route methods
	//
	// Default: [HEAD, GET]
	Methods []string `yaml:"methods,omitempty" json:"methods,omitempty"`

	// HTTP Accept header
	//
	// Default: node
	Accepts string `yaml:"accepts,omitempty" json:"accepts,omitempty"`

	Response struct {
		// Engine to render the response body with
		//
		// Default: none
		Engine RouteEngine `yaml:"engine,omitempty" json:"engine,omitempty"`

		// HTTP response delay
		//
		// Default: none
		Latency *time.Duration `yaml:"latency,omitempty" json:"latency,omitempty"`

		// HTTP response status code
		//
		// Default: none
		StatusCode int `yaml:"status,omitempty" json:"status,omitempty"`

		// HTTP response body
		//
		// Default: none
		Body any `yaml:"body,omitempty" json:"body,omitempty"`

		// Output the application's environment as the response body
		EnvAsBody bool `yaml:"envAsBody" json:"env_as_body"`

		// HTTP response headers
		//
		// Default: none
		Headers map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	} `yaml:"response,omitempty" json:"response,omitempty"`
}

type RouteEngine string

const (
	// Render the response body as-is, without any processing.
	Literal RouteEngine = "literal" // default

	// Support for the following syntax:
	//
	// {{ range 5 }}
	// {{ . }}
	// {{ end }}
	GoTemplate RouteEngine = "go-template"

	// Support for the following syntax:
	//
	// {{ for i in std.range(5) }}
	// {{ i }}
	// {{ end }}
	Jsonnet RouteEngine = "jsonnet"

	// Support for the following syntax:
	//
	// package router
	//
	// func Route() string {
	//   return "Hello, World!"
	// }
	Go RouteEngine = "go"
)
