package checkers

import (
	"crypto/tls"
	"net"
	"time"
)

func SSLCheck(rawAddr string, timeout int) (status string, duration int64, err error) {
	start := time.Now()

	addr, errNormalize := NormalizeAddress(rawAddr, "443")
	if errNormalize != nil {
		return "down", 0, errNormalize
	}

	dialer := &net.Dialer{Timeout: time.Duration(timeout) * time.Second}

	conn, err := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
		InsecureSkipVerify: true, // не проверяем сертификат, просто проверяем соединение
	})
	duration = time.Since(start).Milliseconds()

	if err != nil {
		return "down", duration, err
	}
	defer conn.Close()

	return "up", duration, nil
}
