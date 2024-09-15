package utils

import (
	"fmt"
	"net/url"
)

func GetQueryString(queryParams map[string]interface{}) (string, error) {
	values := url.Values{}

	for key, value := range queryParams {
		switch v := any(value).(type) {
		case string:
			values.Add(key, v)
		case []string:
			for _, item := range v {
				values.Add(key, item)
			}
		default:
			return "", fmt.Errorf("invalid query parameter type for key %s: %T", key, value)
		}
	}

	return values.Encode(), nil
}
