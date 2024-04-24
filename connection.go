package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/websocket"
)


var user_list = []string{}


type connection struct {
    ws   *websocket.Conn
    send   chan []byte
    msg *Msg
}

var wu = &websocket.Upgrader{
    ReadBufferSize: 512,
    WriteBufferSize: 512,
    CheckOrigin: func(r *http.Request) bool { return true },
}


func myws(w http.ResponseWriter, r *http.Request) {
    cookie := http.Cookie{Name: "wsId", Value: "111", Path: "/ws"}
    http.SetCookie(w, &cookie)
    header := http.Header{}
    header.Add("Set-Cookie", cookie.String())

    ws, err := wu.Upgrade(w, r, header)
    if err != nil {
        log.Println(err)
        return
    }
    c := &connection{send: make(chan []byte, 256), ws: ws, msg: &Msg{}}

    h.register <- c
    go c.writer()
    c.reader()
    defer func() {
        c.msg.Type = "logout"
        user_list = del(user_list, c.msg.User)
        c.msg.UserList = user_list
        c.msg.Timestamp = time.Now().Unix()
        jd, _ := json.Marshal(c.msg)
        h.broadcast <- jd
        h.unregister <- c
    }()
}


func (c *connection) writer() {
    for msg := range c.send {
        c.ws.WriteMessage(websocket.TextMessage, msg)
    }
    c.ws.Close()
}


func (c *connection) reader() {
    for {
        _, msg, err := c.ws.ReadMessage()
        if err != nil {
            h.unregister <- c
            break
        }
        json.Unmarshal(msg, &c.msg)

        switch c.msg.Type {
        case "login":
            user_list = append(user_list, c.msg.User)
            // c.data.User = c.data.Content
            c.msg.UserList = user_list
            c.msg.Timestamp = time.Now().Unix()
            jd, _ := json.Marshal(c.msg)
            h.broadcast <- jd
        case "user":
            c.msg.Timestamp = time.Now().Unix()
            jd, _ := json.Marshal(c.msg)
            h.broadcast <- jd
        case "logout":
            // c.msg.Timestamp = time.Now().Unix()
            // user_list = del(user_list, c.msg.User)
            // c.msg.UserList = user_list
            // jd, _ := json.Marshal(c.msg)
            // h.unregister <- c
            // h.broadcast <- jd
        default:
            log.Print("unknown type, ditch")
        }
    }
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