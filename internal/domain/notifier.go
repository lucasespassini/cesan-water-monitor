package domain

import "context"

type Notifier interface {
	NotifyWaterOutage(
		ctx context.Context,
		title string,
		url string,
	) error
}
