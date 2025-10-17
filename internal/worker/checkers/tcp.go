package checkers

import (
	"net"
	"time"
)

func TCPCheck(rawAddr string, timeout int) (status string, duration int64, err error) {
	start := time.Now()

	addr, errNormalize := NormalizeAddress(rawAddr, "80")
	if errNormalize != nil {
		return "down", 0, errNormalize
	}

	conn, err := net.DialTimeout("tcp", addr, time.Duration(timeout)*time.Second)
	duration = time.Since(start).Milliseconds()

	if err != nil {
		return "down", duration, err
	}
	defer conn.Close()

	return "up", duration, nil
}
