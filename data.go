package main

type Msg struct {
    User     string   `json:"user"`
    Type     string   `json:"type"`
    Content  string   `json:"content"`
    Timestamp int64   `json:"timestamp"`
    UserList []string `json:"user_list"`
}
