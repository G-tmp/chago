package main

import (
    "sync"
)


type clientSet struct {
    mu      sync.RWMutex
    numbers map[*Client]struct{}
}

func (cs *clientSet) add(c *Client){
    cs.mu.Lock()
    defer cs.mu.Unlock()

    cs.numbers[c] = struct{}{}
}

func (cs *clientSet) del(c *Client){
    cs.mu.Lock()
    defer cs.mu.Unlock()

    delete(cs.numbers, c)
}

func (cs *clientSet) size() int {
    cs.mu.RLock()
    defer cs.mu.RUnlock()

    return len(cs.numbers)
}

func (cs *clientSet) clear(){
    cs.mu.Lock()
    defer cs.mu.Unlock()

    cs.numbers = make(map[*Client]struct{})
}

func (cs *clientSet) all() []*Client {
    cs.mu.RLock()
    defer cs.mu.RUnlock()

    list := make([]*Client, 0, len(cs.numbers))
    for i := range cs.numbers {
        list = append(list, i)
    }

    return list
}

func (cs *clientSet) each(f func(c *Client)){
    cs.mu.RLock()
    defer cs.mu.RUnlock()

    for i := range cs.numbers {
        f(i)
    }
}



type Hub struct {
    clients     *clientSet
    broadcast   chan []byte
    register    chan *Client
    unregister  chan *Client
    bf          func()
}


func newHub() *Hub {
    return &Hub{
        broadcast:  make(chan []byte, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        clients:    &clientSet{
            numbers: make(map[*Client]struct{}),
        },
        bf:         func(){},
    }
}

func (hub *Hub) broadcastF(f func()){
    hub.bf = f
}

func (hub *Hub) run() {
    for {
        select {
        case c := <-hub.register:
            hub.clients.add(c)
        case c := <-hub.unregister:
            hub.clients.del(c)
        case  m := <-hub.broadcast:
            // hub.bf()

            hub.clients.each(func(c *Client){
                c.send <- m
            })
        }
    }
}