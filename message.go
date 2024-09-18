package main

import(
    "encoding/json"
    "log"
)


const (
    SendMessageAction = "send-message"
    UserJoinedAction  = "user-join"
    UserLeftAction    = "user-left"
    JoinRoomAction    = "join-room"
    LeaveRoomAction   = "leave-room"
)

type Msg struct {
    // Self      bool     `json:"self"`
    Action    string   `json:"action"`
    Sender    string   `json:"sender"`
    Target    string   `json:"target"`
    Content   string   `json:"content"`
    Timestamp int64    `json:"timestamp"`
    UserList  []string `json:"user_list"`
}

func (msg *Msg) encode() []byte {
    json, err := json.Marshal(msg)
    if err != nil {
        log.Println(err)
        return nil
    }

    return json
}

func decode(jsondata []byte) *Msg {
    var msg *Msg
    err := json.Unmarshal(jsondata, &msg)
    if err != nil {
        log.Println(err)
        return nil
    }

    return msg
}