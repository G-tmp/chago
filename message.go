package main

import(
    "encoding/json"
    "log"
)

const (
    TypeUserList        string = "userlist"
    TypeMessage         string = "message"
    TypeUserJoin        string = "user-join"
    TypeUserLeft        string = "user-left"
    TypeRejectUsername  string = "rejectUsername"
    TypeUploadImage     string = "upload-image"
)


type Msg struct {
    Type      string   `json:"type"`
    Sender    string   `json:"sender"`
    Content   string   `json:"content"`
    Timestamp int64    `json:"timestamp"`
    Self      bool     `json:"self"`
    // Id        string   `json:"id"`
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