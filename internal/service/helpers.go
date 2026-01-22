package service

import (
	"monitor-bot/internal/models"
	"net/url"
)

func GetMaxSubscriptions(user *models.User) int {
	if user.IsPremium {
		return 10 // Premium план
	}
	return 1 // Free план
}

func IsValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	if u.Scheme == "" || u.Host == "" {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	return true
}
