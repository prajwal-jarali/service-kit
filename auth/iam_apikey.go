package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/prajwal-jarali/service-kit/api"
)

type IamApiKey struct {
	mu              sync.Mutex
	Apikey          string
	generatedAt     time.Time
	token           string
	testEndpoint    bool
	privateEndpoint bool
}

func NewApiKey(key string, testEndpoint, privateEndpoint bool) (*IamApiKey, error) {
	if key == "" {
		return nil, fmt.Errorf("api-key cannot be empty")
	}

	return &IamApiKey{
		Apikey:          key,
		testEndpoint:    testEndpoint,
		privateEndpoint: privateEndpoint,
	}, nil
}

func (a *IamApiKey) GetToken(ctx context.Context) (string, error) {
	// Fast path
	if a.token != "" && time.Since(a.generatedAt) < 30*time.Minute {
		return a.token, nil
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	// Double-check
	if a.token != "" && time.Since(a.generatedAt) < 30*time.Minute {
		return a.token, nil
	}

	// Build form body
	data := url.Values{}
	data.Set("grant_type", "urn:ibm:params:oauth:grant-type:apikey")
	data.Set("apikey", a.Apikey)

	urlStr := "https://"
	if a.privateEndpoint {
		urlStr += "private."
	}
	urlStr += "iam."
	if a.testEndpoint {
		urlStr += "test."
	}
	urlStr += "cloud.ibm.com/identity/token"

	timeout := int64(5000) // 5 seconds

	req := &api.HttpRequest{
		Method: "POST",
		URL:    urlStr,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body:      []byte(data.Encode()),
		TimeoutMs: &timeout,
		RetryDelayFunction: func(statusCode int, err error, attempt int) *time.Duration {
			// simple retry: retry on 5xx or network error
			if attempt >= 3 {
				return nil
			}
			if err != nil || statusCode >= 500 {
				d := time.Duration(100*(attempt+1)) * time.Millisecond
				return &d
			}
			return nil
		},
	}

	resp, err := api.Execute(req)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", fmt.Errorf("IAM call failed: empty response")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf(
			"IAM call failed: status: %d\nbody: %s",
			resp.StatusCode,
			string(resp.Body),
		)
	}

	var parsed struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.Unmarshal(resp.Body, &parsed); err != nil {
		return "", err
	}

	if parsed.AccessToken == "" {
		return "", fmt.Errorf("access_token missing in response")
	}

	a.token = "Bearer " + parsed.AccessToken
	a.generatedAt = time.Now()

	return a.token, nil
}
