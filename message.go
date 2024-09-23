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

    TextType        = "text"
    ImageType       = "image"
)

type Msg struct {
    Action    string   `json:"action"`
    Type      string   `json:"type"`
    Sender    string   `json:"sender"`
    Content   string   `json:"content"`
    Timestamp int64    `json:"timestamp"`
    UserList  []string `json:"user_list"`
    // Self      bool     `json:"self"`
    // Target    string   `json:"target"`
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