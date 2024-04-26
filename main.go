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
    
    port := "12345"
    fmt.Printf("http://127.0.0.1:%v \n", port)
    if err := http.ListenAndServe(":" + port, nil); err != nil {
        fmt.Println("err:", err)
    }
}
