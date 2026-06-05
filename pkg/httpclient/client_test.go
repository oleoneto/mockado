package httpclient_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cleopatrio/mockado/pkg/httpclient"
)

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
