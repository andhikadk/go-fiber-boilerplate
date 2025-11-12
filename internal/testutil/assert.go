package testutil

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// AssertEqual checks if two values are equal
func AssertEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		message := fmt.Sprintf("Expected: %v\nActual: %v", expected, actual)
		if len(msgAndArgs) > 0 {
			message = fmt.Sprintf("%v\n%s", msgAndArgs[0], message)
		}
		t.Error(message)
	}
}

// AssertNotEqual checks if two values are not equal
func AssertNotEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if reflect.DeepEqual(expected, actual) {
		message := fmt.Sprintf("Did not expect: %v", actual)
		if len(msgAndArgs) > 0 {
			message = fmt.Sprintf("%v\n%s", msgAndArgs[0], message)
		}
		t.Error(message)
	}
}

// AssertNil checks if a value is nil
func AssertNil(t *testing.T, value interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if value != nil {
		message := fmt.Sprintf("Expected nil, but got: %v", value)
		if len(msgAndArgs) > 0 {
			message = fmt.Sprintf("%v\n%s", msgAndArgs[0], message)
		}
		t.Error(message)
	}
}

// AssertNotNil checks if a value is not nil
func AssertNotNil(t *testing.T, value interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if value == nil {
		message := "Expected non-nil value"
		if len(msgAndArgs) > 0 {
			message = fmt.Sprintf("%v", msgAndArgs[0])
		}
		t.Error(message)
	}
}

// AssertTrue checks if a condition is true
func AssertTrue(t *testing.T, condition bool, msgAndArgs ...interface{}) {
	t.Helper()
	if !condition {
		message := "Expected condition to be true"
		if len(msgAndArgs) > 0 {
			message = fmt.Sprintf("%v", msgAndArgs[0])
		}
		t.Error(message)
	}
}

// AssertFalse checks if a condition is false
func AssertFalse(t *testing.T, condition bool, msgAndArgs ...interface{}) {
	t.Helper()
	if condition {
		message := "Expected condition to be false"
		if len(msgAndArgs) > 0 {
			message = fmt.Sprintf("%v", msgAndArgs[0])
		}
		t.Error(message)
	}
}

// AssertStatusCode checks if the HTTP response status code matches expected
func AssertStatusCode(t *testing.T, expected int, resp *http.Response, msgAndArgs ...interface{}) {
	t.Helper()
	if resp.StatusCode != expected {
		message := fmt.Sprintf("Expected status code: %d\nActual status code: %d", expected, resp.StatusCode)
		if len(msgAndArgs) > 0 {
			message = fmt.Sprintf("%v\n%s", msgAndArgs[0], message)
		}
		t.Error(message)
	}
}

// AssertLen checks if a slice/array/map has the expected length
func AssertLen(t *testing.T, value interface{}, expectedLen int, msgAndArgs ...interface{}) {
	t.Helper()
	v := reflect.ValueOf(value)
	if v.Len() != expectedLen {
		message := fmt.Sprintf("Expected length: %d\nActual length: %d", expectedLen, v.Len())
		if len(msgAndArgs) > 0 {
			message = fmt.Sprintf("%v\n%s", msgAndArgs[0], message)
		}
		t.Error(message)
	}
}

// AssertContains checks if a string contains a substring
func AssertContains(t *testing.T, haystack, needle string, msgAndArgs ...interface{}) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		message := fmt.Sprintf("Expected string to contain: %q\nActual string: %q", needle, haystack)
		if len(msgAndArgs) > 0 {
			message = fmt.Sprintf("%v\n%s", msgAndArgs[0], message)
		}
		t.Error(message)
	}
}

// AssertNotContains checks if a string does not contain a substring
func AssertNotContains(t *testing.T, haystack, needle string, msgAndArgs ...interface{}) {
	t.Helper()
	if strings.Contains(haystack, needle) {
		message := fmt.Sprintf("Expected string not to contain: %q\nActual string: %q", needle, haystack)
		if len(msgAndArgs) > 0 {
			message = fmt.Sprintf("%v\n%s", msgAndArgs[0], message)
		}
		t.Error(message)
	}
}

// AssertNoError checks if an error is nil
func AssertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if err != nil {
		message := fmt.Sprintf("Expected no error, but got: %v", err)
		if len(msgAndArgs) > 0 {
			message = fmt.Sprintf("%v\n%s", msgAndArgs[0], message)
		}
		t.Error(message)
	}
}

// AssertError checks if an error is not nil
func AssertError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if err == nil {
		message := "Expected an error, but got nil"
		if len(msgAndArgs) > 0 {
			message = fmt.Sprintf("%v", msgAndArgs[0])
		}
		t.Error(message)
	}
}

// ParseJSONResponse parses a JSON response body into the provided interface
func ParseJSONResponse(t *testing.T, body io.Reader, target interface{}) {
	t.Helper()
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if err := json.Unmarshal(bodyBytes, target); err != nil {
		t.Fatalf("Failed to parse JSON response: %v\nBody: %s", err, string(bodyBytes))
	}
}

// GetJSONField extracts a field value from a JSON response
func GetJSONField(t *testing.T, body io.Reader, field string) interface{} {
	t.Helper()
	var result map[string]interface{}
	ParseJSONResponse(t, body, &result)
	return result[field]
}
