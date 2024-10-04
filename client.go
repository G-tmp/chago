package main

import (
    "log"
    "net/http"
    "time"
    "encoding/json"
    "github.com/gorilla/websocket"
    "github.com/google/uuid"
)


type Client struct {
    hub     *Hub
    conn    *websocket.Conn
    send    chan []byte
    user    *User
}

var upgrader = &websocket.Upgrader{
    ReadBufferSize: 512,
    WriteBufferSize: 512,
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    
    client := &Client{
        send: make(chan []byte, 256),
        conn: conn, 
        hub: hub,
        user: &User{},
    }

    client.hub.register <- client

    go client.writePump()
    go client.readPump()
}


func (client *Client) writePump() {
    for msg := range client.send {
        if err := client.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
            log.Println(err)
            break
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

        switch msg.Type {
        case TypeMessage:
            msg.Timestamp = time.Now().Unix()

            jd := msg.encode()
            client.hub.broadcastE(client, jd) 

            msg.Self = true
            jd = msg.encode()
            client.send <- jd

        case TypeUserJoin:
            id := uuid.New().String()
            user := &User{
                Username: msg.Sender,
                Id: id,
            }
            client.user = user

            msg.Timestamp = time.Now().Unix()
            jd := msg.encode()
            client.hub.broadcast <- jd

            var userList []string
            for _, c := range client.hub.clients.all() {
                userList = append(userList, c.user.Username)
            }
            usermap := make(map[string]any)
            usermap["type"] = TypeUserList
            usermap["userlist"] = userList
            listJ, _ := json.Marshal(usermap)
            client.hub.broadcast <- listJ

        case TypeUploadImage:
            msg.Timestamp = time.Now().Unix()
            
            jd := msg.encode()
            client.hub.broadcastE(client, jd) 
            
            msg.Self = true
            jd = msg.encode()
            client.send <- jd
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
            Type: TypeUserLeft,
            Sender: client.user.Username,
            Timestamp: time.Now().Unix(),
        }
        jd := msg.encode()
        client.hub.broadcast <- jd

        // 1. call and get clients map first
        var userList []string
        for _, c := range client.hub.clients.all() {
            if c != client {
                userList = append(userList, c.user.Username)
            }
        }
        usermap := make(map[string]any)
        usermap["type"] = TypeUserList
        usermap["userlist"] = userList
        listJ, _ := json.Marshal(usermap)
        client.hub.broadcast <- listJ

    }()
}
