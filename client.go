package main

import (
    "log"
    "net/http"
    "time"
    "github.com/gorilla/websocket"
)


type Client struct {
    hub     *Hub
    conn    *websocket.Conn
    send    chan []byte
    nickname   string
}

var upgrader = &websocket.Upgrader{
    ReadBufferSize: 512,
    WriteBufferSize: 512,
    // CheckOrigin: func(r *http.Request) bool { return true },
}


func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
    // cookie := http.Cookie{Name: "session", Value: "111", Path: "/ws"}
    // http.SetCookie(w, &cookie)
    // header := http.Header{}
    // header.Add("Set-Cookie", cookie.String())

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    
    client := &Client{
        send: make(chan []byte, 256),
        conn: conn, 
        hub: hub,
    }

    client.hub.register <- client

    go client.writePump()
    go client.readPump()
}


func (client *Client) writePump() {
    for msg := range client.send {
        if err := client.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
            log.Println(err)
        }
    }

    client.conn.Close()
}


func (client *Client) readPump() {
    for {
        _, m, err := client.conn.ReadMessage()
        if err != nil {
            log.Println(err)
            break
        }

        var msg *Msg 
        msg = decode(m)

        switch msg.Action {
        case SendMessageAction:
            msg.Timestamp = time.Now().Unix()
            jd := msg.encode()

            client.hub.broadcast <- jd
        case UserJoinedAction:
            client.nickname = msg.Sender
            msg.Timestamp = time.Now().Unix()
            for _, c := range client.hub.clients.all() {
                msg.UserList = append(msg.UserList, c.nickname)
            }
            jd := msg.encode()

            client.hub.broadcast <- jd
        // case UserLeftAction:
            
        default:
            log.Print("unknown type, ditch")
        }
    }

    defer func() {
        // 2. other goroutine, call later
        client.hub.unregister <- client
        close(client.send)
        client.conn.Close()

        msg := &Msg{
            Action: UserLeftAction,
            Sender: client.nickname,
            Timestamp: time.Now().Unix(),
        }
        // 1. call and get clients map first
        for _, c := range client.hub.clients.all() {
            if c != client {
                msg.UserList = append(msg.UserList, c.nickname)
            }
        }

        jd := msg.encode()
        client.hub.broadcast <- jd

    }()
}
