package main

import "encoding/json"


type Hub struct {
    clients     map[*Client]bool
    broadcast   chan *Client
    register    chan *Client
    unregister  chan *Client
    user_list   []string
}


func newHub() *Hub {
    return &Hub{
        broadcast:  make(chan *Client),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        clients:    make(map[*Client]bool),
        user_list:  []string{},
    }
}


func (hub *Hub) run() {
    for {
        select {
        case c := <-hub.register:
            hub.clients[c] = true
            c.msg.Type = "handshake"
            jd, _ := json.Marshal(c.msg)
            c.send <- jd
        case c := <-hub.unregister:
            if _, ok := hub.clients[c]; ok {
                delete(hub.clients, c)
                close(c.send)
                c.conn.Close()
            }
        case c := <-hub.broadcast:
            for client := range hub.clients {
                if c == client {
                    c.msg.Self = true
                }else {
                    c.msg.Self = false    
                }
                
                jd, _ := json.Marshal(c.msg)

                select {
                case client.send <- jd:
                
                default:
                    delete(hub.clients, client)
                    close(client.send)
                    c.conn.Close()
                }
            }
 
        }
    }
}
