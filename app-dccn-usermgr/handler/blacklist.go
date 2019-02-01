package handler

import (
	"sync"
	"time"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
)

// Blacklist user for  logout
type Blacklist struct {
	mu             *sync.Mutex
	caches         map[string]int64
	checkInterval  int // minute
	tokenValidTime int
	term           chan struct{}
}

func NewBlacklist() *Blacklist {
	blacklist := &Blacklist{
		mu:             new(sync.Mutex),
		caches:         make(map[string]int64),
		checkInterval:  5,
		tokenValidTime: ankr_default.AccessTokenValidTime,
		term:           make(chan struct{}),
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
	if _, ok := p.caches[token]; ok {
		delete(p.caches, token)
	}
	p.mu.Unlock()
}

func (p *Blacklist) Available(token string) bool {

	lastAccessTime, ok := p.caches[token]
	return ok && time.Now().Unix()-lastAccessTime < int64(p.tokenValidTime*60)
}

func (p *Blacklist) Refresh(token string) {

	p.mu.Lock()
	p.caches[token] = time.Now().Unix()
	p.mu.Unlock()
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
				if (time.Now().Unix() - startTime) >= int64(p.tokenValidTime*60) {
					p.Remove(token)
				}
			}
			time.Sleep(time.Duration(p.checkInterval) * time.Minute)
		}
	}
}
