package checkers

import (
	"io"
	"net/http"
	"time"
)

func HTTPCheck(url string, expectedStatus int, timeout int) (status string, code int, duration int64, err error) {
	client := http.Client{Timeout: time.Duration(timeout) * time.Second}
	start := time.Now()
	resp, err := client.Get(url)
	duration = time.Since(start).Milliseconds()
	if err != nil {
		return "down", 0, duration, err
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)

	if resp.StatusCode == expectedStatus {
		return "up", resp.StatusCode, duration, nil
	}
	return "down", resp.StatusCode, duration, nil
}
