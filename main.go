package main

import (
    "fmt"
    "net/http"
)


func main() {
    hub := newHub()
    go hub.run()

    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request){
        http.ServeFile(w, r, "chat.html")
    })
    http.HandleFunc("/chat", func (w http.ResponseWriter, r *http.Request){
        serveWs(hub, w, r)
    })
    
    if err := http.ListenAndServe(":12345", nil); err != nil {
        fmt.Println("err:", err)
    }
}
