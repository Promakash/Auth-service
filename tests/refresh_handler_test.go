package tests

import (
	apitypes "auth_service/api/http/types"
	"auth_service/config"
	pkgConfig "auth_service/pkg/config"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mhale/smtpd"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"testing"
	"time"
)

const refreshHandlerURL = "http://app:8080/api/v1/auth/tokens/refresh"

const serverConfigEnv = "HTTP_CONFIG_PATH"

func sendRequest(method, url string, requestBody interface{}, IP string) (*http.Response, error) {
	var body *bytes.Reader
	if requestBody != nil {
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if IP != "" {
		req.Header.Set("X-Forwarded-For", IP)
	}

	client := &http.Client{}
	return client.Do(req)
}

func TestRefreshHandler_Correct(t *testing.T) {
	testUUID := "123e4567-e89b-12d3-a456-426614174005"

	resp, err := http.Get(fmt.Sprintf("%s?uuid=%s", genHandlerURL, testUUID))
	assert.NoError(t, err, "error while requesting")
	defer func() { _ = resp.Body.Close() }()

	var result apitypes.GenerateResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err, "error while parsing responses body")

	time.Sleep(time.Second)

	req := apitypes.RefreshRequest{RefreshToken: result.RefreshToken}
	respRefresh, err := sendRequest(http.MethodPost, refreshHandlerURL, req, "")
	assert.NoError(t, err, "error while sending request")
	defer func() { _ = respRefresh.Body.Close() }()

	expected := http.StatusOK
	actual := respRefresh.StatusCode
	assert.Equal(t, expected, actual, fmt.Sprintf("handler returned different code. expected: %d, got: %d", expected, actual))

	var respBody apitypes.RefreshResponse
	err = json.NewDecoder(respRefresh.Body).Decode(&result)
	assert.NoError(t, err, "error while parsing responses body")

	assert.NotEqual(t, result.AccessToken, respBody.AccessToken, "access tokens must be different")
	assert.NotEqual(t, result.RefreshToken, respBody.RefreshToken, "refresh tokens must be different")
}

func TestRefreshHandler_InvalidToken(t *testing.T) {
	testToken := "fsafasfafafafsafsafafafasfa"

	req := apitypes.RefreshRequest{RefreshToken: testToken}
	respRefresh, err := sendRequest(http.MethodPost, refreshHandlerURL, req, "")
	assert.NoError(t, err, "error while sending request")
	defer func() { _ = respRefresh.Body.Close() }()

	expected := http.StatusUnauthorized
	actual := respRefresh.StatusCode
	assert.Equal(t, expected, actual, fmt.Sprintf("handler returned different code. expected: %d, got: %d", expected, actual))
}

func TestGenerateHandler_SameTokenTwice(t *testing.T) {
	testUUID := "123e4567-e89b-12d3-a456-426614174005"

	resp, err := http.Get(fmt.Sprintf("%s?uuid=%s", genHandlerURL, testUUID))
	assert.NoError(t, err, "error while requesting")
	defer func() { _ = resp.Body.Close() }()

	var result apitypes.GenerateResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err, "error while parsing responses body")

	time.Sleep(time.Second)

	req := apitypes.RefreshRequest{RefreshToken: result.RefreshToken}
	respRefresh, err := sendRequest(http.MethodPost, refreshHandlerURL, req, "")
	assert.NoError(t, err, "error while sending request")
	_ = respRefresh.Body.Close()

	expected := http.StatusOK
	actual := respRefresh.StatusCode
	assert.Equal(t, expected, actual, fmt.Sprintf("handler returned different code. expected: %d, got: %d", expected, actual))

	time.Sleep(time.Second)

	respRefresh, err = sendRequest(http.MethodPost, refreshHandlerURL, req, "")
	assert.NoError(t, err, "error while sending request")
	defer func() { _ = respRefresh.Body.Close() }()

	expected = http.StatusUnauthorized
	actual = respRefresh.StatusCode
	assert.Equal(t, expected, actual, fmt.Sprintf("handler returned different code. expected: %d, got: %d", expected, actual))
}

func TestGenerateHandler_CheckUserNotification(t *testing.T) {
	cfg := pkgConfig.ParseAppConfig[config.HTTPConfig](serverConfigEnv)
	received := make(chan string, 1)
	go func() {
		err := smtpd.ListenAndServe(
			fmt.Sprintf("%s:%s", cfg.SMTP.Host, cfg.SMTP.Port),
			func(origin net.Addr, from string, to []string, data []byte) error {
				received <- string(data)
				return nil
			},
			"",
			"",
		)
		if err != nil {
			t.Errorf("failed to start SMTP server: %v", err)
			return
		}
	}()

	time.Sleep(time.Second * 5)

	testUUID := "123e4567-e89b-12d3-a456-426614174005"

	resp, err := http.Get(fmt.Sprintf("%s?uuid=%s", genHandlerURL, testUUID))
	assert.NoError(t, err, "error while requesting")
	defer func() { _ = resp.Body.Close() }()

	var result apitypes.GenerateResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err, "error while parsing responses body")

	IP := "192.168.1.1"
	req := apitypes.RefreshRequest{RefreshToken: result.RefreshToken}
	respRefresh, err := sendRequest(http.MethodPost, refreshHandlerURL, req, IP)
	assert.NoError(t, err, "error while sending request")
	defer func() { _ = respRefresh.Body.Close() }()

	expected := http.StatusOK
	actual := respRefresh.StatusCode
	assert.Equal(t, expected, actual, fmt.Sprintf("handler returned different code. expected: %d, got: %d", expected, actual))

	select {
	case msg := <-received:
		assert.Contains(t, msg, "Suspicious activity detected")
	case <-time.After(5 * time.Second):
		t.Fatalf("did not receive email in time")
	}
}
