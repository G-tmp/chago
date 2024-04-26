package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/websocket"
)


type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte
    msg  *Msg
}

var upgrader = &websocket.Upgrader{
    ReadBufferSize: 512,
    WriteBufferSize: 512,
    CheckOrigin: func(r *http.Request) bool { return true },
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
        msg: &Msg{},
        hub: hub,
    }

    client.hub.register <- client
    go client.writer()
    go client.reader()
}


func (client *Client) writer() {
    for msg := range client.send {
        client.conn.WriteMessage(websocket.TextMessage, msg)
    }
    client.conn.Close()
}


func (client *Client) reader() {
    for {
        _, m, err := client.conn.ReadMessage()
        if err != nil {
            log.Println(err)
            break
        }

        json.Unmarshal(m, &client.msg)

        switch client.msg.Type {
        case "login":
            client.hub.user_list = append(client.hub.user_list, client.msg.User)
            // user_list = append(user_list, client.msg.User)
            client.msg.UserList = client.hub.user_list
            client.msg.Timestamp = time.Now().Unix()
            client.hub.broadcast <- client
        case "user":
            client.msg.Timestamp = time.Now().Unix()
            client.hub.broadcast <- client
        case "logout":
            // c.msg.Timestamp = time.Now().Unix()
            // user_list = del(user_list, c.msg.User)
            // c.msg.UserList = user_list
            // jd, _ := json.Marshal(c.msg)
            // client.hub.unregister <- c
            // client.hub.broadcast <- jd
        default:
            log.Print("unknown type, ditch")
        }
    }

    defer func() {
        client.msg.Type = "logout"
        client.hub.user_list = del(client.hub.user_list, client.msg.User)
        client.msg.UserList = client.hub.user_list
        client.msg.Timestamp = time.Now().Unix()
        client.hub.broadcast <- client
        client.hub.unregister <- client
    }()
}


func del(slice []string, user string) []string {
    count := len(slice)
    if count == 0 {
        return slice
    }
    if count == 1 && slice[0] == user {
        return []string{}
    }
    var n_slice = []string{}
    for i := range slice {
        if slice[i] == user && i == count {
            return slice[:count]
        } else if slice[i] == user {
            n_slice = append(slice[:i], slice[i+1:]...)
            break
        }
    }
    return n_slice
}


func init(){
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}