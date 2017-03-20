package tools

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"

	"cloud.google.com/go/pubsub"
)

// Pubsub ...
type Pubsub struct {
	ProjectID    string
	Topic        string
	Subscription string
	client       *pubsub.Client
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
	iter         *pubsub.MessageIterator
}

// Initial ...
func (p Pubsub) Initial() (*pubsub.MessageIterator, error) {
	ctx := context.Background()
	var err error
	err = Retry(5, func() error {
		p.client, err = pubsub.NewClient(ctx, p.ProjectID)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	err = Retry(20, func() (err error) {
		isExists, err := p.client.Topic(p.Topic).Exists(ctx)
		if err != nil {
			return errors.Wrapf(err, "Error occurs on test topic %s whether exists", p.Topic)
		}
		if isExists {
			p.topic = p.client.Topic(p.Topic)
			return nil
		}
		p.topic, err = p.client.CreateTopic(ctx, p.Topic)
		if err != nil {
			return err
		}
		return err
	})

	if err != nil {
		log.Fatal(err)
	}

	err = Retry(20, func() (err error) {
		isExists, err := p.client.Subscription(p.Subscription).Exists(ctx)
		if err != nil {
			return errors.Wrapf(err, "Error occurs on test topic %s whether exists", p.Subscription)
		}
		if isExists {
			p.subscription = p.client.Subscription(p.Subscription)
		} else {
			p.subscription, err = p.client.CreateSubscription(ctx, p.Subscription, p.topic, 10*time.Second, nil)
		}
		return err
	})

	if err != nil {
		log.Fatal(err)
	}

	p.iter, err = p.subscription.Pull(ctx)
	if err != nil {
		return nil, err
	}
	return p.iter, nil
}
