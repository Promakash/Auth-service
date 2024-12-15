package tests

import (
	apitypes "auth_service/api/http/types"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const genHandlerURL = "http://app:8080/api/v1/auth/tokens/generate"

func TestGenerateHandler_Correct(t *testing.T) {
	testUUID := "123e4567-e89b-12d3-a456-426614174005"

	resp, err := http.Get(fmt.Sprintf("%s?uuid=%s", genHandlerURL, testUUID))
	assert.NoError(t, err, "error while requesting")
	defer func() { _ = resp.Body.Close() }()

	expected := http.StatusOK
	actual := resp.StatusCode
	assert.Equal(t, expected, actual, fmt.Sprintf("handler returned different code. expected: %d, got: %d", expected, actual))

	var result apitypes.GenerateResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err, "error while parsing responses body")

	assert.NotEmpty(t, result.AccessToken, "AccessToken is empty")
	assert.NotEmpty(t, result.RefreshToken, "RefreshToken is empty")
}

func TestGenerateHandler_InvalidUUID(t *testing.T) {
	testUUID := "123e46614174005"

	resp, err := http.Get(fmt.Sprintf("%s?uuid=%s", genHandlerURL, testUUID))
	assert.NoError(t, err, "error while requesting")
	defer func() { _ = resp.Body.Close() }()

	expected := http.StatusBadRequest
	actual := resp.StatusCode
	assert.Equal(t, expected, actual, fmt.Sprintf("handler returned different code. expected: %d, got: %d", expected, actual))
}

func TestGenerateHandler_MissingUUID(t *testing.T) {
	resp, err := http.Get(genHandlerURL)
	assert.NoError(t, err, "request error")
	defer func() { _ = resp.Body.Close() }()

	expected := http.StatusBadRequest
	actual := resp.StatusCode
	assert.Equal(t, expected, actual, fmt.Sprintf("handler returned different code. expected: %d, got: %d", expected, actual))
}

func TestGenerateHandler_NonExistedUser(t *testing.T) {
	testUUID := "f011a1a5-0d60-4467-b254-4ecb7e71f75a"

	resp, err := http.Get(fmt.Sprintf("%s?uuid=%s", genHandlerURL, testUUID))
	assert.NoError(t, err, "error while requesting")
	defer func() { _ = resp.Body.Close() }()

	expected := http.StatusNotFound
	actual := resp.StatusCode
	assert.Equal(t, expected, actual, fmt.Sprintf("handler returned different code. expected: %d, got: %d", expected, actual))
}
