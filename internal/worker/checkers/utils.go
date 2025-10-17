package checkers

import (
	"net/url"
	"strings"
)

// NormalizeAddress приводит URL к виду host:port
func NormalizeAddress(raw string, defaultPort string) (string, error) {
	if strings.Contains(raw, "://") {
		parsed, err := url.Parse(raw)
		if err != nil {
			return "", err
		}
		host := parsed.Host
		if !strings.Contains(host, ":") {
			host += ":" + defaultPort
		}
		return host, nil
	}

	if !strings.Contains(raw, ":") {
		raw += ":" + defaultPort
	}
	return raw, nil
}
