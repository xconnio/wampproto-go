package wampproto

import (
	"fmt"
	"sync"

	"github.com/hashicorp/go-immutable-radix/v2"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/util"
)

const (
	OptAcknowledge = "acknowledge"
)

type Broker struct {
	subscriptionsByTopic   map[string]*Subscription
	subscriptionsBySession map[int64]map[int64]*Subscription
	sessions               map[int64]*SessionDetails
	prefixTree             *iradix.Tree[*Subscription]

	idGen *SessionScopeIDGenerator
	sync.Mutex
}

func NewBroker() *Broker {
	return &Broker{
		sessions:               map[int64]*SessionDetails{},
		subscriptionsByTopic:   make(map[string]*Subscription),
		subscriptionsBySession: make(map[int64]map[int64]*Subscription),
		idGen:                  &SessionScopeIDGenerator{},
		prefixTree:             iradix.New[*Subscription](),
	}
}

func (b *Broker) AddSession(details *SessionDetails) error {
	b.Lock()
	defer b.Unlock()

	_, exists := b.subscriptionsBySession[details.ID()]
	if exists {
		return fmt.Errorf("broker: cannot add session %b, it already exists", details.ID())
	}

	b.subscriptionsBySession[details.ID()] = map[int64]*Subscription{}
	b.sessions[details.ID()] = details
	return nil
}

func (b *Broker) RemoveSession(id int64) error {
	b.Lock()
	defer b.Unlock()

	subscriptions, exists := b.subscriptionsBySession[id]
	if !exists {
		return fmt.Errorf("broker: cannot remove session %b, it doesn't exist", id)
	}

	delete(b.subscriptionsBySession, id)
	for _, v := range subscriptions {
		subscription, ok := b.subscriptionsByTopic[v.Topic]
		if ok {
			delete(subscription.Subscribers, id)
		}

		if len(subscription.Subscribers) == 0 {
			delete(b.subscriptionsByTopic, v.Topic)
		}

		if subscription.Match == MatchPrefix {
			b.prefixTree.Delete([]byte(subscription.Topic))
		}
	}

	delete(b.sessions, id)

	return nil
}

func (b *Broker) HasSubscription(topic string) bool {
	b.Lock()
	defer b.Unlock()

	_, exists := b.subscriptionsByTopic[topic]
	return exists
}

func (b *Broker) ReceiveMessage(sessionID int64, msg messages.Message) (*MessageWithRecipient, error) {
	b.Lock()
	defer b.Unlock()

	switch msg.Type() {
	case messages.MessageTypeSubscribe:
		_, exists := b.subscriptionsBySession[sessionID]
		if !exists {
			return nil, fmt.Errorf("broker: cannot subscribe, session %d doesn't exist", sessionID)
		}

		subscribe := msg.(*messages.Subscribe)
		subscription, exists := b.subscriptionsByTopic[subscribe.Topic()]
		if exists {
			subscription.Subscribers[sessionID] = sessionID
		} else {
			subscription = &Subscription{
				ID:          b.idGen.NextID(),
				Topic:       subscribe.Topic(),
				Subscribers: map[int64]int64{sessionID: sessionID},
			}
			match := util.ToString(subscribe.Options()[OptionMatch])
			switch match {
			case MatchPrefix:
				subscription.Match = match
				b.prefixTree, _, _ = b.prefixTree.Insert([]byte(subscription.Topic), subscription)
			default:
				subscription.Match = MatchExact
			}
			b.subscriptionsByTopic[subscribe.Topic()] = subscription
		}

		b.subscriptionsBySession[sessionID][subscription.ID] = subscription

		subscribed := messages.NewSubscribed(subscribe.RequestID(), subscription.ID)
		result := &MessageWithRecipient{Message: subscribed, Recipient: sessionID}
		return result, nil
	case messages.MessageTypeUnsubscribe:
		unsubscribe := msg.(*messages.Unsubscribe)
		subscriptions, exists := b.subscriptionsBySession[sessionID]
		if !exists {
			return nil, fmt.Errorf("broker: cannot unsubscribe, session %d doesn't exist", sessionID)
		}

		subscription, exists := subscriptions[unsubscribe.SubscriptionID()]
		if !exists {
			return nil, fmt.Errorf("broker: cannot unsubscribe non-existent subscription %d",
				unsubscribe.SubscriptionID())
		}

		delete(subscription.Subscribers, sessionID)
		if len(subscription.Subscribers) == 0 {
			delete(b.subscriptionsByTopic, subscription.Topic)
		}
		if subscription.Match == MatchPrefix {
			b.prefixTree.Delete([]byte(subscription.Topic))
		}

		delete(b.subscriptionsBySession[sessionID], subscription.ID)

		unsubscribed := messages.NewUnsubscribed(unsubscribe.RequestID())
		result := &MessageWithRecipient{Message: unsubscribed, Recipient: sessionID}
		return result, nil
	case messages.MessageTypeError:
		return nil, fmt.Errorf("broker: error handling is not implemented yet")
	default:
		return nil, fmt.Errorf("broker: received unexpected message of type %T", msg)
	}
}

func (b *Broker) ReceivePublish(sessionID int64, publish *messages.Publish) (*Publication, error) {
	b.Lock()
	defer b.Unlock()

	_, exists := b.subscriptionsBySession[sessionID]
	if !exists {
		return nil, fmt.Errorf("broker: cannot publish, session %d doesn't exist", sessionID)
	}

	result := &Publication{}
	publicationID := b.idGen.NextID()

	var subscription *Subscription
	subscription, exists = b.subscriptionsByTopic[publish.Topic()]
	if !exists || len(subscription.Subscribers) == 0 {
		if b.prefixTree.Len() > 0 {
			_, sub, ok := b.prefixTree.Root().LongestPrefix([]byte(publish.Topic()))
			exists = ok
			if ok {
				subscription = sub
			}
		}
	}
	if exists && len(subscription.Subscribers) > 0 {
		event := messages.NewEvent(subscription.ID, publicationID, nil, publish.Args(), publish.KwArgs())
		result.Event = event
		for _, subscriber := range subscription.Subscribers {
			result.Recipients = append(result.Recipients, subscriber)
		}
	}

	ack, ok := publish.Options()[OptAcknowledge].(bool)
	if ok && ack {
		published := messages.NewPublished(publish.RequestID(), publicationID)
		result.Ack = &MessageWithRecipient{Message: published, Recipient: sessionID}
	}

	return result, nil
}
