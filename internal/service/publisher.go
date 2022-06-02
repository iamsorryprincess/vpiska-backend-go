package service

import "sync"

type publisher struct {
	mutex         sync.Mutex
	subscriptions map[string][]Subscriber
}

func newPublisher() Publisher {
	return &publisher{
		mutex:         sync.Mutex{},
		subscriptions: make(map[string][]Subscriber),
	}
}

func (p *publisher) Subscribe(eventId string, subscriber Subscriber) {
	p.mutex.Lock()
	p.subscriptions[eventId] = append(p.subscriptions[eventId], subscriber)
	p.mutex.Unlock()
}

func (p *publisher) Unsubscribe(eventId string, subscriber Subscriber) {
	for i, sub := range p.subscriptions[eventId] {
		if sub == subscriber {
			p.mutex.Lock()
			p.subscriptions[eventId] = append(p.subscriptions[eventId][:i], p.subscriptions[eventId][i+1:]...)
			p.mutex.Unlock()
			return
		}
	}
}

func (p *publisher) Publish(eventId string, message []byte) {
	for _, subscriber := range p.subscriptions[eventId] {
		subscriber.OnReceive(message)
	}
}

func (p *publisher) Close(eventId string) {
	temp := make([]Subscriber, len(p.subscriptions[eventId]))

	for i, item := range p.subscriptions[eventId] {
		temp[i] = item
	}

	p.mutex.Lock()
	delete(p.subscriptions, eventId)
	p.mutex.Unlock()

	for _, subscriber := range temp {
		subscriber.OnClose()
	}
}
