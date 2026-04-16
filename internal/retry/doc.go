// Package retry provides a simple configurable retry mechanism with
// exponential backoff for use in portwatch components that perform
// fallible operations such as webhook delivery or state persistence.
//
// Usage:
//
//	p := retry.DefaultPolicy()
//	err := retry.Do(ctx, p, func() error {
//		return sendWebhook(event)
//	})
package retry
