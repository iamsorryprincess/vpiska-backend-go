package service

import "sync"

type publisher struct {
	mutex         sync.RWMutex
	subscriptions map[string][]Subscriber
}

func newPublisher() Publisher {
	return &publisher{
		mutex:         sync.RWMutex{},
		subscriptions: make(map[string][]Subscriber),
	}
}

func (p *publisher) Subscribe(eventId string, subscriber Subscriber) {
	p.mutex.Lock()
	p.subscriptions[eventId] = append(p.subscriptions[eventId], subscriber)
	p.mutex.Unlock()
}

func (p *publisher) Unsubscribe(eventId string, subscriber Subscriber) {
	p.mutex.RLock()
	subs := p.subscriptions[eventId]
	p.mutex.RUnlock()
	for i, sub := range subs {
		if sub == subscriber {
			p.mutex.Lock()
			p.subscriptions[eventId] = append(p.subscriptions[eventId][:i], p.subscriptions[eventId][i+1:]...)
			p.mutex.Unlock()
			return
		}
	}
}

func (p *publisher) Publish(eventId string, message []byte) {
	p.mutex.RLock()
	subs := p.subscriptions[eventId]
	p.mutex.RUnlock()
	for _, subscriber := range subs {
		subscriber.OnReceive(message)
	}
}

func (p *publisher) Close(eventId string) {
	p.mutex.RLock()
	subs := p.subscriptions[eventId]
	p.mutex.RUnlock()

	for _, subscriber := range subs {
		subscriber.OnClose()
	}

	p.mutex.Lock()
	delete(p.subscriptions, eventId)
	p.mutex.Unlock()
}

func (p *publisher) CloseAll() {
	for eventId := range p.subscriptions {
		p.Close(eventId)
	}
}
