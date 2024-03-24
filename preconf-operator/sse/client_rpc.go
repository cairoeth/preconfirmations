package sse

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// InternalClient is a client for the matchmaker
type InternalClient struct {
	BaseURL string // BaseURL is the base URL for the matchmaker
}

// New creates a new InternalClient for the matchmaker with the given base URL
func New(baseURL string) SSEClient {
	return &InternalClient{
		BaseURL: baseURL,
	}
}

// Subscription represents a subscription to matchmaker events
type Subscription struct {
	eventChan chan<- Event
	pubsub    *redis.PubSub
}

// Subscribe to matchmaker events and returns a type that can be used to control the subscription
func (c *InternalClient) Subscribe(eventChan chan<- Event) (SSESubscription, error) {
	rds := redis.NewClient(&redis.Options{
		Addr: c.BaseURL,
	})

	// There is no error because go-redis automatically reconnects on error.
	pubsub := rds.Subscribe(context.Background(), "hints")

	sub := &Subscription{
		pubsub:    pubsub,
		eventChan: eventChan,
	}

	go sub.readEvents()

	return sub, nil
}

// readEvents reads the events and sends them to the event channel
func (s *Subscription) readEvents() {
	for {
		msg, err := s.pubsub.ReceiveMessage(context.Background())
		if err != nil {
			fmt.Printf("Error occurred receiving message from preconf-share stream %v", err)
		}

		var event MatchMakerEvent
		err = json.Unmarshal([]byte(msg.Payload), &event)

		if err != nil {
			s.eventChan <- Event{
				Error: err,
			}
		} else {
			s.eventChan <- Event{
				Data: &event,
			}
		}
	}
}

// Stop stops the subscription to matchmaker events
func (s *Subscription) Stop() {
	s.pubsub.Close()
}
