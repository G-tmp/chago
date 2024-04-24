package main

import "encoding/json"


var h = hub{
    clients: make(map[*connection]bool),
    broadcast: make(chan []byte),
    register: make(chan *connection),
    unregister: make(chan *connection),
}

type hub struct {
    clients     map[*connection]bool
    broadcast   chan []byte
    register    chan *connection
    unregister  chan *connection
}


func (h *hub) run() {
    for {
        select {
        case c := <-h.register:
            h.clients[c] = true
            c.msg.Type = "handshake"
            jd, _ := json.Marshal(c.msg)
            c.send <- jd
        case c := <-h.unregister:
            if _, ok := h.clients[c]; ok {
                delete(h.clients, c)
                close(c.send)
            }
        case msg := <-h.broadcast:
            for c := range h.clients {
                select {
                case c.send <- msg:

                default:
                    delete(h.clients, c)
                    close(c.send)
                }
            }
        }
    }
}
