package handler

import (
	"sync"
	"time"
)

type Blacklist struct {
	mu            *sync.Mutex
	caches        map[string]int64
	checkInterval int // minute
	term          chan struct{}
}

func NewBlacklist() *Blacklist {
	blacklist := &Blacklist{
		mu:            new(sync.Mutex),
		caches:        make(map[string]int64),
		checkInterval: 5,
		term:          make(chan struct{}),
	}
	go blacklist.check()

	return blacklist
}

func (p *Blacklist) Add(token string) {

	p.mu.Lock()
	p.caches[token] = time.Now().Unix()
	p.mu.Unlock()
}

func (p *Blacklist) Remove(token string) {

	p.mu.Lock()
	if p.Available(token) {
		delete(p.caches, token)
	}
	p.mu.Unlock()
}

func (p *Blacklist) Available(token string) bool {

	_, ok := p.caches[token]
	return ok
}

func (p *Blacklist) destroy() {
	close(p.term)
	p.caches = nil
}

func (p *Blacklist) check() {
	for {
		select {
		case <-p.term:
			return

		default:
			for token, startTime := range p.caches {
				if (time.Now().Unix() - startTime) >= int64(p.checkInterval*60) {
					p.Remove(token)
				}
			}
			time.Sleep(time.Duration(5) * time.Minute)
		}
	}
}
