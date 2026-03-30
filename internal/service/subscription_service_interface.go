package service

import "context"

type SubscriptionServiceInterface interface {
	SubscribeByURL(ctx context.Context, chatID int64, url string) error
	UnsubscribeByURL(ctx context.Context, chatID int64, url string) error
}
