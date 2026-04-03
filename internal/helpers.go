package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func getModuleName(config map[string]any) string {
	if v, ok := config["module"].(string); ok && v != "" {
		return v
	}
	return "okta"
}

func resolveValue(key string, current, config map[string]any) string {
	if v, ok := current[key].(string); ok && v != "" {
		return v
	}
	if v, ok := config[key].(string); ok && v != "" {
		return v
	}
	return ""
}

func resolveStringSlice(key string, current, config map[string]any) []string {
	for _, m := range []map[string]any{current, config} {
		switch v := m[key].(type) {
		case []string:
			return v
		case []any:
			result := make([]string, 0, len(v))
			for _, item := range v {
				if s, ok := item.(string); ok {
					result = append(result, s)
				}
			}
			return result
		}
	}
	return nil
}

func resolveMap(key string, current, config map[string]any) map[string]any {
	if v, ok := current[key].(map[string]any); ok {
		return v
	}
	if v, ok := config[key].(map[string]any); ok {
		return v
	}
	return nil
}

func resolveBool(key string, current, config map[string]any) bool {
	for _, m := range []map[string]any{current, config} {
		switch v := m[key].(type) {
		case bool:
			return v
		case string:
			return v == "true" || v == "1" || v == "yes"
		}
	}
	return false
}

func resolveInt(key string, current, config map[string]any) int {
	if v := toInt64(current[key]); v != 0 {
		return int(v)
	}
	return int(toInt64(config[key]))
}

func toInt64(v any) int64 {
	switch t := v.(type) {
	case int64:
		return t
	case int:
		return int64(t)
	case int32:
		return int64(t)
	case float64:
		return int64(t)
	case float32:
		return int64(t)
	case string:
		var n int64
		fmt.Sscanf(t, "%d", &n)
		return n
	}
	return 0
}

// getHTTPClient returns the SDK-managed HTTP client, falling back to
// http.DefaultClient when no SDK client is configured (e.g. unit tests).
func getHTTPClient(client *OktaClient) *http.Client {
	if client.SdkClient != nil {
		return client.SdkClient.GetConfig().HTTPClient
	}
	return http.DefaultClient
}

// rateLimitWait parses Okta's X-Rate-Limit-Reset header (unix timestamp)
// to determine how long to wait before retrying a 429 response.
func rateLimitWait(resp *http.Response) time.Duration {
	if reset := resp.Header.Get("X-Rate-Limit-Reset"); reset != "" {
		if ts, err := strconv.ParseInt(reset, 10, 64); err == nil {
			wait := time.Until(time.Unix(ts, 0))
			if wait > 0 && wait < 60*time.Second {
				return wait + time.Second
			}
		}
	}
	return 2 * time.Second
}

// oktaRequest performs an authenticated HTTP request to the Okta API.
// Uses the SDK's HTTP client for transport and adds automatic rate-limit retry.
func oktaRequest(client *OktaClient, method, path string, body map[string]any, queryParams url.Values) (any, int, error) {
	endpoint := client.OrgURL + path
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("okta: failed to marshal request body: %w", err)
		}
	}

	httpClient := getHTTPClient(client)
	const maxRetries = 3

	for attempt := 0; attempt <= maxRetries; attempt++ {
		var reqBody io.Reader
		if bodyBytes != nil {
			reqBody = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequest(method, endpoint, reqBody)
		if err != nil {
			return nil, 0, fmt.Errorf("okta: failed to create request: %w", err)
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "SSWS "+client.APIToken)

		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, 0, fmt.Errorf("okta: request failed: %w", err)
		}

		if resp.StatusCode == http.StatusTooManyRequests && attempt < maxRetries {
			resp.Body.Close()
			time.Sleep(rateLimitWait(resp))
			continue
		}

		respData, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, resp.StatusCode, fmt.Errorf("okta: failed to read response: %w", err)
		}

		if resp.StatusCode >= 400 {
			return nil, resp.StatusCode, fmt.Errorf("okta: API error %d: %s", resp.StatusCode, string(respData))
		}

		if len(respData) == 0 || string(respData) == "null" {
			return nil, resp.StatusCode, nil
		}

		var result any
		if err := json.Unmarshal(respData, &result); err != nil {
			return nil, resp.StatusCode, fmt.Errorf("okta: failed to parse response: %w", err)
		}

		return result, resp.StatusCode, nil
	}

	return nil, http.StatusTooManyRequests, fmt.Errorf("okta: rate limited after %d retries", maxRetries)
}

// oktaGet performs a GET request to the Okta API.
func oktaGet(client *OktaClient, path string, queryParams url.Values) (any, error) {
	result, _, err := oktaRequest(client, http.MethodGet, path, nil, queryParams)
	return result, err
}

// oktaPost performs a POST request to the Okta API.
func oktaPost(client *OktaClient, path string, body map[string]any) (any, error) {
	result, _, err := oktaRequest(client, http.MethodPost, path, body, nil)
	return result, err
}

// oktaPut performs a PUT request to the Okta API.
func oktaPut(client *OktaClient, path string, body map[string]any) (any, error) {
	result, _, err := oktaRequest(client, http.MethodPut, path, body, nil)
	return result, err
}

// oktaDelete performs a DELETE request to the Okta API.
func oktaDelete(client *OktaClient, path string) error {
	_, _, err := oktaRequest(client, http.MethodDelete, path, nil, nil)
	return err
}

// oktaPostEmpty performs a POST with no body (for lifecycle operations like activate, deactivate).
func oktaPostEmpty(client *OktaClient, path string) (any, error) {
	result, _, err := oktaRequest(client, http.MethodPost, path, map[string]any{}, nil)
	return result, err
}

// toMap converts an any to map[string]any safely.
func toMap(v any) map[string]any {
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return nil
}

// toSlice converts an any to []any safely.
func toSlice(v any) []any {
	if s, ok := v.([]any); ok {
		return s
	}
	return nil
}

// mapResult wraps a map result in a StepResult output.
func mapResult(m map[string]any) map[string]any {
	if m == nil {
		return map[string]any{}
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// listResult wraps a slice result in a StepResult output with items + count.
func listResult(key string, items []any) map[string]any {
	return map[string]any{
		key:     items,
		"count": len(items),
	}
}

// errResult returns a StepResult output with an error message.
func errResult(msg string) map[string]any {
	return map[string]any{"error": msg}
}
