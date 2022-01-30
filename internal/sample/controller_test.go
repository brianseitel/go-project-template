package sample_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianseitel/go-project-template/internal/sample"
	"github.com/stretchr/testify/assert"
)

func TestHello(t *testing.T) {
	// Create a sample request
	req := httptest.NewRequest(http.MethodGet, "/v1/hello", nil)

	// Create a recorder. This records the HTTP response from the server and
	// acts as a test response, effectively.
	recorder := httptest.NewRecorder()

	// Instantiate the controller
	controller := sample.Controller{}

	// Simulate the request
	controller.Hello().ServeHTTP(recorder, req)

	// Make sure we got an HTTP 200 OK response
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Read the response body
	body, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		panic(err)
	}

	// Convert the raw JSON response into a SampleResponse struct
	var response sample.SampleResponse
	err = json.Unmarshal(body, &response)
	assert.Nil(t, err)

	// Validate the response
	assert.Equal(t, "Hello World", response.Greeting)
}
