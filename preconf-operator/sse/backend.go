// Package sse is the SSE Client for preconf-share
package sse

// SSEClient is the SSE Client abstraction
type SSEClient interface {
	// Subscribe to events and returns a subscription
	Subscribe(eventChan chan<- Event) (SSESubscription, error)
}

type SSESubscription interface {
	// To stop the subscription
	Stop()
}
