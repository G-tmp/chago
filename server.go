package main

import (
    "fmt"
    "net/http"

)

func main() {
    http.HandleFunc("/chat", myws)
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request){
        http.ServeFile(w, r, "chat.html")
    })
    go h.run()
    if err := http.ListenAndServe("127.0.0.1:12345", nil); err != nil {
        fmt.Println("err:", err)
    }
}
